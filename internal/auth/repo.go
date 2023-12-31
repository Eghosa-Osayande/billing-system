package auth

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)



type AuthRepo struct {
	db *Queries
}

func NewAuthRepo(conn *pgx.Conn) (*AuthRepo) {

	db := New(conn)

	return &AuthRepo{
		db: db,
	}

}


func (repo *AuthRepo) CreateOrUpdateUserEmailVerificationData(input *CreateOrUpdateUserEmailVerificationParams) error {
	ctx := context.Background()

	err:=repo.db.CreateOrUpdateUserEmailVerification(ctx, CreateOrUpdateUserEmailVerificationParams{
		Email:     input.Email,
		Code:      input.Code,
		ExpiresAt: input.ExpiresAt,
	})

	return err
}
func (repo *AuthRepo) CreateUser(user *CreateUserParams) (*User, error) {
	ctx := context.Background()
	
	newuser,err:= repo.db.CreateUser(ctx, CreateUserParams{
		ID:            user.ID,
		Fullname:      user.Fullname,
		Email:         user.Email,
		Phone:         user.Phone,
		Password:      user.Password,
		EmailVerified: user.EmailVerified,
	})
	return &newuser,err;
	
}

func (repo *AuthRepo) GetUserVerificationDataByEmail(email string) (*UserEmailVerification, error) {
	ctx := context.Background()
	
	user,err:= repo.db.FindUserEmailVerificationByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user,nil

}

func (repo *AuthRepo) UpdateUserEmailVerified(email string, emailVerified bool) (*User, error) {
	ctx := context.Background()
	
	user,err:= repo.db.UpdateUserEmailVerifiedByEmail(ctx, UpdateUserEmailVerifiedByEmailParams{
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

func (repo *AuthRepo) GetUserByEmail(email string) (*User, error) {
	ctx := context.Background()
	
	user,err:= repo.db.FindUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
