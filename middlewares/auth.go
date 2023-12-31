package middlewares

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

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
				log.Println("token expires at", r)

				if mapClaims, ok := claims.(jwt.MapClaims); ok {
					
					c.Context().SetUserValue("user_id", mapClaims["user_id"].(string))
					return c.Next()
				}
				return errors.New("invalid Bearer Token Claims")
			}
		} else {
			return errors.New("invalid Bearer Token ")
		}
	}

	log.Println(authHeaders)
	return nil
}
