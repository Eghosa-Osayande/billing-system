package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var customValidator *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type ValidationError struct {
	ErrArr []error
}

func (v ValidationError) Error() string {
	return errors.Join(v.ErrArr...).Error()
}

func validateStruct(val interface{}) error {

	if err := customValidator.Struct(val); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errArr := make([]error, len(validationErrors))
		for index, validationErr := range validationErrors {
			field := validationErr.Field()
			msg := msgForTag(validationErr.Tag(), field, validationErr.Error())
			errArr[index] = errors.New(msg)

		}
		return ValidationError{ErrArr: errArr}
	}
	return nil
}

func msgForTag(tag string, field string, errorMessage string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%v is required", field)
	case "email":
		return "Email is invalid"
	default:
		return fmt.Sprintf("Error for %v, %v", field, errorMessage)
	}

}

func ValidateRequestBody[T any](body []byte, output T) (T, error) {

	if err := json.Unmarshal(body, output); err != nil {
		return output, fiber.NewError(fiber.ErrBadRequest.Code,err.Error())
	}
	if valErr := validateStruct(output); valErr != nil {
		return output, fiber.NewError(fiber.ErrBadRequest.Code,valErr.Error())
	}

	return output, nil

}
