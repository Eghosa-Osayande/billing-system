package api

import (
	"blanq_invoice/database"
	"blanq_invoice/util"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var errorUserExists = fiber.NewError(400, "User already exists")
var errorUserCreationFailed = fiber.NewError(400, "Creation Failed")

func (server *ApiServer) HandleSignup(w http.ResponseWriter, r *http.Request) {
	body := r.Body

	input := new(database.SendVerifyEmailOtpInput)

	marshalError := json.NewDecoder(body).Decode(input)

	if marshalError != nil {
		util.RespondWithJson(w, 400, nil, []util.ApiError{
			{
				Field: "",
				Msg:   "Invalid Input",
			},
		})
		return
	}

	validationError := customValidator.Struct(*input)

	if validationError != nil {
		errors := []util.ApiError{}
		for _, err := range validationError.(validator.ValidationErrors) {
			errors = append(errors, util.ApiError{
				Field: err.Field(),
				Msg:   msgForTag(err.Tag()),
			})

		}
		util.RespondWithJson(w, 400, nil, errors)
		return
	}

	emailExists, _ := server.Repo.CheckExistingEmail(input.Email)

	if emailExists {
		util.RespondWithJson(w, 400, nil, []util.ApiError{
			{
				Field: "",
				Msg:   errorUserExists.Error(),
			},
		})
		return
	} else {
		err := server.Repo.PutEmailVerificationData(input.Email)
		if err != nil {
			util.RespondWithJson(w, 400, nil, []util.ApiError{
				{
					Field: "",
					Msg:   errorUserCreationFailed.Error(),
				},
			})
		}
		
		util.RespondWithJson(w, 200, struct {
			Msg string `json:"msg"`
		}{
			Msg: "OTP sent to email",
		}, nil)
		return
	}
}

func (server *ApiServer) HandleResendOtp(w http.ResponseWriter, r *http.Request) {

}

func (server *ApiServer) HandleVerifyEmailOtp(w http.ResponseWriter, r *http.Request) {

}

func (server *ApiServer) HandleLogin(w http.ResponseWriter, r *http.Request) {

}
