// Package jobs provides the domain object for jobs
// dto - data transfer object; a pattern for transferring data between processes
package jobs

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"jobd/domain/status"
	"jobd/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
)

var DEBUG = os.Getenv("DEBUG") == "true"

type Job struct {
	ID          string `json:"ID"`
	Status      string
	Path        string
	Input       string `json:"Input"`
	Output      string
	Message     string
	Slurml      bool `json:"Slurml"`
	SlurmID     int
	LastUpdated time.Time
}

type JobList struct {
	Jobs []Job
}

type SlurmPostResponse struct {
	Jobid int `json:"jobid"`
}

type SlurmGetResponse struct {
	Output string `json:"output"`
}

// Prepare uncompress the input file and prepares the job for execution
//   - uncompress the input file
//   - create a directory for the job
//   - decompress input in the job directory
func (j *Job) Prepare() error {

	// Create the job directory
	_ = os.MkdirAll(j.Path, 0755)

	// Copy the input file to the job directory
	err := utils.Unzip(j.Input, j.Path)
	if err != nil {
		j.AddMessage("could not unzip file, is it base64 encoded? error: " + err.Error())
		j.UpdateStatus(status.Failed)
		return err
	}

	// Check if run.sh exists
	if _, err := os.Stat(j.Path + "/run.sh"); os.IsNotExist(err) {
		j.AddMessage("run.sh does not exist in the input file")
		j.UpdateStatus(status.Failed)
		return err
	}

	// Make the run.sh file executable
	_ = os.Chmod(j.Path+"/run.sh", 0775)

	j.UpdateStatus(status.Prepared)

	return nil

}

// Run executes the job by running the run.sh script in the job directory
func (j *Job) Run() string {

	glog.Infof("Going into %s and executing run.sh", j.Path)

	// Run the job
	errRun := utils.RunScript(j.Path, "run.sh")

	// Compress the output regardless of the error
	bArr, _ := utils.Zip(j.Path)
	// if err != nil {
	// 	j.UpdateStatus(status.Failure)
	// 	return err
	// }

	// glog.Info(bArr)
	encodedOutput := base64.StdEncoding.EncodeToString(bArr)
	j.AddOutput(encodedOutput)

	// Delete the job directory
	if DEBUG {
		glog.Info("DEBUG is true, not deleting the job directory: ", j.Path)
	} else {
		_ = os.RemoveAll(j.Path)
	}

	if errRun != nil {
		glog.Info(j.ID, " Error running script: ", errRun.Error())
		j.AddMessage("could not finish the job, error: " + errRun.Error())
		j.UpdateStatus(status.Failed)
	} else {
		glog.Info(j.ID, " finished successfully")
		j.AddMessage("job finished successfully")
		j.UpdateStatus(status.Success)
	}

	return j.Status
}

