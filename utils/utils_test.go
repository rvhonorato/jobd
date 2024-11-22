package utils

import (
	"math/rand"
	"os"
	"reflect"
	"testing"
)

func TestUniqueID(t *testing.T) {
	// Seed the random number generator with a fixed value so that the test is deterministic
	rand.Seed(42)

	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "generate a unique ID with n=10",
			args: args{10},
			want: "HRukpTTueZ",
		},
		{
			name: "generate a unique ID with n=20",
			args: args{20},
			want: "PtNeuvunhuksqVGzAdxl",
		},
		{
			name: "generate a unique ID with n=0",
			args: args{0},
			want: "gghEjkMV",
		},
		{
			name: "generate a unique ID with n=-1",
			args: args{-1},
			want: "eZJpmKqa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueID(tt.args.n); got != tt.want {
				t.Errorf("UniqueID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnzip(t *testing.T) {

	tempDir, err := os.MkdirTemp("/tmp", "jobd-test-")
	if err != nil {
		t.Errorf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	validBase64Zip := "UEsDBBQAAAAIAIRTSVZwXkK2oQAAANkAAAAGABwAcnVuLnNoVVQJAAOovORj0LzkY3V4CwABBPYBAAAEFAAAAFWOTQrCMBSE9znF2B9oKNL2AC4EXXRRheIF0ibSSJuWvAh6e5OgC+FthjfzzaS7atCmGgRNjJyw7pAVUjiFMifOUtwmTfAnQGLZZoXHOoBGqzeH+2qDlMHWdufeR4uiP15O1w45mpqjRMM5o1mpzWvvU+O0ImkhtcTyjjBtkMV4QkwZ+ddvn8bpRQWwf2GPONEjfxwrTJyRfZ2B8dIONfsAUEsDBAoAAAAAANtaSVYAAAAAAAAAAAAAAAAJABwAdGVtcC1kaXIvVVQJAAN+yeRjf8nkY3V4CwABBPYBAAAEFAAAAFBLAwQKAAAAAADhWklWAAAAAAAAAAAAAAAAEwAcAHRlbXAtZGlyL3RlbXAtZGlyMi9VVAkAA4XJ5GOGyeRjdXgLAAEE9gEAAAQUAAAAUEsDBAoAAAAAAOFaSVYAAAAAAAAAAAAAAAAYABwAdGVtcC1kaXIvdGVtcC1kaXIyL2ZpbGUyVVQJAAOFyeRjhcnkY3V4CwABBPYBAAAEFAAAAFBLAwQKAAAAAADYWklWAAAAAAAAAAAAAAAADgAcAHRlbXAtZGlyL2ZpbGUxVVQJAAN4yeRjeMnkY3V4CwABBPYBAAAEFAAAAFBLAQIeAxQAAAAIAIRTSVZwXkK2oQAAANkAAAAGABgAAAAAAAEAAACkgQAAAABydW4uc2hVVAUAA6i85GN1eAsAAQT2AQAABBQAAABQSwECHgMKAAAAAADbWklWAAAAAAAAAAAAAAAACQAYAAAAAAAAABAA7UHhAAAAdGVtcC1kaXIvVVQFAAN+yeRjdXgLAAEE9gEAAAQUAAAAUEsBAh4DCgAAAAAA4VpJVgAAAAAAAAAAAAAAABMAGAAAAAAAAAAQAO1BJAEAAHRlbXAtZGlyL3RlbXAtZGlyMi9VVAUAA4XJ5GN1eAsAAQT2AQAABBQAAABQSwECHgMKAAAAAADhWklWAAAAAAAAAAAAAAAAGAAYAAAAAAAAAAAApIFxAQAAdGVtcC1kaXIvdGVtcC1kaXIyL2ZpbGUyVVQFAAOFyeRjdXgLAAEE9gEAAAQUAAAAUEsBAh4DCgAAAAAA2FpJVgAAAAAAAAAAAAAAAA4AGAAAAAAAAAAAAKSBwwEAAHRlbXAtZGlyL2ZpbGUxVVQFAAN4yeRjdXgLAAEE9gEAAAQUAAAAUEsFBgAAAAAFAAUApgEAAAsCAAAAAA=="

	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "unzip a file",
			args: args{
				source:      validBase64Zip,
				destination: tempDir,
			},
			wantErr: false,
		},
		{
			name: "unzip a file with invalid base64",
			args: args{
				source:      "not-base64",
				destination: tempDir,
			},
			wantErr: true,
		},
		{
			name: "unzip a file to an invalid destination",
			args: args{
				source:      validBase64Zip,
				destination: "/invalid/destination",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unzip(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("Unzip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZip(t *testing.T) {
	// Create a temporary directory to zip
	// tempDir, err := os.MkdirTemp("/tmp", "jobd-test-")
	testDir := "/tmp/test-zip"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Errorf("Error creating temporary directory: %v", err)
	}

	// Create a file to zip
	file, err := os.Create(testDir + "/test-file")
	if err != nil {
		t.Errorf("Error creating temporary file: %v", err)
	}
	defer file.Close()

	// Write some content to the file
	_, err = file.WriteString("test")
	if err != nil {
		t.Errorf("Error writing to temporary file: %v", err)
	}
	// Create a subdirectory to zip
	err = os.MkdirAll(testDir+"/subdir", 0755)
	if err != nil {
		t.Errorf("Error creating temporary subdirectory: %v", err)
	}

	// Create a file to zip in the subdirectory
	subFile, err := os.Create(testDir + "/subdir/subfile")
	if err != nil {
		t.Errorf("Error creating temporary subfile: %v", err)
	}
	defer subFile.Close()
	type args struct {
		srcDir string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "zip-dir",
			args: args{
				srcDir: testDir,
			},
			want:    []byte{80, 75, 3, 4, 20, 0, 8, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 14, 0, 0, 0, 115, 117, 98, 100, 105, 114, 47, 115, 117, 98, 102, 105, 108, 101, 1, 0, 0, 255, 255, 80, 75, 7, 8, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 80, 75, 3, 4, 20, 0, 8, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 0, 0, 0, 116, 101, 115, 116, 45, 102, 105, 108, 101, 42, 73, 45, 46, 1, 4, 0, 0, 255, 255, 80, 75, 7, 8, 12, 126, 127, 216, 10, 0, 0, 0, 4, 0, 0, 0, 80, 75, 1, 2, 20, 0, 20, 0, 8, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 117, 98, 100, 105, 114, 47, 115, 117, 98, 102, 105, 108, 101, 80, 75, 1, 2, 20, 0, 20, 0, 8, 0, 8, 0, 0, 0, 0, 0, 12, 126, 127, 216, 10, 0, 0, 0, 4, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 65, 0, 0, 0, 116, 101, 115, 116, 45, 102, 105, 108, 101, 80, 75, 5, 6, 0, 0, 0, 0, 2, 0, 2, 0, 115, 0, 0, 0, 130, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "zip-dir-not-exists",
			args: args{
				srcDir: "/tmp/does-not-exist",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Zip(tt.args.srcDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Zip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Zip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunScript(t *testing.T) {
	// Write a script that prints "hello world"
	cmd := []byte("#!/bin/bash\necho \"hello\"")
	err := os.WriteFile("./script.sh", cmd, 0775)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove("./script.sh")

	script := `#!/bin/bash
echo "10"`
	// Create a temporary directory
	// _, _ = os.MkdirTemp("/tmp", "jobd-test")
	_ = os.MkdirAll("/tmp/jobd-test", 0755)
	// Write the script to the temporary directory
	_ = os.WriteFile("/tmp/jobd-test/script.sh", []byte(script), 0755)
	defer os.RemoveAll("/tmp/jobd-test")

	type args struct {
		dir    string
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "run-script",
			args: args{
				dir:    "./",
				script: "script.sh",
			},
			wantErr: false,
		},
		{
			name: "run-script-unexisting",
			args: args{
				dir:    "./",
				script: "unexisting.sh",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunScript(tt.args.dir, tt.args.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
