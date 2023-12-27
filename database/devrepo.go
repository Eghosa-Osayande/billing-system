package database

type devRepo struct {
	userMap              map[string]UserModel
	emailVerificationMap map[string]EmailVerificationModel
}

func NewDevRepo() Repository {
	return &devRepo{
		userMap:              map[string]UserModel{},
		emailVerificationMap: map[string]EmailVerificationModel{},
	}
}

func (repo *devRepo) CheckExistingEmail(email string) (bool, error) {
	_, hasUser := repo.userMap[email]
	return hasUser, nil
}

func (repo *devRepo) PutEmailVerificationData(email string) error {
	repo.emailVerificationMap[email] = EmailVerificationModel{
		Email:     email,
		Code:      "1234",
		ExpiresAt: newEmailVerificationExpiration(),
	}

	return nil
}
func (repo *devRepo) CreateUser(input *CreateUserInput) (*UserModel, error) {

	user := &UserModel{
		Fullname:      input.Fullname,
		Email:         input.Email,
		Phone:         input.Phone,
		Password:      input.Password,
		EmailVerified: true,
	}
	repo.userMap[user.Email] = *user

	return user, nil
}
