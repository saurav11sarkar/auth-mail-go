package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(v any) error {
	if err := validate.Struct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			messages := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				message := fmt.Sprintf("%s: failed %s validation", e.Field(), e.Tag())
				if e.Kind().String() == "string" && e.Tag() == "min" {
					message = fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
				}
				messages = append(messages, message)
			}
			return NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", strings.Join(messages, ", "))
		}
		return err
	}
	return nil
}
