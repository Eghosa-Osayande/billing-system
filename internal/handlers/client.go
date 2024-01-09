package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	router.Post("/update", handler.HandleUpdateClient)

}

type FetchClientFilter struct {
	ClientID   *uuid.UUID `json:"client_id" validate:"omitnil"`
	Limit      *int32     `json:"limit" validate:"omitnil"`
	Cursor     *string    `json:"cursor"`
	SearchTerm *string    `json:"search_term" validate:"omitnil"`
}

func (handler *ClientHandler) HandleAll(ctx *fiber.Ctx) error {
	input, err := ValidateRequestBody[*FetchClientFilter](ctx.Body(), &FetchClientFilter{})
	if err != nil {
		return err
	}

	businessId, err := util.GetUserBusinessIdFromContext(ctx)
	if err != nil {
		log.Println("Error getting business id", err)
		return fiber.NewError(fiber.ErrUnauthorized.Code, "unauthorized")
	}
	var ctime pgtype.Timestamptz
	var cid *int64
	if input.Cursor != nil {
		cursortime, cursorId, err := util.DecodeCursor(*input.Cursor)
		if err != nil {
			log.Println("Error decoding cursor", err)
			return fiber.NewError(fiber.ErrBadRequest.Code, "invalid-cursor")
		}

		ctime = pgtype.Timestamptz{Time: cursortime, Valid: true}
		cid = &cursorId

		if err != nil {
			log.Println("Error decoding cursor timestamp", err)
			return fiber.NewError(fiber.ErrBadRequest.Code, "invalid-cursor")
		}
		cid = &cursorId

	}

	clients, err := handler.config.ClientRepo.FindClientsWhere(&database.FindClientsWhereParams{
		ID:         input.ClientID,
		BusinessID: businessId,
		Fullname:   input.SearchTerm,
		Email:      input.SearchTerm,
		Phone:      input.SearchTerm,
		CursorTime: ctime,
		CursorID:   cid,
		Limit:      input.Limit,
	})
	if err != nil {
		log.Println("Error finding clients", err)
		if isErrNoRows(err) {
			return ctx.JSON(util.NewSuccessResponseWithData("No clients found", make([]any, 0)))
		}
		return fiber.NewError(fiber.ErrInternalServerError.Code)
	}

	return ctx.JSON(util.NewSuccessResponseWithData("Clients found", util.ListToPagedResult[database.Client](
		clients,
		func(item database.Client) (time.Time, int64) {
			return item.CreatedAt.Time, item.CountID
		},
	)))

}

type CreateClientInput struct {
	Fullname string  `json:"fullname" validate:"required"`
	Email    *string `json:"email" validate:"omitnil,email"`
	Phone    *string `json:"phone" validate:"omitnil,e164"`
}

func (handler *ClientHandler) HandleCreateClient(ctx *fiber.Ctx) error {
	input, valErr := ValidateRequestBody(ctx.Body(), &CreateClientInput{})
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

type UpdateClientInput struct {
	Fullname *string   `json:"fullname" validate:"omitnil,min=1"`
	Email    *string   `json:"email" validate:"omitnil,email"`
	Phone    *string   `json:"phone" validate:"omitnil,e164"`
	ClientID uuid.UUID `json:"client_id" validate:"required"`
}

func (h *ClientHandler) HandleUpdateClient(ctx *fiber.Ctx) error {
	input, valerr := ValidateRequestBody[*UpdateClientInput](ctx.Body(), &UpdateClientInput{})

	if valerr != nil {
		return valerr
	}

	r, err := h.config.ClientRepo.UpdateClient(&database.UpdateClientParams{
		ID:       input.ClientID,
		Fullname: input.Fullname,
		Email:    input.Email,
		Phone:    input.Phone,
	})

	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code)
	}

	return ctx.JSON(util.NewSuccessResponseWithData[*database.Client]("Client updated", r))

}
