package repository

import (
	"blanq_invoice/sql_gen"
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql_gen.Queries
}

func NewPostgresRepo(connString string) (*PostgresRepo, error) {

	conn, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	db := sql_gen.New(conn)

	return &PostgresRepo{
		db: db,
	}, nil

}

func (repo *PostgresRepo) Close() {

}

func (repo *PostgresRepo) CreateOrUpdateUserEmailVerificationData(input *sql_gen.UserEmailVerification) error {
	ctx := context.Background()

	err:=repo.db.CreateOrUpdateUserEmailVerification(ctx, sql_gen.CreateOrUpdateUserEmailVerificationParams{
		Email:     input.Email,
		Code:      input.Code,
		ExpiresAt: input.ExpiresAt,
	})

	return err
}
func (repo *PostgresRepo) CreateUser(user *sql_gen.User) (*sql_gen.User, error) {
	ctx := context.Background()
	
	newuser,err:= repo.db.CreateUser(ctx, sql_gen.CreateUserParams{
		ID:            user.ID,
		Fullname:      user.Fullname,
		Email:         user.Email,
		Phone:         user.Phone,
		Password:      user.Password,
		EmailVerified: user.EmailVerified,
	})
	return &newuser,err;
	
}

func (repo *PostgresRepo) GetUserVerificationDataByEmail(email string) (*sql_gen.UserEmailVerification, error) {
	ctx := context.Background()
	
	user,err:= repo.db.FindUserEmailVerificationByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user,nil

}

func (repo *PostgresRepo) UpdateUserEmailVerified(email string, emailVerified bool) (*sql_gen.User, error) {
	ctx := context.Background()
	
	user,err:= repo.db.UpdateUserEmailVerifiedByEmail(ctx, sql_gen.UpdateUserEmailVerifiedByEmailParams{
		EmailVerified: emailVerified,
		Email:         email,
	})
	return &user,err
}

func (repo *PostgresRepo) DeleteEmailVerificationDataByEmail(email string) error {
	ctx := context.Background()
	
	err:= repo.db.DeleteUserEmailVerificationByEmail(ctx,
		email) 
	return err
}

func (repo *PostgresRepo) GetUserByEmail(email string) (*sql_gen.User, error) {
	ctx := context.Background()
	
	user,err:= repo.db.FindUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
