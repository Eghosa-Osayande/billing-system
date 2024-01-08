package repos

import (
	"blanq_invoice/database"
	"context"

	"github.com/google/uuid"
)

type UserRepo struct {
	db *database.Queries
}

func NewUserRepo(db *database.Queries) *UserRepo {
	return &UserRepo{
		db: db,
	}

}

func (repo *UserRepo) FindUserById(id uuid.UUID) (*database.User, error) {
	db := repo.db
	ctx := context.Background()
	user, err := db.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}
