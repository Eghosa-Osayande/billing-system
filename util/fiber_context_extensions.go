package util

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserIdFromContext(ctx *fiber.Ctx) (*uuid.UUID, error) {
	if userId, ok := ctx.Context().UserValue("user_id").(uuid.UUID); !ok {
		return nil, fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {
		return &userId, nil
	}
}

func GetUserBusinessIdFromContext(ctx *fiber.Ctx) (*uuid.UUID, error) {

	if businessId, ok := ctx.Context().UserValue("business_id").(uuid.UUID); !ok {
		log.Println("Business ID not provided")
		return nil, fiber.NewError(fiber.ErrInternalServerError.Code)
	} else {
		return &businessId, nil
	}
}
