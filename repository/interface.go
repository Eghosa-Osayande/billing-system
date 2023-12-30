package repository

import "blanq_invoice/sql_gen"

type RepoInterface interface {
	CreateUser(input *sql_gen.User) (*sql_gen.User, error)

	CreateOrUpdateUserEmailVerificationData(input *sql_gen.UserEmailVerification) error

	DeleteEmailVerificationDataByEmail(email string) error

	GetUserVerificationDataByEmail(email string) (*sql_gen.UserEmailVerification, error)

	UpdateUserEmailVerified(email string, emailVerified bool) (*sql_gen.User, error)

	GetUserByEmail(email string) (*sql_gen.User, error)
}
