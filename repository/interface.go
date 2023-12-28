package repository

import "time"

func newEmailVerificationExpiration() time.Time {
	return time.Time.Add(time.Now().UTC(), time.Duration(time.Duration.Seconds(20)))
}

func NewRepo() RepoInterface {
	return &devRepo{
		userMap:              map[string]UserModel{},
		emailVerificationMap: map[string]EmailVerificationModel{},
	}
}

type RepoInterface interface {
	CheckExistingEmail(email string) (bool, error)
	CreateUser(input *UserModel) (*UserModel, error)
	PutEmailVerificationData(email string) error
}


