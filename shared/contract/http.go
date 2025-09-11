package contract

import "time"

// APIResponse represents the response structure for the API.
type APIResponse struct {
	Data      any       `json:"data,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// APIError represents the error structure for the API.
type APIError struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Details []APIValidationError `json:"details,omitempty"`
}

// APIValidationError represents a validation error for the API.
type APIValidationError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
	Value   any    `json:"value,omitempty"`
}

const (
	ErrorCodeValidation   = "VALIDATION_ERROR"
	ErrorCodeNotFound     = "NOT_FOUND"
	ErrorCodeUnauthorized = "UNAUTHORIZED"
	ErrorCodeForbidden    = "FORBIDDEN"
	ErrorCodeInternal     = "INTERNAL_ERROR"
	ErrorCodeBadRequest   = "BAD_REQUEST"
	ErrorCodeConflict     = "CONFLICT"
	ErrorCodeRateLimit    = "RATE_LIMIT_EXCEEDED"
)

// NewSuccessResponse creates a new success response with the given data.
func NewSuccessResponse(data any) APIResponse {
	return APIResponse{
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse creates a new error response with the given code and message.
func NewErrorResponse(code, message string) APIResponse {
	return APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
	}
}

// NewValidationErrorResponse creates a new validation error response with the given details.
func NewValidationErrorResponse(details []APIValidationError) APIResponse {
	return APIResponse{
		Error: &APIError{
			Code:    ErrorCodeValidation,
			Message: "Validation failed",
			Details: details,
		},
		Timestamp: time.Now(),
	}
}
