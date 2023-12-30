package business

import (
	"github.com/gofiber/fiber/v2"
)

type BusinessHandler struct {
	
}

func NewBusinessHandler() (*BusinessHandler) {
	return &BusinessHandler{
		
	}
}

func (handler *BusinessHandler) RegisterHandlers(router fiber.Router) {
	router.Get("/me", handler.HandleGetBusiness)

	router.Get("/update", handler.HandleGetBusiness)

	

}

func (handler *BusinessHandler) HandleGetBusiness(ctx *fiber.Ctx) error {
	// body := ctx.Body()
	// input := &LoginUserInput{}

	// if err := json.Unmarshal(body, input); err != nil {
	// 	return util.ErrorInvalidJsonInput
	// }
	// if valErr := util.ValidateStruct(input); valErr != nil {
	// 	return valErr
	// }
	// repo := handler.Repo

	return nil
}
