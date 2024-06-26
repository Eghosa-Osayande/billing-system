package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/util"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	router = router.Group("/auth")
	router.Post("/signup", handler.HandleSignup)
	router.Post("/login", handler.HandleLogin)
	router.Post("/verifyEmail", handler.HandleVerifyEmail)
	router.Post("/resendEmailOtp", handler.HandleResendEmailOtp)
}

type CreateUserInput struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,len=7"`
}



func (handler *AuthHandler) HandleSignup(ctx *fiber.Ctx) error {

	input, valErr := ValidateRequestBody(ctx.Body(), &CreateUserInput{})

	if valErr != nil {
		return valErr
	}

	existingUser, err := handler.config.AuthRepo.GetUserByEmail(input.Email)

	if err != nil {
		log.Println(err)

	}

	if existingUser != nil {
		return fiber.NewError(400, "Email already exists")
	} else {
		hashedPass, hashErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		if hashErr != nil {
			return hashErr
		}
		newUser := database.CreateUserParams{
			Fullname:      input.Fullname,
			Email:         input.Email,
			Password:      string(hashedPass),
			EmailVerified: false,
		}

		createdUser, err := handler.config.AuthRepo.CreateUser(&newUser)
		if err != nil {
			return fiber.NewError(500, "User creation failed")
		}

		verificationData := &database.CreateOrUpdateUserEmailVerificationParams{
			Email:     input.Email,
			Code:      generateOTP(),
			ExpiresAt: pgtype.Timestamptz{Time: generateOtpExpiration(), Valid: true},
		}

		verificationErr := handler.config.AuthRepo.CreateOrUpdateUserEmailVerificationData(verificationData)

		if verificationErr != nil {
			return fiber.NewError(500, "OTP sending failed")
		}

		go sendEmailOTP(verificationData.Email, verificationData.Code)

		updatedUser, err := handler.config.UserRepo.GetUserProfileWhere(database.GetUserProfileWhereParams{
			Email: &createdUser.Email,
		})

		if err != nil {
			return fiber.NewError(400, "Unknown Error")
		}
		if len(updatedUser) > 1 {
			log.Println("Multiple Accounts found")
			return fiber.NewError(400, "Unknown Error")
		}

		userprof, err := updatedUser[0].ToFullUser()
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
		}

		return ctx.JSON(util.NewSuccessResponseWithData("User created successfully", userprof))
	}

}

type ResendEmailOtpInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (handler *AuthHandler) HandleResendEmailOtp(ctx *fiber.Ctx) error {

	input, valErr := ValidateRequestBody(ctx.Body(), &ResendEmailOtpInput{})

	if valErr != nil {
		return valErr
	}

	user, err := handler.config.AuthRepo.GetUserByEmail(input.Email)

	if err != nil || user == nil {
		return fiber.NewError(404, "Account with email does not exists")
	}

	verificationData := &database.CreateOrUpdateUserEmailVerificationParams{
		Email:     input.Email,
		Code:      generateOTP(),
		ExpiresAt: pgtype.Timestamptz{Time: generateOtpExpiration(), Valid: true}}

	err = handler.config.AuthRepo.CreateOrUpdateUserEmailVerificationData(verificationData)

	if err != nil {
		return fiber.NewError(500, "OTP sending failed")
	}

	sendEmailErr := sendEmailOTP(verificationData.Email, verificationData.Code)

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

	input, valErr := ValidateRequestBody(ctx.Body(), &VerifyEmailOtpInput{})

	if valErr != nil {
		return valErr
	}

	verificationData, err := handler.config.AuthRepo.GetUserVerificationDataByEmail(input.Email)

	if err != nil {
		return fiber.NewError(400, "Record not found")
	}

	if verificationData.Code == input.Code {

		if verificationData.ExpiresAt.Time.Before(time.Now().UTC()) {
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

		updatedUser, err := handler.config.UserRepo.GetUserProfileWhere(database.GetUserProfileWhereParams{
			Email: &input.Email,
		})

		if err != nil {
			return fiber.NewError(400, "Unknown Error")
		}
		if len(updatedUser) > 1 {
			log.Println("Multiple Accounts found")
			return fiber.NewError(400, "Unknown Error")
		}

		userprof, err := updatedUser[0].ToFullUser()
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
		}

		return ctx.JSON(util.NewSuccessResponseWithData("Email Verified", userprof))

	} else {
		return fiber.NewError(400, "Incorrect OTP")
	}

}

type LoginUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (handler *AuthHandler) HandleLogin(ctx *fiber.Ctx) error {

	input, valErr := ValidateRequestBody(ctx.Body(), &LoginUserInput{})

	if valErr != nil {
		return valErr
	}

	userProfile, err := handler.config.UserRepo.GetUserProfileWhere(database.GetUserProfileWhereParams{
		Email: &input.Email,
	})

	if err != nil {
		log.Println(err)
		return fiber.NewError(400, "Invalid login details")
	}

	if noOfUsers := len(userProfile); noOfUsers != 1 {
		log.Println("Multiple or Zero Accounts found: ", noOfUsers)
		return fiber.NewError(400, "Invalid login details")
	}

	user := userProfile[0]

	incorrectPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if incorrectPassword == nil {
		accessDuration := time.Duration(24) * time.Hour * 30
		claims := jwt.MapClaims{}
		claims["user_id"] = user.ID
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(accessDuration).UTC())

		key := os.Getenv("TOKENKEY")

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		accessToken, err := token.SignedString([]byte(key))
		if err != nil {
			log.Println("Error signing token", err)
			return fiber.NewError(500)
		}

		userprof, err := user.ToFullUser()
		if err != nil {
			log.Println("Error converting user to full user", err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		return ctx.JSON(util.NewSuccessResponseWithData(
			"Logged In",
			map[string]any{
				"user": userprof,
				"auth": map[string]string{
					"accessToken": accessToken,
					"expires_by":  fmt.Sprintf("%v", time.Now().Add(accessDuration).UTC()),
				},
			},
		))
	} else {
		return fiber.NewError(400, "Invalid login details")
	}
}

func generateOtpExpiration() time.Time {
	return time.Now().UTC().Add(time.Duration(5) * time.Minute)
}

func generateOTP() string {
	arr := make([]int, 4)
	otpString := ""

	for range arr {
		otpString += fmt.Sprintf("%v", rand.Intn(10))
	}

	return "1234"
}

func sendEmailOTP(email string, otp string) error {
	fmt.Printf("Sent otp (%v) to %v\n", otp, email)
	return nil
}
