package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	StatusCode   int            `json:"statusCode"`
	Message      string         `json:"message"`
	ErrorMessage []ErrorMessage `json:"error"`
}

func (err *ValidationError) Error() string {
	return "validation failed"
}

func getErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fieldError.Param()
	case "gte":
		return "Should be greater than " + fieldError.Param()
	}
	return "Unknown error"
}

func ValidateRequest(err error) error {
	if err == nil {
		return nil
	}
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMessages := make([]ErrorMessage, len(validationErrors))
		for i, fieldError := range validationErrors {
			errorMessages[i] = ErrorMessage{fieldError.Field(), getErrorMessage(fieldError)}
		}
		return &ValidationError{
			StatusCode:   http.StatusNotAcceptable,
			Message:      "validation failed",
			ErrorMessage: errorMessages,
		}
	}
	return err
}
