package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	config *repos.ApiRepos
}

func NewUserHandler(config *repos.ApiRepos) *UserHandler {
	return &UserHandler{
		config: config,
	}
}

func (handler *UserHandler) RegisterHandlers(router fiber.Router) {
	router = router.Use(middlewares.AuthenticatedUserMiddleware)
	router.Get("/me", handler.HandleMe)
}

func (handler *UserHandler) HandleMe(ctx *fiber.Ctx) error {
	id, err := util.GetUserIdFromContext(ctx)

	if err != nil {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "unauthorized")
	}

	user, err := handler.config.UserRepo.GetUserProfileWhere(database.GetUserProfileWhereParams{
		ID: id,
	})

	if len(user) > 1 {
		log.Println("Multiple Accounts found")
		return fiber.NewError(fiber.ErrInternalServerError.Code, "User account not found")
	}

	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	userprof, err := user[0].ToFullUser()
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	return ctx.JSON(util.NewSuccessResponseWithData[any]("User profile", userprof))
}
