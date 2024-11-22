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

// UploadJob is the controller for uploading a job to the queue
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

// GetJob is the controller for getting a job from the queue
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
