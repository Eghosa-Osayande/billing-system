package repos

import (
	"blanq_invoice/database"
	"context"
)

type UserRepo struct {
	db *database.Queries
}

func NewUserRepo(db *database.Queries) *UserRepo {
	return &UserRepo{
		db: db,
	}

}

func (repo *UserRepo) GetUserProfileWhere(input database.GetUserProfileWhereParams) ([]database.GetUserProfileWhereRow, error) {
	db := repo.db
	ctx := context.Background()
	user, err := db.GetUserProfileWhere(ctx, input)
	if err != nil {
		return nil, err
	}

	return user, nil
}
