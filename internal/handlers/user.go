package handlers

import (
	"blanq_invoice/internal/repos"
	"blanq_invoice/util"

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
	router.Get("/me", handler.HandleMe)
}


func (handler *UserHandler) HandleMe(ctx *fiber.Ctx)error{
	return ctx.JSON(util.NewSuccessResponseWithData[any]("yes",[]any{}))
}