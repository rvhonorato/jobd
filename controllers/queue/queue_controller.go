// Package queue provides the queue for the jobd application
package queue

import (
	"jobd/domain/jobs"
	"jobd/domain/status"
	"jobd/errors"
	"jobd/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// UploadJob godoc
// @Summary Upload a new job to the queue
// @Description Upload a payload. `id` is a unique user-provided job identificator. The `input` field must contain a base64 encoded`.zip` file with a `run.sh` script and the input data. `slurml` marks the job for redirection to the `slurml` endpoint (wip)
// @Accept json
// @Produce json
// @Param job body jobs.Upload true "Job to be uploaded"
// @Success 201 {object} jobs.Job "Job successfully created"
// @Failure 400 {object} errors.RestErr "Bad request - validation error"
// @Failure 500 {object} errors.RestErr "Internal server error"
// @Router /api/upload [post]
func UploadJob(c *gin.Context) {
	var j jobs.Job
	var uploadRequest jobs.Upload
	err := c.BindJSON(&uploadRequest)
	if err != nil {
		glog.Error(err)
		err := errors.NewBadRequestError("error reading job from request " + err.Error())
		c.JSON(err.Status, err)
		return
	}

	j = jobs.Job{
		ID:     uploadRequest.Id,
		Input:  uploadRequest.Input,
		Slurml: uploadRequest.Slurml,
	}

	if err := j.Validate(); err != nil {
		glog.Error(err)
		err := errors.NewBadRequestError("error validating job " + err.Error())
		c.JSON(err.Status, err)
		return
	}

	result, errPost := services.CreateJob(j)
	if errPost != nil {
		glog.Error(errPost)
		c.JSON(errPost.Status, errPost)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// RetrieveJob godoc
// @Summary Retrieve a job from the queue
// @Description Fetches a job by its `id` (provided by the user) with partial content handling
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} jobs.Job "Successfully retrieved job"
// @Success 206 {object} jobs.Job "Partially completed job"
// @Failure 404 {object} errors.RestErr "Job not found"
// @Failure 500 {object} errors.RestErr "Internal server error"
// @Router /api/get/{id} [get]
func RetrieveJob(c *gin.Context) {
	id := c.Param("id")
	j := jobs.Job{ID: id}

	result, err := services.GetJob(j)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	// Do things related to getting the job?
	// Clear the input and path before returning the job
	result.Input = ""
	result.Path = ""

	switch result.Status {
	case status.Partial:
		c.JSON(http.StatusPartialContent, result)
	default:
		c.JSON(http.StatusOK, result)
	}
}
