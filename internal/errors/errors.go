package errors

import "fmt"

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s - %s", e.StatusCode, e.Code, e.Message)
}
func (e *APIError) ExitCode() int { return ExitAPI }

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string { return fmt.Sprintf("authentication error: %s", e.Message) }
func (e *AuthError) ExitCode() int { return ExitAuth }

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string { return fmt.Sprintf("%s not found: %s", e.Resource, e.ID) }
func (e *NotFoundError) ExitCode() int { return ExitNotFound }

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
}
func (e *ValidationError) ExitCode() int { return ExitValidation }

type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string { return fmt.Sprintf("network error: %s", e.Err) }
func (e *NetworkError) Unwrap() error { return e.Err }
func (e *NetworkError) ExitCode() int { return ExitNetwork }
