package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/util"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	config *repos.ApiRepos
}

func NewAuthHandler(config *repos.ApiRepos) *AuthHandler {
	return &AuthHandler{
		config: config,
	}
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
}

// Signup godoc
// @Tags Authentication
// @Summary Signup
// @Description Create a new user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the Bearer prefix"
// @Param CreateUserInput body CreateUserInput true " "
// @Success 200 {object}  util.SuccessResponseWithData[database.User]
// @Failure 500 {object}  util.ErrorResponse
// @Router /auth/signup [post]
func (handler *AuthHandler) HandleSignup(ctx *fiber.Ctx) error {
	body := ctx.Body()
	createuserInput := &CreateUserInput{}

	if err := json.Unmarshal(body, createuserInput); err != nil {
		return util.ErrorInvalidJsonInput
	}
	if valErr := util.ValidateStruct(createuserInput); valErr != nil {
		return valErr
	}

	existingUser, err := handler.config.AuthRepo.GetUserByEmail(createuserInput.Email)

	if err != nil {
		log.Println(err)

	}

	if existingUser != nil {
		return fiber.NewError(400, "Email already exists")
	} else {
		hashedPass, hashErr := bcrypt.GenerateFromPassword([]byte(createuserInput.Password), bcrypt.DefaultCost)

		if hashErr != nil {
			return hashErr
		}
		newUser := database.CreateUserParams{
			ID:            uuid.New(),
			Fullname:      createuserInput.Fullname,
			Email:         createuserInput.Email,
			Phone:         createuserInput.Phone,
			Password:      string(hashedPass),
			EmailVerified: false,
		}
		
		createdUser, err := handler.config.AuthRepo.CreateUser(&newUser)
		if err != nil {
			return fiber.NewError(500, "User creation failed")
		}

		verificationData := &database.CreateOrUpdateUserEmailVerificationParams{
			Email:     createuserInput.Email,
			Code:      GenerateOTP(),
			ExpiresAt: time.Now().UTC().Add(time.Duration(5) * time.Minute)}

		verificationErr := handler.config.AuthRepo.CreateOrUpdateUserEmailVerificationData(verificationData)

		if verificationErr != nil {
			return fiber.NewError(500, "OTP sending failed")
		}

		go SendEmailOTP(verificationData.Email, verificationData.Code)

		return ctx.JSON(util.NewSuccessResponseWithData("User created successfully", createdUser))
	}

}

type ResendEmailOtpInput struct {
	Email string `json:"email" validate:"required,email"`
}

// Signup godoc
// @Tags Authentication
// @Summary Resend EmailOTP
// @Description Resend EmailOTP
// @Accept json
// @Produce json
// @Param ResendEmailOtpInput body ResendEmailOtpInput true " "
// @Success 200 {object}  util.SuccessResponse
// @Failure 500 {object}  util.ErrorResponse
// @Router /auth/resendEmailOtp [post]
func (handler *AuthHandler) HandleResendEmailOtp(ctx *fiber.Ctx) error {
	body := ctx.Body()
	input := &ResendEmailOtpInput{}

	if err := json.Unmarshal(body, input); err != nil {
		return util.ErrorInvalidJsonInput
	}
	if valErr := util.ValidateStruct(input); valErr != nil {
		return valErr
	}

	user, err := handler.config.AuthRepo.GetUserByEmail(input.Email)

	if err != nil || user == nil {
		return fiber.NewError(404, "Account with email does not exists")
	}

	verificationData := &database.CreateOrUpdateUserEmailVerificationParams{
		Email:     input.Email,
		Code:      GenerateOTP(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(5) * time.Minute)}

	err = handler.config.AuthRepo.CreateOrUpdateUserEmailVerificationData(verificationData)

	if err != nil {
		return fiber.NewError(500, "OTP sending failed")
	}

	sendEmailErr := SendEmailOTP(verificationData.Email, verificationData.Code)

	if sendEmailErr != nil {
		return fiber.NewError(500, "OTP not sent")
	}

	return ctx.JSON(util.NewSuccessResponse("OTP sent successfully"))
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

	verificationData, err := handler.config.AuthRepo.GetUserVerificationDataByEmail(input.Email)

	if err != nil {
		return fiber.NewError(400, "Record not found")
	}

	if verificationData.Code == input.Code {

		if verificationData.ExpiresAt.Before(time.Now().UTC()) {
			// TODO: optimise removal of expired otp
			// handler.config.AuthRepo.DeleteEmailVerificationDataByEmail(input.Email)
			return fiber.NewError(400, "OTP has expired")
		}

		_, err := handler.config.AuthRepo.UpdateUserEmailVerified(input.Email, true)

		if err != nil {
			return fiber.NewError(400, err.Error())
		}

		if deleteErr := handler.config.AuthRepo.DeleteEmailVerificationDataByEmail(input.Email); deleteErr != nil {
			log.Println(deleteErr)
		}

		updatedUser, err := handler.config.AuthRepo.GetUserByEmail(input.Email)

		if err != nil {
			return fiber.NewError(400, "Unknown Error")
		}

		return ctx.JSON(util.NewSuccessResponseWithData("Email Verified", updatedUser))

	} else {
		return fiber.NewError(400, "Incorrect OTP")
	}

}

type LoginUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	User database.User `json:"user"`
	Auth map[string]string
}

func (handler *AuthHandler) HandleLogin(ctx *fiber.Ctx) error {
	body := ctx.Body()
	input := &LoginUserInput{}

	if err := json.Unmarshal(body, input); err != nil {
		return util.ErrorInvalidJsonInput
	}
	if valErr := util.ValidateStruct(input); valErr != nil {
		return valErr
	}
	repo := handler.config.AuthRepo

	user, err := repo.GetUserByEmail(input.Email)
	if err != nil {
		return fiber.NewError(400, "Invalid login details")
	}

	incorrectPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if incorrectPassword == nil {
		accessDuration := time.Duration(24) * time.Hour * 30
		claims := jwt.MapClaims{}
		claims["email"] = user.Email
		claims["user_id"] = user.ID
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(accessDuration).UTC())

		key := os.Getenv("TOKENKEY")

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		accessToken, err := token.SignedString([]byte(key))
		if err != nil {
			log.Println("Error signing token", err)
			return fiber.NewError(500)
		}

		return ctx.JSON(util.NewSuccessResponseWithData("Logged In", LoginUserResponse{
			User: *user,
			Auth: map[string]string{
				"accessToken": accessToken,
				"expires_by":  fmt.Sprintf("%v", time.Now().Add(accessDuration).UTC()),
			},
		}))
	} else {
		return fiber.NewError(400, "Invalid login details")
	}
}

func GenerateOTP() string {
	arr := make([]int, 4)
	otpString := ""

	for range arr {
		otpString += fmt.Sprintf("%v", rand.Intn(10))
	}

	return otpString
}

func SendEmailOTP(email string, otp string) error {
	fmt.Printf("Sent otp (%v) to %v\n", otp, email)
	return nil
}
