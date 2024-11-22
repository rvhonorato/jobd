// Package utils provides the utility functions for the jobd application
package utils

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golang/glog"
)

// uniqueID generates a unique ID for a Job
func UniqueID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if n <= 0 {
		n = 8
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Unzip takes a zip file encoded as base64 string and uncompresses it to a destination directory
func Unzip(source, destination string) error {

	// Decode the base64 string
	inp, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		glog.Info(err)
		return err
	}

	reader, _ := zip.NewReader(bytes.NewReader(inp), int64(len(inp)))

	// Thanks `https://gist.github.com/paulerickson/6d8650947ee4e3f3dbcc28fde10eaae7` (:
	for _, file := range reader.File {
		reader, _ := file.Open()
		// if err != nil {
		// 	return err
		// }
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		_ = os.Remove(path)
		// if err != nil {
		// 	return err
		// }
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		// if err != nil {
		// 	return err
		// }
		defer writer.Close()
		_, _ = io.Copy(writer, reader)
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}

// Zip compresses the a directory
func Zip(srcDir string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Check if the source directory exists
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return nil, err
	}

	_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// if err != nil {
		// 	glog.Error("failed to access path", path, err)
		// 	return fmt.Errorf("failed to access path %q: %w", path, err)
		// }

		relPath, _ := filepath.Rel(srcDir, path)
		// if err != nil {
		// 	glog.Error("failed to get relative path:", err)
		// 	return fmt.Errorf("failed to get relative path: %w", err)
		// }

		// glog.Info("Adding file to zip: ", relPath)
		zipFile, _ := zipWriter.Create(relPath)
		// if err != nil {
		// 	glog.Error("failed to create zip file for", relPath, ":", err)
		// 	return fmt.Errorf("failed to create zip file for %q: %w", relPath, err)
		// }

		file, _ := os.Open(path)
		// if err != nil {
		// 	glog.Error("failed to open file", path, ":", err)
		// 	return fmt.Errorf("failed to open file %q: %w", path, err)
		// }
		defer file.Close()

		// glog.Info("Writing file to zip: ", file)
		_, _ = io.Copy(zipFile, file)
		// if err != nil {
		// 	glog.Error("failed to write to zip file: %w", err)
		// 	return fmt.Errorf("failed to write to zip file: %w", err)
		// }

		return nil
	})

	zipWriter.Close()

	zipBytes := buf.Bytes()

	return zipBytes, nil
}

// RunScript runs a script in a directory
func RunScript(dir, script string) error {

	var out bytes.Buffer

	// Run the job
	cmd := exec.Command("./" + script)
	cmd.Dir = dir
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
