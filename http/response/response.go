package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/smarthomeix/pkg/validator"
)

// HandleStatus writes the given status code to the response.
func HandleStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// HandleNotFound writes a 404 Not Found status.
func HandleNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

// HandleJSON writes the given body as JSON to the response.
func HandleJSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

// HandleJSONWithStatus writes the given body as JSON to the response with a specified status code.
func HandleJSONWithStatus(w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

// HandleServerError logs the error and writes a 500 Internal Server Error status.
func HandleServerError(w http.ResponseWriter, err error) {
	log.Println(err)

	HandleStatus(w, http.StatusInternalServerError)
}

// HandleValidationError checks if the error is a validation error.
// If so, it returns a JSON response with a 400 Bad Request status.
// Otherwise, it falls back to handling a server error.
func HandleValidationError(w http.ResponseWriter, err error) {
	var validationErrors validator.ValidationErrors

	if !errors.As(err, &validationErrors) {
		HandleServerError(w, err)
		return // early return to prevent writing a second response
	}

	HandleJSONWithStatus(w, err, http.StatusBadRequest)
}
