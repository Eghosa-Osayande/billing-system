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
	router.Post("/update", handler.HandleAll)

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
			return ctx.JSON(util.NewSuccessResponseWithData[*util.PagedResult[database.Client]]("Create a business first", nil))
		}
		clients, err := handler.config.ClientRepo.GetClients(&database.GetClientsWhereParams{
			BusinessID: business.ID,
			Fullname:   nil,
			Email:      nil,
			Phone:      nil,
			Limit:      10,
			Offset:     0,
		})
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		if clients == nil {
			return ctx.JSON(util.NewSuccessResponseWithData[*util.PagedResult[database.Client]]("No clients found", nil))
		}

		return ctx.JSON(util.NewSuccessResponseWithData[*util.PagedResult[database.Client]]("Clients found", clients))

	}
}

type CreateClientInput struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    *string `json:"email" validate:"omitnil,email"`
	Phone    *string `json:"phone" validate:"omitnil"`
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
			ID:        uuid.New(),
			Fullname:  &input.Fullname,
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