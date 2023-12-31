package util

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)


var customValidator *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type ValidationError struct {	
	ErrArr []error 
}

func (v ValidationError) Error() string {
	return errors.Join(v.ErrArr...).Error()
}

func ValidateStruct(val interface{}) error {
	if err := customValidator.Struct(val); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errArr := make([]error, len(validationErrors))
		for index, validationErr := range validationErrors {
			field := validationErr.Field()
			msg := msgForTag(validationErr.Tag(), field)
			errArr[index] = errors.New(msg)

		}
		return ValidationError{ErrArr:errArr }
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
