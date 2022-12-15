package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email"
	case "boolean":
		return "this field must be of type boolean"
	case "oneof":
		return fmt.Sprintf("this field must be one of the following values: %v", fe.Param())
	}

	return "unknown error"
}

func FillErrors(ve validator.ValidationErrors) []ErrorMsg {
	out := make([]ErrorMsg, len(ve))

	for i, fe := range ve {
		out[i] = ErrorMsg{
			Field:   fe.Field(),
			Message: GetErrorMsg(fe),
		}
	}

	return out
}
