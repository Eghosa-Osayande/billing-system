package handlers

import (
	"blanq_invoice/repository"
	"blanq_invoice/util"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"time"
)

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

func (handler *AuthHandler) RegisterHandlers(router fiber.Router) {
	router.Post("/signup", handler.HandleSignup)
	router.Post("/login", handler.HandleLogin)
	router.Post("/verifyEmail", handler.HandleVerifyEmail)
	router.Post("/resendEmailOtp", handler.HandleResendEmailOtp)
}

type CreateUserInput struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

func newVerificationData(email string) *repository.EmailVerificationModel {
	return &repository.EmailVerificationModel{
		Email:     email,
		Code:      util.GenerateOTP(),
		ExpiresAt: time.Time.Add(time.Now().UTC(), time.Duration(time.Duration.Seconds(20))),
	}
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
		var verificationData *repository.EmailVerificationModel
		var createdUser *repository.UserModel
		handler.Repo.Tx(func() error {
			user, err := handler.Repo.CreateUser(newUser)
			if err != nil {
				return fiber.NewError(500, "User creation failed")
			}

			createdUser = user

			verificationData = newVerificationData(createdUser.Email)

			verificationErr := handler.Repo.PutEmailVerificationData(verificationData)

			if verificationErr != nil {
				return fiber.NewError(500, "User creation failed, OTP not created")
			}
			return nil
		})

		go util.SendEmailOTP(verificationData.Email, verificationData.Code)

		return ctx.JSON(util.SuccessMessage("User created successfully", createdUser))
	}

}

type SendEmailOtpInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (handler *AuthHandler) HandleResendEmailOtp(ctx *fiber.Ctx) error {
	body := ctx.Body()
	input := &SendEmailOtpInput{}

	if err := json.Unmarshal(body, input); err != nil {
		return util.ErrorInvalidJsonInput
	}
	if valErr := util.ValidateStruct(input); valErr != nil {
		return valErr
	}

	verificationData := newVerificationData(input.Email)

	handler.Repo.PutEmailVerificationData(verificationData)

	sendEmailErr := util.SendEmailOTP(verificationData.Email, verificationData.Code)

	if sendEmailErr != nil {
		return fiber.NewError(500, "OTP not sent")
	}

	return ctx.JSON(util.SuccessMessage("OTP sent successfully", nil))
}

type VerifyEmailOtpInput struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

func (handler *AuthHandler) HandleVerifyEmail(ctx *fiber.Ctx) error {
	body := ctx.Body()
	input := &VerifyEmailOtpInput{}

	if err := json.Unmarshal(body, input); err != nil {
		return util.ErrorInvalidJsonInput
	}
	if valErr := util.ValidateStruct(input); valErr != nil {
		return valErr
	}

	verificationData, err := handler.Repo.GetUserVerificationDataWithEmail(input.Email)

	if err != nil {
		return util.ApiError{Message: "Record not found"}
	}

	if verificationData.Code == input.Code {
		var updatedUser *repository.UserModel

		txErr := handler.Repo.Tx(func() error {
			user, err := handler.Repo.UpdateUserEmailVerified(input.Email, true)
			if err != nil {
				return util.ApiError{Message: err.Error()}
			}
			updatedUser = user
			deleteErr := handler.Repo.DeleteEmailVerificationDataByEmail(input.Email)
			if deleteErr != nil {
				return deleteErr
			}
			return nil
		})

		if txErr != nil {
			return txErr
		}

		return ctx.JSON(util.SuccessMessage("Email Verified", updatedUser))

	} else {
		return util.ApiError{Message: "Incorrect OTP"}
	}

}

func (handler *AuthHandler) HandleLogin(ctx *fiber.Ctx) error {
	return nil
}
