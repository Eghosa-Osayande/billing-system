package repos

import (
	"blanq_invoice/database"
	"context"
	"log"
)



type AuthRepo struct {
	db *database.Queries
}

func NewAuthRepo( db *database.Queries) (*AuthRepo) {
	return &AuthRepo{
		db: db,
	}

}


func (repo *AuthRepo) CreateOrUpdateUserEmailVerificationData(input *database.CreateOrUpdateUserEmailVerificationParams) error {
	ctx := context.Background()

	err:=repo.db.CreateOrUpdateUserEmailVerification(ctx, database.CreateOrUpdateUserEmailVerificationParams{
		Email:     input.Email,
		Code:      input.Code,
		ExpiresAt: input.ExpiresAt,
	})

	return err
}
func (repo *AuthRepo) CreateUser(user *database.CreateUserParams) (*database.User, error) {
	ctx := context.Background()
	
	newuser,err:= repo.db.CreateUser(ctx, database.CreateUserParams{
		ID:            user.ID,
		Fullname:      user.Fullname,
		Email:         user.Email,
		Phone:         user.Phone,
		Password:      user.Password,
		EmailVerified: user.EmailVerified,
	})
	return &newuser,err;
	
}

func (repo *AuthRepo) GetUserVerificationDataByEmail(email string) (*database.UserEmailVerification, error) {
	ctx := context.Background()
	
	user,err:= repo.db.FindUserEmailVerificationByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user,nil

}

func (repo *AuthRepo) UpdateUserEmailVerified(email string, emailVerified bool) (*database.User, error) {
	ctx := context.Background()
	
	user,err:= repo.db.UpdateUserEmailVerifiedByEmail(ctx, database.UpdateUserEmailVerifiedByEmailParams{
		EmailVerified: emailVerified,
		Email:         email,
	})
	return &user,err
}

func (repo *AuthRepo) DeleteEmailVerificationDataByEmail(email string) error {
	ctx := context.Background()
	
	err:= repo.db.DeleteUserEmailVerificationByEmail(ctx,
		email) 
	return err
}

func (repo *AuthRepo) GetUserByEmail(email string) (*database.User, error) {
	ctx := context.Background()
	
	user,err:= repo.db.FindUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
