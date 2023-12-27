package database



type SendVerifyEmailOtpInput struct {
	Email string `json:"email" validate:"required,email"`
}

type CreateUserInput struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type CreateUserResponse struct {
	User *User `json:"user"`
}

type User struct {
	Fullname      string `json:"fullname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	EmailVerified bool   `json:"email_verified"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
