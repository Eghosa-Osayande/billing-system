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

type userBusinessHolder struct {
	User    *database.User    `json:"user"`
	Business *database.Business `json:"business"`
}

// User Account Details godoc
// @Tags Account Details
// @Summary Get User Account Details
// @Description 
// @Accept json
// @Produce json
// @Success 200 {object}  util.SuccessResponseWithData[userBusinessHolder]
// @Failure 500 {object}  util.ErrorResponse
// @Router /me [get]
func (handler *UserHandler) HandleMe(ctx *fiber.Ctx) error {
	id, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "unauthorized")
	}

	user,err:=handler.config.UserRepo.FindUserById(*id)
	if err!=nil{
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code,err.Error())
	}
	business,err:=handler.config.BusinessRepo.FindBusinessByUserID(*id)

	if err!=nil{
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code,err.Error())
	}
	return ctx.JSON(util.NewSuccessResponseWithData[userBusinessHolder]("User account", userBusinessHolder{
		User:    user,
		Business: business,
	}))
}
