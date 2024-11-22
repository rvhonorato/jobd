// dao = data access object, a pattern for accessing data from a database
package jobs

import (
	"jobd/datasource/db"
	"jobd/domain/status"
	"jobd/errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	_ = db.InitDB()
	gin.SetMode(gin.TestMode)
}

func TestJob_Save(t *testing.T) {

	// Add a job to the database
	j := &Job{ID: "existing-id"}
	_ = db.Client.Write(db.NAME, j.ID, j)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		LastUpdated time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *errors.RestErr
	}{
		{
			name: "TestJob_Save",
			fields: fields{
				ID: "123",
			},
			want: nil,
		},
		{
			name: "TestJob_Save_Existing",
			fields: fields{
				ID: "existing-id",
			},
			want: errors.NewBadRequestError("job already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:          tt.fields.ID,
				Status:      tt.fields.Status,
				Path:        tt.fields.Path,
				Input:       tt.fields.Input,
				Output:      tt.fields.Output,
				LastUpdated: tt.fields.LastUpdated,
			}
			if got := j.Save(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Save() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Get(t *testing.T) {
	// Add a job to the database
	j := &Job{ID: "TestJob_Get"}
	_ = db.Client.Write(db.NAME, j.ID, j)
	// trunk-ignore(golangci-lint/errcheck)
	defer db.Client.Delete(db.NAME, j.ID)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		LastUpdated time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *errors.RestErr
	}{
		{
			name: "TestJob_Get",
			fields: fields{
				ID: "TestJob_Get",
			},
		},
		{
			name: "FailTestJob_Get",
			fields: fields{
				ID: "FailTestJob_Get",
			},
			want: errors.NewInternalServerError("error getting job from database"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:          tt.fields.ID,
				Status:      tt.fields.Status,
				Path:        tt.fields.Path,
				Input:       tt.fields.Input,
				Output:      tt.fields.Output,
				LastUpdated: tt.fields.LastUpdated,
			}
			if got := j.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Delete(t *testing.T) {
	// Add a job to the database
	j := &Job{ID: "to-be-deleted"}
	_ = db.Client.Write(db.NAME, j.ID, j)
	// defer db.Client.Delete(db.NAME, j.ID)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		LastUpdated time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *errors.RestErr
	}{
		{
			name: "TestJob_Delete",
			fields: fields{
				ID: "to-be-deleted",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:          tt.fields.ID,
				Status:      tt.fields.Status,
				Path:        tt.fields.Path,
				Input:       tt.fields.Input,
				Output:      tt.fields.Output,
				LastUpdated: tt.fields.LastUpdated,
			}
			if got := j.Delete(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListQueued(t *testing.T) {
	// Add a job to the database as queued
	j := &Job{ID: "queued-job", Status: status.Queued}
	_ = db.Client.Write(db.NAME, j.ID, j)
	// trunk-ignore(golangci-lint/errcheck)
	defer db.Client.Delete(db.NAME, j.ID)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	tests := []struct {
		name  string
		want  []Job
		want1 *errors.RestErr
	}{
		{
			name: "TestListQueued",
			want: []Job{
				{
					ID:     "queued-job",
					Status: status.Queued,
				},
			},
			want1: nil,
		},
		// {
		// 	name:  "TestListQueued with no queued jobs",
		// 	want:  []Job{},
		// 	want1: nil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ListQueued()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListQueued() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ListQueued() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestJob_UpdateStatus(t *testing.T) {
	// Create a job
	j := &Job{ID: "TestJob_UpdateStatus"}
	_ = db.Client.Write(db.NAME, j.ID, j)
	// trunk-ignore(golangci-lint/errcheck)
	defer db.Client.Delete(db.NAME, j.ID)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		LastUpdated time.Time
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *errors.RestErr
	}{
		{
			name: "TestJob_UpdateStatus",
			fields: fields{
				ID:     "TestJob_UpdateStatus",
				Status: status.Queued,
			},
			args: args{
				s: status.Running,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:          tt.fields.ID,
				Status:      tt.fields.Status,
				Path:        tt.fields.Path,
				Input:       tt.fields.Input,
				Output:      tt.fields.Output,
				LastUpdated: tt.fields.LastUpdated,
			}
			if got := j.UpdateStatus(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobList_ListOld(t *testing.T) {

	// Add a job to the database
	j := &Job{ID: "TestJobList_ListOld", Status: status.Queued, LastUpdated: time.Now().Add(-time.Hour * 24 * 2)}
	_ = db.Client.Write(db.NAME, j.ID, j)
	defer db.Client.Delete(db.NAME, j.ID)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	type fields struct {
		Jobs []Job
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *errors.RestErr
	}{
		{
			name: "TestJobList_ListOld",
			fields: fields{
				Jobs: []Job{
					{
						ID:          "TestJobList_ListOld",
						Status:      status.Queued,
						LastUpdated: time.Now().Add(-time.Hour * 24 * 2),
					},
				},
			},
			args: args{
				t: time.Now(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jl := &JobList{
				Jobs: tt.fields.Jobs,
			}
			if got := jl.ListOld(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobList.ListOld() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_AddOutput(t *testing.T) {
	// Delete the database after the test
	defer os.RemoveAll(db.NAME)
	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		Message     string
		LastUpdated time.Time
	}
	type args struct {
		o string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *errors.RestErr
	}{
		{
			name: "TestJob_AddOutput",
			fields: fields{
				ID:     "TestJob_AddOutput",
				Status: status.Queued,
			},
			args: args{
				o: "base64encodedoutput",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:          tt.fields.ID,
				Status:      tt.fields.Status,
				Path:        tt.fields.Path,
				Input:       tt.fields.Input,
				Output:      tt.fields.Output,
				Message:     tt.fields.Message,
				LastUpdated: tt.fields.LastUpdated,
			}
			if got := j.AddOutput(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.AddOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListOld(t *testing.T) {
	// Add a job to the database
	j := &Job{ID: "TestListOld", Status: status.Queued, LastUpdated: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	_ = db.Client.Write(db.NAME, j.ID, j)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	type args struct {
		t time.Time
	}
	tests := []struct {
		name  string
		args  args
		want  []Job
		want1 *errors.RestErr
	}{
		{
			name: "TestListOld",
			args: args{
				t: time.Now(),
			},
			want: []Job{
				{
					ID:          "TestListOld",
					Status:      status.Queued,
					LastUpdated: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ListOld(tt.args.t)
			// for i, _ := range got {
			// 	got[i].LastUpdated = tt.want[i].LastUpdated
			// }
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListOld() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ListOld() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestListSlurml(t *testing.T) {
	// Add a job to the database that is SLURM
	j := &Job{ID: "slurm-job", Status: status.Running, Slurml: true}
	_ = db.Client.Write(db.NAME, j.ID, j)
	// trunk-ignore(golangci-lint/errcheck)
	defer db.Client.Delete(db.NAME, j.ID)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)
	tests := []struct {
		name  string
		want  []Job
		want1 *errors.RestErr
	}{
		{
			name: "TestListSlurml",
			want: []Job{
				{
					ID:     "slurm-job",
					Status: status.Running,
					Slurml: true,
				},
			},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ListSlurml()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListSlurml() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ListSlurml() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestJob_AddSlurmJobid(t *testing.T) {
	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		Message     string
		Slurm       bool
		SlurmID     int
		LastUpdated time.Time
	}
	type args struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *errors.RestErr
	}{
		{
			name: "TestJob_AddSlurmJobid",
			fields: fields{
				ID:      "TestJob_AddSlurmJobid",
				SlurmID: 0,
			},
			args: args{
				id: 42,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				ID:          tt.fields.ID,
				Status:      tt.fields.Status,
				Path:        tt.fields.Path,
				Input:       tt.fields.Input,
				Output:      tt.fields.Output,
				Message:     tt.fields.Message,
				Slurml:      tt.fields.Slurm,
				SlurmID:     tt.fields.SlurmID,
				LastUpdated: tt.fields.LastUpdated,
			}
			if got := j.AddSlurmJobid(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.AddSlurmJobid() = %v, want %v", got, tt.want)
			}
		})
	}
}
