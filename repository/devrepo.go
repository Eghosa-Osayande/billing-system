package repository

type devRepo struct {
	userMap              map[string]UserModel
	emailVerificationMap map[string]EmailVerificationModel
}

func (repo *devRepo) CheckExistingEmail(email string) (bool, error) {
	_, hasUser := repo.userMap[email]
	return hasUser, nil
}

func (repo *devRepo) PutEmailVerificationData(input *EmailVerificationModel) error {
	repo.emailVerificationMap[input.Email] = EmailVerificationModel{
		Email:     input.Email,
		Code:      input.Code,
		ExpiresAt: input.ExpiresAt,
	}

	return nil
}
func (repo *devRepo) CreateUser(input *UserModel) (*UserModel, error) {

	user := &UserModel{
		Fullname:      input.Fullname,
		Email:         input.Email,
		Phone:         input.Phone,
		Password:      input.Password,
		EmailVerified: input.EmailVerified,
	}
	repo.userMap[user.Email] = *user

	return user, nil
}

func (repo *devRepo) Tx(action func() error) error {

	return action()
}
