// dto - data transfer object; a pattern for transferring data between processes
package jobs

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"jobd/datasource/db"
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

func TestJob_Prepare(t *testing.T) {

	testPath := "./test"
	testPath2 := "./test2"
	defer os.RemoveAll(testPath)
	defer os.RemoveAll(testPath2)

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
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestJob_Prepare",
			fields: fields{
				ID:          "123",
				Status:      "pending",
				Path:        testPath,
				Input:       "UEsDBAoAAAAAAKVdUFYAAAAAAAAAAAAAAAAGABwAcnVuLnNoVVQJAAM1Ce5jNQnuY3V4CwABBPUBAAAEAAAAAFBLAQIeAwoAAAAAAKVdUFYAAAAAAAAAAAAAAAAGABgAAAAAAAAAAACkgQAAAABydW4uc2hVVAUAAzUJ7mN1eAsAAQT1AQAABAAAAABQSwUGAAAAAAEAAQBMAAAAQAAAAAAA",
				Output:      "",
				LastUpdated: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "TestJob_Prepare with invalid base64",
			fields: fields{
				ID:          "123",
				Status:      "pending",
				Path:        testPath,
				Input:       "not-base64",
				Output:      "",
				LastUpdated: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "TestJob_Prepare without run.sh",
			fields: fields{
				ID:          "123",
				Status:      "pending",
				Path:        testPath2,
				Input:       "UEsDBAoAAAAAAEZeUFYAAAAAAAAAAAAAAAAHABwAbm90LXJ1blVUCQADYwruY2MK7mN1eAsAAQT1AQAABAAAAABQSwECHgMKAAAAAABGXlBWAAAAAAAAAAAAAAAABwAYAAAAAAAAAAAApIEAAAAAbm90LXJ1blVUBQADYwruY3V4CwABBPUBAAAEAAAAAFBLBQYAAAAAAQABAE0AAABBAAAAAAA=",
				Output:      "",
				LastUpdated: time.Now(),
			},
			wantErr: true,
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
			if err := j.Prepare(); (err != nil) != tt.wantErr {
				t.Errorf("Job.Prepare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Run(t *testing.T) {

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	// Create a path and an executable run.sh file
	testDir := "./test-run"
	_ = os.Mkdir(testDir, 0755)
	defer os.RemoveAll(testDir)

	testWoRunDir := "./test-wo-run"
	_ = os.Mkdir(testWoRunDir, 0755)
	defer os.RemoveAll(testWoRunDir)

	d1 := []byte("#!/bin/bash\necho \"hello\"")
	err := os.WriteFile("./test-run/run.sh", d1, 0775)
	if err != nil {
		t.Errorf("Job.Run() error = %v", err)
	}

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
		want   string
	}{
		{
			name: "TestJob_Run",
			fields: fields{
				ID:          "123",
				Status:      "pending",
				Path:        testDir,
				Input:       "",
				Output:      "",
				LastUpdated: time.Now(),
			},
			want: status.Success,
		},
		{
			name: "TestJob_Run without run.sh",
			fields: fields{
				ID:          "123",
				Status:      "pending",
				Path:        testWoRunDir,
				Input:       "",
				Output:      "",
				LastUpdated: time.Now(),
			},
			want: status.Failed,
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
			got := j.Run()
			if got != tt.want {
				t.Errorf("Job.Run() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestJob_Execute(t *testing.T) {

	_ = os.MkdirAll("./test", 0755)
	defer os.RemoveAll("./test")

	_ = os.MkdirAll("./test2", 0755)
	defer os.RemoveAll("./test2")

	// Create a zip containing a run.sh file

	// Delete the database and file after the test
	defer os.RemoveAll(db.NAME)

	// Create a zip containing a run.sh file
	var buf bytes.Buffer
	var buf2 bytes.Buffer

	// Create a valid run.sh file
	d1 := []byte("#!/bin/bash\necho \"hello\"")
	err := os.WriteFile("test/run.sh", d1, 0775)
	if err != nil {
		t.Errorf("Job.Run() error = %v", err)
	}

	f1, err := os.Open("test/run.sh")
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	zipWriter := zip.NewWriter(&buf)

	w1, err := zipWriter.Create("test/run.sh")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(w1, f1); err != nil {
		panic(err)
	}

	zipWriter.Close()
	zipBytes := buf.Bytes()
	inpB64 := base64.StdEncoding.EncodeToString(zipBytes)

	// -----

	// Create an invalid run.sh file
	d2 := []byte("!")
	err = os.WriteFile("test2/run.sh", d2, 0775)
	if err != nil {
		t.Errorf("Job.Run() error = %v", err)
	}

	f2, err := os.Open("test2/run.sh")
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	zipWriter2 := zip.NewWriter(&buf2)

	w2, err := zipWriter.Create("run.sh")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(w2, f2); err != nil {
		panic(err)
	}

	zipWriter2.Close()
	invalidzipBytes := buf2.Bytes()
	invalidinpB64 := base64.StdEncoding.EncodeToString(invalidzipBytes)

	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		LastUpdated time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestJob_Execute",
			fields: fields{
				ID:          "TestJob_Execute",
				Status:      "pending",
				Path:        "./test",
				Input:       inpB64,
				Output:      "",
				LastUpdated: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "TestJob_Execute with invalid input for prepare",
			fields: fields{
				ID:          "TestJob_Execute-invalid-prepare",
				Status:      "pending",
				Path:        "./test",
				Input:       "invalid",
				Output:      "",
				LastUpdated: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "TestJob_Execute with invalid input for run",
			fields: fields{
				ID:          "TestJob_Execute-invalid-run",
				Status:      "pending",
				Path:        "./test2",
				Input:       invalidinpB64,
				Output:      "",
				LastUpdated: time.Now(),
			},
			wantErr: true,
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
			if err := j.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Job.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Validate(t *testing.T) {
	type fields struct {
		ID          string
		Status      string
		Path        string
		Input       string
		Output      string
		LastUpdated time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestJob_Validate",
			fields: fields{
				ID: "TestJob_Validate",
			},
			wantErr: false,
		},
		{
			name: "TestJob_Validate with invalid ID",
			fields: fields{
				ID: "",
			},
			wantErr: true,
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
			if err := j.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Job.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_AddMessage(t *testing.T) {
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
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestJob_AddMessage",
			fields: fields{
				Message: "",
			},
			args: args{
				message: "TestJob_AddMessage",
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
				Message:     tt.fields.Message,
				LastUpdated: tt.fields.LastUpdated,
			}
			j.AddMessage(tt.args.message)
			if j.Message != tt.args.message {
				t.Errorf("Job.AddMessage() error = %v, wantErr %v", j.Message, tt.args.message)
			}
		})
	}
}
