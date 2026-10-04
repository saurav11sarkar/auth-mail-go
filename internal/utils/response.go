package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func JSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
		Error:   map[string]string{"code": code, "message": message},
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	code := "HTTP-" + http.StatusText(status)
	switch status {
	case http.StatusBadRequest:
		code = "BAD_REQUEST"
	case http.StatusUnauthorized:
		code = "UNAUTHORIZED"
	case http.StatusForbidden:
		code = "FORBIDDEN"
	case http.StatusNotFound:
		code = "NOT_FOUND"
	case http.StatusInternalServerError:
		code = "INTERNAL_SERVER_ERROR"
	}
	JSONError(w, status, code, message)
}

func DecodeJSON(r *http.Request, v any) error {
	const maxBodyBytes = 1048576 // 1MB
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(v)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		message := "Invalid JSON"
		switch {
		case errors.Is(err, io.EOF):
			message = "Request body must not be empty"
		case errors.As(err, &syntaxErr):
			message = fmt.Sprintf("Invalid JSON at byte %d: check commas, quotes, and brackets", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			message = "Incomplete JSON: check closing quotes and brackets"
		case errors.As(err, &typeErr):
			message = fmt.Sprintf("Invalid JSON type for field %q: expected %s", typeErr.Field, typeErr.Type)
		}
		return NewAppError(http.StatusBadRequest, "BAD_REQUEST", message)
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Request body must contain a single JSON value")
	}
	return nil
}
