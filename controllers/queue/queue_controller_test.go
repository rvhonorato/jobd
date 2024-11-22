package queue

import (
	"bytes"
	"jobd/datasource/db"
	"jobd/domain/jobs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	_ = db.InitDB()
	// _ = migrations.Migrate()
	gin.SetMode(gin.TestMode)
}

func TestUploadJob(t *testing.T) {

	// --------------------------------------------------
	// Set things for the test here
	var req *http.Request
	var err error
	var recorder *httptest.ResponseRecorder
	var jsonData []byte

	// Initialize a new gin router
	router := gin.Default()

	// Add the UploadJob endpoint to the router
	router.POST("/upload", UploadJob)

	// Delete the database after the test
	defer os.RemoveAll(db.NAME)

	// --------------------------------------------------
	// Test 1 - Pass the test
	jsonData = []byte(`{
		"id": "test1",
		"file": "base64encodedfile"
	}`)

	req, err = http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error creating request: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Create a new recorder to capture the response from the endpoint
	recorder = httptest.NewRecorder()

	// Call the UploadJob endpoint
	router.ServeHTTP(recorder, req)

	// Check that the status code is 201 (created)
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, recorder.Code)
	}

	// --------------------------------------------------
	// Test 2 - Fail the test with invalid json
	jsonData = []byte(`{
		"id": "test2",
	}`)

	req, err = http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonData))

	if err != nil {
		t.Errorf("Error creating request: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Create a new recorder to capture the response from the endpoint
	recorder = httptest.NewRecorder()

	// Call the UploadJob endpoint
	router.ServeHTTP(recorder, req)

	// Check that the status code is 400 (bad request)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	// --------------------------------------------------
	// Test 3 - Fail by passing an id that already exists

	// create a job in the database
	j := &jobs.Job{ID: "test3"}
	_ = db.Client.Write(db.NAME, j.ID, j)

	jsonData = []byte(`{
		"id": "test3",
		"file": "base64encodedfile"
	}`)

	req, err = http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonData))

	if err != nil {
		t.Errorf("Error creating request: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Create a new recorder to capture the response from the endpoint
	recorder = httptest.NewRecorder()

	// Call the UploadJob endpoint
	router.ServeHTTP(recorder, req)

	// Check that the status code is 400 (bad request)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	// --------------------------------------------------
	// Test 4 - Fail by passing a json with invalid fields
	jsonData = []byte(`{
		}`)

	req, err = http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonData))

	if err != nil {
		t.Errorf("Error creating request: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Create a new recorder to capture the response from the endpoint
	recorder = httptest.NewRecorder()

	// Call the UploadJob endpoint
	router.ServeHTTP(recorder, req)

	// Check that the status code is 400 (bad request)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}

}

func TestRetrieveJob(t *testing.T) {

	// Create a job in the database and get its ID
	j := &jobs.Job{ID: "TestRetrieveJob"}
	_ = db.Client.Write(db.NAME, j.ID, j)
	defer os.RemoveAll(db.NAME)

	// --------------------------------------------------

	router := gin.Default()

	router.GET("/retrieve/:id", RetrieveJob)

	// Pass the test

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/retrieve/"+j.ID, nil)

	router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected status code %d, got %d", http.StatusAccepted, w.Result().StatusCode)
	}

	// Fail the test

	w = httptest.NewRecorder()

	req = httptest.NewRequest("GET", "/retrieve/does-not-exist", nil)

	router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Result().StatusCode)
	}

}
