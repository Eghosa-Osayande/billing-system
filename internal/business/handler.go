package business

import (
	"blanq_invoice/database"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BusinessHandler struct {
	repo *BusinessRepo
}

func NewBusinessHandler(repo *BusinessRepo) *BusinessHandler {
	return &BusinessHandler{repo: repo}
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
		business, err := handler.repo.FindBusinessByUserID(userUUID)
		if business != nil {
			return ctx.JSON(util.SuccessMessage[*database.Business]("Business found", business))
		} else {
			log.Println(err)
			return ctx.JSON(util.SuccessMessage[*database.Business]("No business found", nil))
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
		business, err := handler.repo.FindBusinessByUserID(userUUID)
		if business != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "User already has a business")
		}
		if err != nil {
			log.Println(err)
		}

		business, err = handler.repo.CreateBusiness(&database.CreateBusinessParams{
			ID:             uuid.New(),
			BusinessName:   input.BusinessName,
			BusinessAvatar: input.BusinessAvatar,
			OwnerID:        userUUID,
		})
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code, "Internal Server Error")
		}
		return ctx.JSON(util.SuccessMessage[*database.Business]("Business created successfully", business))
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
		business, err := handler.repo.FindBusinessByUserID(userUUID)
		if business == nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "User does not have a business yet")
		}
		if err != nil {
			log.Println(err)
		}
		
		
		business, err = handler.repo.UpdateBusiness(&database.UpdateBusinessParams{
			OwnerID:             userUUID,
			BusinessName:   input.BusinessName,
			BusinessAvatar: input.BusinessAvatar,
		})
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code, "Internal Server Error")
		}
		return ctx.JSON(util.SuccessMessage[*database.Business]("Business updated successfully", business))
	}
}
