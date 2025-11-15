package internal

import "fmt"

// APIError represents an error returned by the GitHub Models API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("GitHub Models API error: %d - %s", e.StatusCode, e.Message)
}

// NewAPIError creates a new APIError from status code and message
func NewAPIError(statusCode int, message string) error {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// WrapError adds context to an existing error
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
