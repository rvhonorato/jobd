// Package services provides the services for the jobd application
package services

import (
	"jobd/domain/jobs"
	"jobd/domain/status"
	"jobd/errors"
	"os"

	"github.com/golang/glog"
)

// const DATAPATH = "/data" // TODO: make this configurable
var DATAPATH = os.Getenv("DATAPATH")

func init() {
	if DATAPATH == "" {
		glog.Warning("DATAPATH not set, using default `./data`")
		DATAPATH = "./data"
	}
}

// GetJob gets a job from the database
func GetJob(j jobs.Job) (*jobs.Job, *errors.RestErr) {

	result := &jobs.Job{ID: j.ID}
	err := result.Get()
	if err != nil {
		err := errors.NewNotFoundError("job not found")
		return nil, err
	}

	// If the status is success or failed, return the job
	// Else return a 202 Accepted
	validStatus := []string{status.Success, status.Failed, status.Partial}
	for _, s := range validStatus {
		if result.Status == s {
			return result, nil
		}
	}

	return nil, errors.NewStatusAccepted("job not ready")

}

// PostJob posts a job to the database and save it to the server
// func CreateJob(b []byte, id string) (*jobs.Job, *errors.RestErr) {
func CreateJob(j jobs.Job) (*jobs.Job, *errors.RestErr) {

	j.Path = DATAPATH + "/" + j.ID
	j.Status = status.Queued

	err := j.Save()
	if err != nil {
		return nil, err
	}

	// do something?

	return &j, nil

}
