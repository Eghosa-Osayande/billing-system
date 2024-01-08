package middlewares

import (

	// "time"

	"blanq_invoice/util"

	"github.com/gofiber/fiber/v2"
)

func ErrorMessageMiddleware(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		errorList := []error{}

		if validationErr, ok := err.(util.ValidationError); ok {
			errorList = append(errorList, validationErr.ErrArr...)
		} else {
			if fiberErr, ok := err.(*fiber.Error); ok {
				c.Response().SetStatusCode(fiberErr.Code)
			} else {
				c.Response().SetStatusCode(500)
			}
			errorList = append(errorList, err)
		}

		return c.JSON(util.NewErrorMessage("error", errorList...))
	}
	return nil
}
