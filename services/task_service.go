// Package services provides the services for the jobd application
package services

import (
	"jobd/domain/jobs"
	"time"

	"github.com/golang/glog"
)

// RunTasks runs all queued jobs
func RunTasks() error {

	queuedJobs, _ := jobs.ListQueued()
	// if err != nil {
	// 	log.Println(err)
	// 	return nil
	// }

	for _, job := range queuedJobs {
		go func(j jobs.Job) {
			_ = j.Execute()
		}(job)
	}

	return nil
}

// UpdateSlurmlJobs updates all jobs that are running on Slurml
func UpdateSlurmlJobs() error {

	slurmJobs, _ := jobs.ListSlurml()
	// if err != nil {
	// 	log.Println(err)
	// 	return nil
	// }

	for _, job := range slurmJobs {
		go func(j jobs.Job) {
			_ = j.GetFromSlurml()
		}(job)
	}

	return nil
}

// ClearOldJobs clears jobs that are older than 2 days
func ClearOldJobs() error {

	cutoff := time.Now().AddDate(0, 0, -2)

	oldJobs, _ := jobs.ListOld(cutoff)

	for _, job := range oldJobs {
		go func(j jobs.Job) {
			glog.Info("Deleting job ", j.ID, " older than ", cutoff)
			_ = j.Delete()
		}(job)
	}

	return nil
}
