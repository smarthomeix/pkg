package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarthomeix/pkg/validator"
)

// TestHandleStatus verifies that HandleStatus sets the proper status code.
func TestHandleStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	HandleStatus(rec, http.StatusCreated)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

// TestHandleNotFound verifies that HandleNotFound sets a 404 status code.
func TestHandleNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	HandleNotFound(rec)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// TestHandleJSON verifies that HandleJSON sets the JSON header and encodes the body.
func TestHandleJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	body := map[string]string{"hello": "world"}
	HandleJSON(rec, body)

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if result["hello"] != "world" {
		t.Errorf("expected key 'hello' with value 'world', got %q", result["hello"])
	}
}

// TestHandleJSONWithStatus verifies that HandleJSONWithStatus sets the status code and JSON header, and encodes the body.
func TestHandleJSONWithStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	body := map[string]string{"foo": "bar"}
	status := http.StatusAccepted
	HandleJSONWithStatus(rec, body, status)

	if rec.Code != status {
		t.Errorf("expected status %d, got %d", status, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if result["foo"] != "bar" {
		t.Errorf("expected key 'foo' with value 'bar', got %q", result["foo"])
	}
}

// TestHandleServerError verifies that HandleServerError logs the error and sets a 500 status code.
func TestHandleServerError(t *testing.T) {
	rec := httptest.NewRecorder()
	errInput := errors.New("something went wrong")
	HandleServerError(rec, errInput)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	// Note: We are not capturing log output in this test.
}

// TestHandleValidationError_WithValidationError tests that a validation error returns a 400 status with a JSON body.
func TestHandleValidationError_WithValidationError(t *testing.T) {
	rec := httptest.NewRecorder()

	// Create a validation error using the validator package.
	ve := validator.New()
	ve.Field("Name", validator.NewFieldError("Name is required"))

	HandleValidationError(rec, ve)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	// Decode the JSON response and verify its content.
	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if result["Name"] != "Name is required" {
		t.Errorf("expected field 'Name' error message 'Name is required', got %q", result["Name"])
	}
}

// TestHandleValidationError_WithNonValidationError tests that a non-validation error results in a 500 status.
func TestHandleValidationError_WithNonValidationError(t *testing.T) {
	rec := httptest.NewRecorder()
	nonValErr := errors.New("regular error")
	HandleValidationError(rec, nonValErr)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d for non-validation error, got %d", http.StatusInternalServerError, rec.Code)
	}
}
