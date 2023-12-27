package database

import "time"

func newEmailVerificationExpiration() time.Time {
	return time.Time.Add(time.Now().UTC(), time.Duration(time.Duration.Seconds(20)))
}

type Repository interface {
	CheckExistingEmail(email string) (bool, error)
	CreateUser(input *CreateUserInput) (*UserModel, error) 
	PutEmailVerificationData(email string) error
}

