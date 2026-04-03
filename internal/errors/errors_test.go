package errors

import (
	"fmt"
	"testing"
)

func TestAPIError(t *testing.T) {
	err := &APIError{StatusCode: 400, Code: "BAD_REQUEST", Message: "invalid date"}
	if err.Error() != "API error (400): BAD_REQUEST - invalid date" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if err.ExitCode() != ExitAPI {
		t.Errorf("expected exit code %d, got %d", ExitAPI, err.ExitCode())
	}
}

func TestAuthError(t *testing.T) {
	err := &AuthError{Message: "invalid API key"}
	if err.ExitCode() != ExitAuth {
		t.Errorf("expected exit code %d, got %d", ExitAuth, err.ExitCode())
	}
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Resource: "document", ID: "S1234567"}
	if err.Error() != "document not found: S1234567" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if err.ExitCode() != ExitNotFound {
		t.Errorf("expected exit code %d, got %d", ExitNotFound, err.ExitCode())
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "date", Message: "invalid format"}
	if err.ExitCode() != ExitValidation {
		t.Errorf("expected exit code %d, got %d", ExitValidation, err.ExitCode())
	}
}

func TestNetworkError(t *testing.T) {
	err := &NetworkError{Err: fmt.Errorf("connection refused")}
	if err.ExitCode() != ExitNetwork {
		t.Errorf("expected exit code %d, got %d", ExitNetwork, err.ExitCode())
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		err      error
		expected int
	}{
		{&AuthError{Message: "no key"}, ExitAuth},
		{&APIError{StatusCode: 500}, ExitAPI},
		{fmt.Errorf("unknown"), ExitGeneral},
	}
	for _, tt := range tests {
		if got := GetExitCode(tt.err); got != tt.expected {
			t.Errorf("GetExitCode(%v) = %d, want %d", tt.err, got, tt.expected)
		}
	}
}
