package validator

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestNewFieldError checks that a FieldError is created correctly.
func TestNewFieldError(t *testing.T) {
	message := "This field is required"
	err := NewFieldError(message)
	if err.Error() != message {
		t.Errorf("expected error message %q, got %q", message, err.Error())
	}
}

// TestValidationErrorsError verifies that adding errors via Field()
// produces an aggregated error string containing the appropriate messages.
func TestValidationErrorsError(t *testing.T) {
	ve := New()
	if ve == nil {
		t.Error("New() should return a non-nil ValidationErrors")
	}

	ve.Field("Name", NewFieldError("Name cannot be empty"))
	ve.Field("Email", NewFieldError("Invalid email address"))

	// Since map iteration order is not guaranteed,
	// we simply check that each expected substring is present.
	errStr := ve.Error()
	if !strings.Contains(errStr, "Name: Name cannot be empty") {
		t.Errorf("expected error string to contain %q, got %q", "Name: Name cannot be empty", errStr)
	}
	if !strings.Contains(errStr, "Email: Invalid email address") {
		t.Errorf("expected error string to contain %q, got %q", "Email: Invalid email address", errStr)
	}
}

// TestMarshalJSON ensures that MarshalJSON correctly converts ValidationErrors into JSON.
func TestMarshalJSON(t *testing.T) {
	ve := New()
	ve.Field("Username", NewFieldError("Username is required"))
	ve.Field("Password", NewFieldError("Password must be at least 8 characters"))

	data, err := json.Marshal(ve)
	if err != nil {
		t.Errorf("unexpected error during MarshalJSON: %v", err)
	}

	// Unmarshal back into a map[string]string to verify the JSON structure.
	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Errorf("error unmarshalling JSON: %v", err)
	}

	if result["Username"] != "Username is required" {
		t.Errorf("expected Username error %q, got %q", "Username is required", result["Username"])
	}
	if result["Password"] != "Password must be at least 8 characters" {
		t.Errorf("expected Password error %q, got %q", "Password must be at least 8 characters", result["Password"])
	}
}
