package middlewares

import (

	// "time"

	"blanq_invoice/util"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorMessageMiddleware(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {

		if fiberErr, ok := err.(*fiber.Error); ok {
			c.Response().SetStatusCode(fiberErr.Code)
		} else {
			c.Response().SetStatusCode(500)
		}
		err = errors.Join(err)

		return c.JSON(util.NewErrorMessage("error", err))
	}
	return nil
}
