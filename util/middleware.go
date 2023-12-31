package util

import (
	"errors"
	"log"
	"os"
	"strings"
	// "time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ErrorMessageMiddleware(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		errorList := []error{}

		if validationErr, ok := err.(ValidationError); ok {
			errorList = append(errorList, validationErr.ErrArr...)
		} else {
			if fiberErr, ok := err.(*fiber.Error); ok {
				c.Response().SetStatusCode(fiberErr.Code)
			} else {
				c.Response().SetStatusCode(500)
			}
			errorList = append(errorList, err)
		}

		

		return c.JSON(errorMessage("error", errorList))
	}
	return nil
}

func AuthenticatedUserMiddleware(c *fiber.Ctx) error {
	authHeaders := c.GetReqHeaders()["Authorization"]

	if len(authHeaders) > 0 {
		authHeader := authHeaders[0]
		splitString := strings.Split(authHeader, "Bearer ")
		unauthenticatedError := fiber.ErrUnauthorized
		if len(splitString) == 2 {
			accessToken := splitString[1]

			token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Println("Unexpected signing method in auth token")
					return nil, errors.New("unexpected signing method in auth token")
				}
				return []byte(os.Getenv("TOKENKEY")), nil
			})
			if err != nil {
				log.Println("unable to parse claims", "error", err)
				return unauthenticatedError
			}

			claims := token.Claims
			if !token.Valid {
				return unauthenticatedError
			} else {
				r, _ := claims.GetExpirationTime()
				log.Println(r.UTC())
				c.Next()
			}
		} else {
			return errors.New("invalid Bearer Token ")
		}
	}

	log.Println(authHeaders)
	return nil
}
