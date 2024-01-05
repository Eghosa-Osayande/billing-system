package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClientHandler struct {
	config *repos.ApiRepos
}

func NewClientHandler(config *repos.ApiRepos) *ClientHandler {
	return &ClientHandler{config: config}
}

func (handler *ClientHandler) RegisterHandlers(router fiber.Router) {
	router.Get("/all", handler.HandleAll)
	router.Post("/new", handler.HandleCreateClient)
	

}



func (handler *ClientHandler) HandleAll(ctx *fiber.Ctx) error {
	if userId, ok := ctx.Context().UserValue("user_id").(string); !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {
		
		userId, err := uuid.Parse(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		business, err := handler.config.BusinessRepo.FindBusinessByUserID(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		if business == nil {
			return ctx.JSON(util.NewSuccessResponseWithData[any]("Create a business first", nil))
		}
		
		clients, err := handler.config.ClientRepo.GetClients(business.ID,)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		if clients == nil {
			return ctx.JSON(util.NewSuccessResponseWithData[any]("No clients found", nil))
		}

		return ctx.JSON(util.NewSuccessResponseWithData[[]database.Client]("Clients found", clients))

	}
}

type CreateClientInput struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    *string `json:"email" validate:"omitnil,email"`
	Phone    *string `json:"phone" validate:"omitnil,e164"`
}

func (handler *ClientHandler) HandleCreateClient(ctx *fiber.Ctx) error {
	input, valErr := util.ValidateRequestBody(ctx.Body(), &CreateClientInput{})

	if valErr != nil {
		return valErr
	}

	if userId, ok := ctx.Context().UserValue("user_id").(string); !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {
		userId, err := uuid.Parse(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		business, err := handler.config.BusinessRepo.FindBusinessByUserID(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		if business == nil {
			return ctx.JSON(util.NewSuccessResponseWithData[*database.Client]("Create a business first", nil))
		}
		createClientParams:=&database.CreateClientParams{
			Fullname:  input.Fullname,
			Email:     input.Email,
			Phone:     input.Phone,
			BusinessID: business.ID,
		}

		newclient, err := handler.config.ClientRepo.CreateClient(createClientParams)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		return ctx.JSON(util.NewSuccessResponseWithData[*database.Client]("Client created", newclient))

	}
}