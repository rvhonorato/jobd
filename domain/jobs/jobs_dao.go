// Package jobs provides the domain object for jobs
// dao = data access object, a pattern for accessing data from a database
package jobs

import (
	"encoding/json"
	"jobd/datasource/db"
	"jobd/domain/status"
	"jobd/errors"
	"time"
)

func (j *Job) Save() *errors.RestErr {
	j.LastUpdated = time.Now()
	// Check if the ID exists in the database already
	err := db.Client.Read(db.NAME, j.ID, &j)
	if err == nil {
		// this id exists already, raise an error
		return errors.NewBadRequestError("job already exists")
	}

	_ = db.Client.Write(db.NAME, j.ID, &j)
	// if err != nil {
	// 	return errors.NewInternalServerError("error saving job to database")
	// }
	return nil
}

func (j *Job) Get() *errors.RestErr {
	err := db.Client.Read(db.NAME, j.ID, &j)
	if err != nil {
		return errors.NewInternalServerError("error getting job from database")
	}
	return nil
}

func (j *Job) Delete() *errors.RestErr {
	_ = db.Client.Delete(db.NAME, j.ID)
	// if err != nil {
	// 	return errors.NewInternalServerError("error deleting job from database")
	// }
	return nil
}

// ListQueued lists all jobs in the database with a status of "queued"
func ListQueued() ([]Job, *errors.RestErr) {

	// Read all records from the database
	records, _ := db.Client.ReadAll(db.NAME)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }

	jobs := []Job{}
	for _, j := range records {
		// Unmarshal the record into a Job
		foundj := Job{}
		_ = json.Unmarshal([]byte(j), &foundj)
		// if err := json.Unmarshal([]byte(j), &foundj); err != nil {
		// fmt.Println("Error", err)
		// err := errors.NewInternalServerError("error getting jobs from database (unmarshal)")
		// return nil, err
		// }

		// Check if the job is queued and has not already been added to the list
		if foundj.Status == status.Queued {
			jobs = append(jobs, foundj)
		}

	}

	return jobs, nil
}

// ListSlurml lists all jobs in the database with a status of "slurm"
func ListSlurml() ([]Job, *errors.RestErr) {

	// Read all records from the database
	records, _ := db.Client.ReadAll(db.NAME)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }

	jobs := []Job{}
	for _, j := range records {
		// Unmarshal the record into a Job
		foundj := Job{}
		_ = json.Unmarshal([]byte(j), &foundj)
		// if err := json.Unmarshal([]byte(j), &foundj); err != nil {
		// fmt.Println("Error", err)
		// err := errors.NewInternalServerError("error getting jobs from database (unmarshal)")
		// return nil, err
		// }

		// Filter the jobs that are slurm and that are still running
		if foundj.Slurml && foundj.Status == status.Running {
			jobs = append(jobs, foundj)
		}

	}

	return jobs, nil

}

// ListOld lists all jobs in the database that are older than the specified time
func ListOld(t time.Time) ([]Job, *errors.RestErr) {

	// Read all from the database, unmarshaling the response.
	records, _ := db.Client.ReadAll(db.NAME)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }

	jobs := []Job{}
	for _, j := range records {
		// Unmarshal the record into a Job
		foundj := Job{}
		_ = json.Unmarshal([]byte(j), &foundj)
		// if err := json.Unmarshal([]byte(j), &foundj); err != nil {
		// fmt.Println("Error", err)
		// err := errors.NewInternalServerError("error getting jobs from database (unmarshal)")
		// return nil, err
		// }

		// Check if the job is older than the cutoff time and has not already been added to the list
		if foundj.LastUpdated.Before(t) {
			jobs = append(jobs, foundj)
		}

	}

	return jobs, nil
}

// UpdateStatus updates the status of the job
func (j *Job) UpdateStatus(s string) *errors.RestErr {
	j.Status = s
	j.LastUpdated = time.Now()
	_ = db.Client.Write(db.NAME, j.ID, &j)
	// if err != nil {
	// 	return errors.NewInternalServerError("error saving job to database")
	// }
	return nil
}

// AddOutput adds output to the job
func (j *Job) AddOutput(o string) *errors.RestErr {
	j.Output = o
	j.LastUpdated = time.Now()
	_ = db.Client.Write(db.NAME, j.ID, &j)
	// if err != nil {
	// 	return errors.NewInternalServerError("error saving job to database")
	// }
	return nil
}

// AddSlurmJobid adds the slurm jobid to the job
func (j *Job) AddSlurmJobid(id int) *errors.RestErr {
	j.SlurmID = id
	j.LastUpdated = time.Now()
	_ = db.Client.Write(db.NAME, j.ID, &j)
	// if err != nil {
	// 	return errors.NewInternalServerError("error saving job to database")
	// }
	return nil
}

// ListOld lists all jobs in the database that are older than the specified time
func (jl *JobList) ListOld(t time.Time) *errors.RestErr {

	// Read all from the database, unmarshaling the response.
	records, _ := db.Client.ReadAll(db.NAME)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }

	for _, j := range records {
		foundj := Job{}
		_ = json.Unmarshal([]byte(j), &foundj)
		// if err := json.Unmarshal([]byte(j), &foundj); err != nil {
		// 	fmt.Println("Error", err)
		// 	err := errors.NewInternalServerError("error getting jobs from database (unmarshal)")
		// 	return err
		// }
		if foundj.LastUpdated.Before(t) {
			jl.Jobs = append(jl.Jobs, foundj)
		}
	}

	return nil
}
