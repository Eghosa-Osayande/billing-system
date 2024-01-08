package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"io/fs"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

type BusinessHandler struct {
	config *repos.ApiRepos
}

func NewBusinessHandler(config *repos.ApiRepos) *BusinessHandler {
	return &BusinessHandler{config: config}
}

func (handler *BusinessHandler) RegisterHandlers(router fiber.Router) {
	router = router.Group("/business").Use(middlewares.AuthenticatedUserMiddleware)

	router.Get("/all", handler.HandleGetBusiness)
	router.Post("/new", handler.HandleCreateBusiness)

	router.Use(
		middlewares.UserMustHaveBusinessMiddlewareInstance().Use,
	).Post("/update", handler.HandleUpdateBusiness)

	router.Use(
		middlewares.UserMustHaveBusinessMiddlewareInstance().Use,
	).Post("/avatar", handler.HandleUploadBusinessAvatar)

}

// write a handler for each route
func (handler *BusinessHandler) HandleGetBusiness(ctx *fiber.Ctx) error {

	userUUID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	business, err := handler.config.BusinessRepo.FindBusinessByUserID(*userUUID)
	if business != nil {
		return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("Business found", business))
	} else {
		log.Println(err)
		return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("No business found", nil))
	}

}

type CreateBusinessInput struct {
	BusinessName   string  `db:"business_name" json:"business_name" validate:"required"`
	BusinessAvatar *string `db:"business_avatar" json:"business_avatar"`
}

func (handler *BusinessHandler) HandleUploadBusinessAvatar(ctx *fiber.Ctx) (error) {
	fileHeader, err := ctx.FormFile("avatar")
	
	if err != nil {
		log.Println("Error getting file header")
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println("Error opening file")
		return err
	}

	contents := make([]byte, fileHeader.Size)
	_, err = file.Read(contents)
	if err != nil {
		log.Println("Error reading file")
		return err
	}
	os.Mkdir("uploads",fs.ModePerm)
	savedFile, err := os.Create("uploads/" + fileHeader.Filename)
	if err != nil {
		log.Println("Error creating file")
		return err
	}

	_, err = savedFile.Write(contents)
	if err != nil {
		log.Println("Error writing file")
		return err	
	}
	err=savedFile.Close()
	if err != nil {
		log.Println("Error closing file")
		return err
	}

	return ctx.JSON(util.NewSuccessResponseWithData("File uploaded successfully", "avatar.png"))
}

func (handler *BusinessHandler) HandleCreateBusiness(ctx *fiber.Ctx) error {

	input, valErr := util.ValidateRequestBody(ctx.Body(), &CreateBusinessInput{})

	if valErr != nil {
		return valErr
	}

	userUUID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	business, err := handler.config.BusinessRepo.FindBusinessByUserID(*userUUID)
	if business != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "User already has a business")
	}
	if err != nil {
		log.Println(err)
	}

	business, err = handler.config.BusinessRepo.CreateBusiness(&database.CreateBusinessParams{
		BusinessName:   input.BusinessName,
		BusinessAvatar: input.BusinessAvatar,
		OwnerID:        *userUUID,
	})
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code, "Internal Server Error")
	}
	return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("Business created successfully", business))

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

	userUUID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	business, err := handler.config.BusinessRepo.UpdateBusiness(&database.UpdateBusinessParams{
		OwnerID:        *userUUID,
		BusinessName:   input.BusinessName,
		BusinessAvatar: input.BusinessAvatar,
	})
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code, "Internal Server Error")
	}

	return ctx.JSON(util.NewSuccessResponseWithData[*database.Business]("Business updated successfully", business))

}
