package repository

import "time"

type UserModel struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password    string `json:"-"`
	EmailVerified bool `json:"emailVerified"`
}

type EmailVerificationModel struct {
	Email    string `json:"email"`
	Code    string `json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
}
