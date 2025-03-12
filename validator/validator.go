package validator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldError represents an error associated with a specific field.
type FieldError struct {
	Message string
}

// Error implements the error interface for FieldError.
func (fe FieldError) Error() string {
	return fe.Message
}

// NewFieldError creates a new FieldError with the given message.
func NewFieldError(message string) error {
	return FieldError{
		Message: message,
	}
}

// ValidationErrors holds a collection of errors keyed by field names.
type ValidationErrors map[string]error

// New creates and returns a new ValidationErrors instance.
func New() ValidationErrors {
	return make(ValidationErrors)
}

// Field adds an error for a given field.
func (ve ValidationErrors) Field(name string, err error) {
	ve[name] = err
}

// Error implements the error interface for ValidationErrors.
func (ve ValidationErrors) Error() string {
	var sb strings.Builder

	for field, err := range ve {
		sb.WriteString(fmt.Sprintf("%s: %s\n", field, err.Error()))
	}
	return sb.String()
}

// MarshalJSON customizes JSON marshalling for ValidationErrors.
func (ve ValidationErrors) MarshalJSON() ([]byte, error) {
	errorMessages := make(map[string]string, len(ve))

	for field, err := range ve {
		errorMessages[field] = err.Error()
	}

	return json.Marshal(errorMessages)
}
