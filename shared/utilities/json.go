package utilities

import (
	"encoding/json"
	"net/http"
)

// ReadJSON reads JSON from the request body and decodes it into the provided value.
func ReadJSON(w http.ResponseWriter, r *http.Request, value any) error {
	maxBytes := int64(1_048_576) // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(value)
}

// WriteJSON writes the provided value as JSON to the response writer with the given status code.
func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(value)
}
