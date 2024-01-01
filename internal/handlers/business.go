package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BusinessHandler struct {
	config *repos.ApiRepos
}

func NewBusinessHandler(config *repos.ApiRepos) *BusinessHandler {
	return &BusinessHandler{config: config}
}

func (handler *BusinessHandler) RegisterHandlers(router fiber.Router) {
	router.Get("/", handler.HandleGetBusiness)
	router.Post("/new", handler.HandleCreateBusiness)
	router.Post("/update", handler.HandleUpdateBusiness)

}

// write a handler for each route
func (handler *BusinessHandler) HandleGetBusiness(ctx *fiber.Ctx) error {
	if userId, ok := ctx.Context().UserValue("user_id").(string); !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
		}
		business, err := handler.config.BusinessRepo.FindBusinessByUserID(userUUID)
		if business != nil {
			return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("Business found", business))
		} else {
			log.Println(err)
			return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("No business found", nil))
		}

	}
}

type CreateBusinessInput struct {
	BusinessName   string  `db:"business_name" json:"business_name" validate:"required"`
	BusinessAvatar *string `db:"business_avatar" json:"business_avatar"`
}

func (handler *BusinessHandler) HandleCreateBusiness(ctx *fiber.Ctx) error {

	input, valErr := util.ValidateRequestBody(ctx.Body(), &CreateBusinessInput{})

	if valErr != nil {
		return valErr
	}

	if userId, ok := ctx.Context().UserValue("user_id").(string); !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
		}
		business, err := handler.config.BusinessRepo.FindBusinessByUserID(userUUID)
		if business != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "User already has a business")
		}
		if err != nil {
			log.Println(err)
		}

		business, err = handler.config.BusinessRepo.CreateBusiness(&database.CreateBusinessParams{
			ID:             uuid.New(),
			BusinessName:   input.BusinessName,
			BusinessAvatar: input.BusinessAvatar,
			OwnerID:        userUUID,
		})
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code, "Internal Server Error")
		}
		return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("Business created successfully", business))
	}

}

type UpdateBusinessInput struct {
	BusinessName   string  `db:"business_name" json:"business_name" validate:"required"`
	BusinessAvatar *string `db:"business_avatar" json:"business_avatar"`
}

func (handler *BusinessHandler) HandleUpdateBusiness(ctx *fiber.Ctx) error {
	input, valErr := util.ValidateRequestBody(ctx.Body(), &UpdateBusinessInput{})

	if valErr != nil {
		return valErr
	}

	if userId, ok := ctx.Context().UserValue("user_id").(string); !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
		}
		business, err := handler.config.BusinessRepo.FindBusinessByUserID(userUUID)
		if business == nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "User does not have a business yet")
		}
		if err != nil {
			log.Println(err)
		}

		business, err = handler.config.BusinessRepo.UpdateBusiness(&database.UpdateBusinessParams{
			OwnerID:        userUUID,
			BusinessName:   input.BusinessName,
			BusinessAvatar: input.BusinessAvatar,
		})
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code, "Internal Server Error")
		}
		return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("Business updated successfully", business))
	}
}
