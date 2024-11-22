package services

import (
	"jobd/datasource/db"
	"jobd/domain/jobs"
	"jobd/domain/status"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	_ = db.InitDB()
	gin.SetMode(gin.TestMode)
}

func TestRunTasks(t *testing.T) {
	// Add some queued jobs to the database
	j := &jobs.Job{ID: "TestRunTasks", Status: status.Queued, Input: "UEsDBAoAAAAAAKVdUFYAAAAAAAAAAAAAAAAGABwAcnVuLnNoVVQJAAM1Ce5jNQnuY3V4CwABBPUBAAAEAAAAAFBLAQIeAwoAAAAAAKVdUFYAAAAAAAAAAAAAAAAGABgAAAAAAAAAAACkgQAAAABydW4uc2hVVAUAAzUJ7mN1eAsAAQT1AQAABAAAAABQSwUGAAAAAAEAAQBMAAAAQAAAAAAA"}
	// _ = j.Save()
	// j := &jobs.Job{ID: "TestRunTasks", Status: status.Queued}
	_ = db.Client.Write(db.NAME, j.ID, j)

	// Remove the testing file
	defer os.RemoveAll("run.sh")

	// // Remove the database after the test
	defer os.RemoveAll(db.NAME)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "run tasks",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunTasks(); (err != nil) != tt.wantErr {
				t.Errorf("RunTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClearOldJobs(t *testing.T) {
	// Add some queued jobs to the database
	j := &jobs.Job{ID: "TestClearOldJobs", LastUpdated: time.Now().AddDate(0, 0, -99)}
	_ = db.Client.Write(db.NAME, j.ID, j)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "clear old jobs",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ClearOldJobs(); (err != nil) != tt.wantErr {
				t.Errorf("ClearOldJobs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
