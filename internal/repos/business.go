package repos

import (
	"blanq_invoice/database"
	"context"

	"github.com/google/uuid"
)

type BusinessRepo struct {
	db *database.Queries
}

func NewBusinessRepo(db *database.Queries) *BusinessRepo {

	return &BusinessRepo{
		db: db,
	}

}

func (repo *BusinessRepo) CreateBusiness(input *database.CreateBusinessParams) (*database.Business, error) {
	ctx := context.Background()

	newbusiness, err := repo.db.CreateBusiness(ctx, *input)
	if err != nil {
		return nil, err
	}

	return &newbusiness, nil

}

func (repo *BusinessRepo) UpdateBusiness(input *database.UpdateBusinessParams) (*database.Business, error) {
	ctx := context.Background()

	newbusiness, err := repo.db.UpdateBusiness(ctx, *input)
	if err != nil {
		return nil, err
	}

	return &newbusiness, nil

}

func (repo *BusinessRepo) FindBusinessByUserID(id uuid.UUID) (*database.Business, error) {
	ctx := context.Background()

	business, err := repo.db.FindBusinessByUserID(ctx, id)
	if isErrNoRows(err) {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}

	return &business, nil

}
