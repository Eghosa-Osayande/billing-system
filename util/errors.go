package util

import (
	"github.com/gofiber/fiber/v2"
)



var (
	ErrorInvalidJsonInput= fiber.NewError(400, "Invalid Json Input")
)