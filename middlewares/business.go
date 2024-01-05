package middlewares

import (
	"blanq_invoice/internal/repos"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserMustHaveBusinessMiddleware struct {
	*repos.ApiRepos
}

var instance *UserMustHaveBusinessMiddleware

func NewUserMustHaveBusinessMiddleware(repos *repos.ApiRepos) *UserMustHaveBusinessMiddleware {
	if instance == nil {
		instance = &UserMustHaveBusinessMiddleware{repos}
	}
	return instance
}

func UserMustHaveBusinessMiddlewareInstance() *UserMustHaveBusinessMiddleware {
	return instance
}

func (v *UserMustHaveBusinessMiddleware) Use(ctx *fiber.Ctx) error {

	userId, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	business, err := v.BusinessRepo.FindBusinessByUserID(*userId)
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code)
	}
	if business == nil {
		return ctx.JSON(util.NewSuccessResponseWithData[any]("Create a business first", nil))
	}
	ctx.Context().SetUserValue("business_id", business.ID)

	return ctx.Next()
}
