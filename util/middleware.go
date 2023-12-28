package util

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorMessageMiddleware(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		errorList := []error{}
		if apiErrorList, ok := err.(ApiErrorList); ok {
			errorList = append(errorList, apiErrorList...)
			c.Response().SetStatusCode(400)
		} else {
			if fiberErr, ok := err.(*fiber.Error); ok {
				c.Response().SetStatusCode(fiberErr.Code)
			}else{
				c.Response().SetStatusCode(500)
			}
			errorList = append(errorList, err)
		}

		return c.JSON(ErrorMessage("error", errorList))
	}
	return nil
}
