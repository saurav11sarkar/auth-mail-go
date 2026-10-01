package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

type AppError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap preserves errors.Is/errors.As through application errors.
func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(status int, code string, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func HandlerError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)
	json.NewEncoder(w).Encode(appErr)
}
