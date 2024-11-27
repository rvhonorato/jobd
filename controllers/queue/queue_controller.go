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
// @Description Validates and creates a new job in the system
// @Tags Jobs
// @Accept json
// @Produce json
// @Param job body jobs.Job true "Job to be uploaded"
// @Success 201 {object} jobs.Job "Job successfully created"
// @Failure 400 {object} errors.RestErr "Bad request - validation error"
// @Failure 500 {object} errors.RestErr "Internal server error"
// @Router /api/upload [post]
func UploadJob(c *gin.Context) {
	var j jobs.Job

	err := c.BindJSON(&j)
	if err != nil {
		glog.Error(err)
		err := errors.NewBadRequestError("error reading job from request " + err.Error())
		c.JSON(err.Status, err)
		return
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
// @Description Fetches a job by its ID with partial content handling
// @Tags Jobs
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
