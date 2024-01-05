package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ClientHandler struct {
	config *repos.ApiRepos
}

func NewClientHandler(config *repos.ApiRepos) *ClientHandler {
	return &ClientHandler{config: config}
}

func (handler *ClientHandler) RegisterHandlers(router fiber.Router) {
	router = router.Group("/clients").Use(middlewares.AuthenticatedUserMiddleware).Use(middlewares.UserMustHaveBusinessMiddlewareInstance().Use)
	
	router.Get("/all", handler.HandleAll)
	router.Post("/new", handler.HandleCreateClient)

}

func (handler *ClientHandler) HandleAll(ctx *fiber.Ctx) error {
	log.Println(util.GetUserIdFromContext(ctx))
	businessId, err := util.GetUserBusinessIdFromContext(ctx)
	if err != nil {
		return err
	}

	clients, err := handler.config.ClientRepo.GetClients(*businessId)
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code)
	}

	if clients == nil {
		return ctx.JSON(util.NewSuccessResponseWithData[any]("No clients found", nil))
	}

	return ctx.JSON(util.NewSuccessResponseWithData[[]database.Client]("Clients found", clients))

}

type CreateClientInput struct {
	Fullname string  `json:"fullname" validate:"required"`
	Email    *string `json:"email" validate:"omitnil,email"`
	Phone    *string `json:"phone" validate:"omitnil,e164"`
}

func (handler *ClientHandler) HandleCreateClient(ctx *fiber.Ctx) error {
	input, valErr := util.ValidateRequestBody(ctx.Body(), &CreateClientInput{})
	if valErr != nil {
		return valErr
	}

	businessId, err := util.GetUserBusinessIdFromContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	createClientParams := &database.CreateClientParams{
		Fullname:   input.Fullname,
		Email:      input.Email,
		Phone:      input.Phone,
		BusinessID: *businessId,
	}

	newclient, err := handler.config.ClientRepo.CreateClient(createClientParams)
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code)
	}
	return ctx.JSON(util.NewSuccessResponseWithData[*database.Client]("Client created", newclient))

}
