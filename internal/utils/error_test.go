package utils

import (
	"errors"
	"fmt"
	"testing"
)

func TestAppErrorChain(t *testing.T) {
	cause := errors.New("database failure")
	appErr := &AppError{Status: 500, Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error", Err: cause}
	wrapped := fmt.Errorf("operation: %w", appErr)
	if !errors.Is(wrapped, cause) {
		t.Fatal("underlying cause lost")
	}
	var found *AppError
	if !errors.As(wrapped, &found) || found != appErr {
		t.Fatal("application error lost")
	}
}