// PostToSlurml posts the job to the SLURML API
func (j *Job) PostToSlurml() string {

	// Get the SLURML API URL
	slurmAPIURL := os.Getenv("SLURML_API_URL")
	if slurmAPIURL == "" {
		glog.Error("SLURML_API_URL is not set")
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	slurmAPIURL = slurmAPIURL + "/api/submit"

	// Get the TOKEN for the SLURML API
	slurmAPIToken := os.Getenv("SLURML_API_TOKEN")
	if slurmAPIToken == "" {
		glog.Error("SLURML_API_TOKEN is not set")
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	// Create the JSON payload
	jsonStr := []byte(
		`{
			"payload": "` + j.Input + `"
		}`)

	req, err := http.NewRequest("POST", slurmAPIURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		glog.Error("could not post the job to the SLURML API: ", err.Error())
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	req.Header.Add("Authorization", slurmAPIToken)

	// Create a new request using http
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		glog.Error("could not post the job to the SLURML API: ", err.Error())
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	defer r.Body.Close()

	// Check the response
	if r.StatusCode != http.StatusCreated {
		glog.Error("could not post the job to the SLURML API: ", r.StatusCode)
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	body, _ := io.ReadAll(r.Body)

	// Unmarshal result
	slurmResponse := SlurmPostResponse{}
	err = json.Unmarshal(body, &slurmResponse)
	if err != nil {
		glog.Error("could not unmarshal the response from the SLURML API: ", err.Error())
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	// Add the SLURML ID to the job
	j.AddSlurmJobid(slurmResponse.Jobid)

	j.UpdateStatus(status.Running)

	return j.Status

}

// GetFromSlurml gets the job from the SLURML API
func (j *Job) GetFromSlurml() string {
	glog.Info("Updating job " + j.ID)

	// Retrieve the job from the SLURML API
	slurmAPIURL := os.Getenv("SLURML_API_URL")
	if slurmAPIURL == "" {
		glog.Error("SLURML_API_URL is not set")
		return status.Failed
	}

	// Get the TOKEN for the SLURML API
	slurmAPIToken := os.Getenv("SLURML_API_TOKEN")
	if slurmAPIToken == "" {
		glog.Error("SLURML_API_TOKEN is not set")
		return status.Failed
	}

	// Make the request
	req, err := http.NewRequest("GET", slurmAPIURL+"/api/download/"+strconv.Itoa(j.SlurmID), nil)
	if err != nil {
		glog.Error("could not get the job from the SLURML API: ", err.Error())
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	req.Header.Add("Authorization", slurmAPIToken)

	// Create a new request using http
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		glog.Error("could not get the job from the SLURML API: ", err.Error())
		j.UpdateStatus(status.Failed)
		return j.Status
	}

	defer r.Body.Close()

	switch r.StatusCode {
	case http.StatusNoContent:
		glog.Info("Job " + j.ID + " is not ready for download yet")
		return j.Status

	case http.StatusInternalServerError:
		glog.Error("Job " + j.ID + " has failed with  status code 500")
		j.UpdateStatus(status.Failed)
		return j.Status

	case http.StatusPartialContent:
		glog.Info("Job " + j.ID + " is ready for download [partial content]")

		slurmReponse, err := GetSlurmResponse(r)
		if err != nil {
			glog.Error("could not unmarshal the response from the SLURML API: ", err.Error())
			j.UpdateStatus(status.Failed)
			return status.Failed
		}

		j.AddOutput(slurmReponse.Output)
		j.UpdateStatus(status.Partial)
		return j.Status

	case http.StatusOK:
		glog.Info("Job " + j.ID + " is ready for download")

		slurmReponse, err := GetSlurmResponse(r)
		if err != nil {
			glog.Error("could not unmarshal the response from the SLURML API: ", err.Error())
			j.UpdateStatus(status.Failed)
			return status.Failed
		}

		j.AddOutput(slurmReponse.Output)
		j.UpdateStatus(status.Success)
		return j.Status

	default:
		glog.Error("could not get the job from the SLURML API: ", r.StatusCode)
		j.UpdateStatus(status.Failed)
		return j.Status
	}
}

// Execute prepares and runs the job
func (j *Job) Execute() error {

	// Prepare the job
	err := j.Prepare()
	if err != nil {
		return err
	}

	// Run the job
	if j.Slurml {
		runStatus := j.PostToSlurml()
		if runStatus != status.Running {
			return errors.New("job failed")
		}
	} else {
		runStatus := j.Run()
		if runStatus == status.Failed {
			return errors.New("job failed")
		}
	}

	return nil
}

// Validate checks if the job is valid
func (j *Job) Validate() error {

	if j.ID == "" {
		return errors.New("job id is required")
	}

	return nil

}

// AddMessage adds a message to the job
func (j *Job) AddMessage(message string) {
	j.Message = message
}

// GetSlurmResponse unmarshals the response from the SLURML API
func GetSlurmResponse(r *http.Response) (SlurmGetResponse, error) {

	slurmReponse := SlurmGetResponse{}
	err := json.NewDecoder(r.Body).Decode(&slurmReponse)
	if err != nil {
		return slurmReponse, err
	}

	return slurmReponse, nil
}
