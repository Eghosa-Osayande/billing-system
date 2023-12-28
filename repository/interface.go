package repository




func NewRepo() RepoInterface {
	return &devRepo{
		userMap:              map[string]UserModel{},
		emailVerificationMap: map[string]EmailVerificationModel{},
	}
}

type RepoInterface interface {
	Tx(action func() error) (error)

	CheckExistingEmail(email string) (bool, error)

	CreateUser(input *UserModel) (*UserModel, error)

	PutEmailVerificationData(input *EmailVerificationModel) error
	
	DeleteEmailVerificationDataByEmail(email string) error

	GetUserVerificationDataWithEmail(email string) (*EmailVerificationModel, error)

	UpdateUserEmailVerified (email string, emailVerified bool) (*UserModel, error)
}


