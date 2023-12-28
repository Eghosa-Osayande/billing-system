package repository

import "errors"

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

func (repo *devRepo) GetUserVerificationDataWithEmail(email string) (*EmailVerificationModel, error) {

	if data, ok := repo.emailVerificationMap[email]; ok {
		return &data, nil
	}
	return nil, errors.New("not found")
}

func (repo *devRepo) UpdateUserEmailVerified(email string, emailVerified bool) (*UserModel, error) {

	if user, ok := repo.userMap[email]; ok {
		user.EmailVerified = emailVerified
		repo.userMap[email] = user
		return &user, nil
	}
	return nil, errors.New("user not found")
}

func (repo *devRepo) DeleteEmailVerificationDataByEmail(email string) error {

	delete(repo.userMap, email)

	return nil
}
