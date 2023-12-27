package api

import "github.com/go-playground/validator/v10"

var customValidator *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func msgForTag(tag string) string {
    switch tag {
    case "required":
        return "This field is required"
    case "email":
        return "Invalid email"
    }
    return ""
}