package handlers

import (
	"blanq_invoice/repository"
	"blanq_invoice/util"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// type SendVerifyEmailOtpInput struct {
// 	Email string `json:"email" validate:"required,email"`
// }

type CreateUserInput struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

// type CreateUserResponse struct {
// 	User *User `json:"user"`
// }

// type User struct {
// 	Fullname      string `json:"fullname"`
// 	Email         string `json:"email"`
// 	Phone         string `json:"phone"`
// 	Password      string `json:"password"`
// 	EmailVerified bool   `json:"email_verified"`
// }

// type LoginUserInput struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

type AuthHandler struct {
	Repo repository.RepoInterface
}

func (handler *AuthHandler) HandleSignup(ctx *fiber.Ctx) error {
	body := ctx.Body()
	createuserInput := &CreateUserInput{}

	if err := json.Unmarshal(body, createuserInput); err != nil {
		return util.ErrorInvalidJsonInput
	}
	if valErr := util.ValidateStruct(createuserInput); valErr != nil {
		return valErr
	}

	isExisting, err := handler.Repo.CheckExistingEmail(createuserInput.Email)

	if err != nil {
		return fiber.NewError(500, "Email lookup failed")
	}

	if isExisting {
		return fiber.NewError(400, "Email already exists")
	} else {
		newUser := &repository.UserModel{
			Fullname:      createuserInput.Fullname,
			Email:         createuserInput.Email,
			Phone:         createuserInput.Phone,
			Password:      createuserInput.Password,
			EmailVerified: false,
		}
		createdUser, err := handler.Repo.CreateUser(newUser)
		if err != nil {
			return fiber.NewError(500, "User creation failed")
		}
		return ctx.JSON(util.SuccessMessage("User created successfully", createdUser))
	}

}

func (handler *AuthHandler) HandleResendOtp(ctx *fiber.Ctx) error {
	return nil
}

func (handler *AuthHandler) HandleVerifyEmailOtp(ctx *fiber.Ctx) error {
	return nil
}

func (handler *AuthHandler) HandleLogin(ctx *fiber.Ctx) error {
	return nil
}
