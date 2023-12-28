package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var customValidator *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(val interface{}) error {

	if err := customValidator.Struct(val); err != nil {

		validationErrors := err.(validator.ValidationErrors)

		errArr := make(ApiErrorList, len(validationErrors))
		for index, validationErr := range validationErrors {
			field := validationErr.Field()
			msg := msgForTag(validationErr.Tag(), field)
			errArr[index] = ApiError{Field: &field, Message: msg}

		}
		return errArr
	}
	return nil
}

func msgForTag(tag string, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%v is required", field)
	case "email":
		return "Email is invalid"
	default:
		return fmt.Sprintf("Error for %v", field)
	}

}
