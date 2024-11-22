package services

import (
	"jobd/datasource/db"
	"jobd/domain/jobs"
	"jobd/domain/status"
	"jobd/errors"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	_ = db.InitDB()
	gin.SetMode(gin.TestMode)
}

func TestGetJob(t *testing.T) {
	// Add jobs to the database

	// Successful
	successJ := &jobs.Job{ID: "TestGetJobSuccess", Status: status.Success}
	_ = db.Client.Write(db.NAME, successJ.ID, successJ)

	// Running
	runningJ := &jobs.Job{ID: "TestGetJobRunning", Status: status.Running}
	_ = db.Client.Write(db.NAME, runningJ.ID, runningJ)

	// Partial
	partialJ := &jobs.Job{ID: "TestGetJobPartial", Status: status.Partial}
	_ = db.Client.Write(db.NAME, partialJ.ID, partialJ)

	// Remove the database after the test
	defer os.RemoveAll(db.NAME)

	type args struct {
		j jobs.Job
	}
	tests := []struct {
		name  string
		args  args
		want  *jobs.Job
		want1 *errors.RestErr
	}{
		{
			name: "GetJobSuccess",
			args: args{
				j: jobs.Job{
					ID: "TestGetJobSuccess",
				},
			},
			want:  successJ,
			want1: nil,
		},
		{
			name: "GetJobNonExisting",
			args: args{
				j: jobs.Job{
					ID: uuid.New().String(),
				},
			},
			want:  nil,
			want1: errors.NewNotFoundError("job not found"),
		},
		{
			name: "GetJobRunning",
			args: args{
				j: jobs.Job{
					ID: "TestGetJobRunning",
				},
			},
			want:  nil,
			want1: errors.NewStatusAccepted("job not ready"),
		},
		{
			name: "GetJobPartial",
			args: args{
				j: jobs.Job{
					ID: "TestGetJobPartial",
				},
			},
			want:  partialJ,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetJob(tt.args.j)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJob() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetJob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

// func TestClearOldJobs(t *testing.T) {
// 	// Remove the database after the test
// 	defer os.RemoveAll(db.NAME)

// 	type args struct {
// 		jobs *jobs.JobList
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *errors.RestErr
// 	}{
// 		{
// 			name: "ClearOldJobs",
// 			args: args{
// 				jobs: &jobs.JobList{
// 					Jobs: []jobs.Job{
// 						{
// 							ID:     "TestClearOldJobs",
// 							Status: "completed",
// 						},
// 					},
// 				},
// 			},
// 			want: nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := ClearOldJobs(tt.args.jobs); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ClearOldJobs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestCreateJob(t *testing.T) {
	// Add a job to the database
	j := &jobs.Job{ID: "existing-job-test-create-job"}
	_ = db.Client.Write(db.NAME, j.ID, j)

	defer os.RemoveAll(db.NAME)

	type args struct {
		j jobs.Job
	}
	tests := []struct {
		name  string
		args  args
		want  *jobs.Job
		want1 *errors.RestErr
	}{

		{
			name: "CreateJob",
			args: args{
				j: jobs.Job{
					ID: "TestCreateJob",
				},
			},
			want: &jobs.Job{
				ID:     "TestCreateJob",
				Status: "QUEUED",
				Path:   DATAPATH + "/TestCreateJob",
			},
			want1: nil,
		},
		{
			name: "FailCreateJob",
			args: args{
				j: jobs.Job{
					ID: "existing-job-test-create-job",
				},
			},
			want:  nil,
			want1: errors.NewBadRequestError("job already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CreateJob(tt.args.j)
			if got != nil {
				// Warning: This is a hack to get the test to pass
				//   Overwrite the time with the expected time
				got.LastUpdated = tt.want.LastUpdated
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateJob() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreateJob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
