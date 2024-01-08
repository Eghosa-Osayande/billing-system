package repos

import (
	"blanq_invoice/database"
	"context"

	"github.com/google/uuid"
)

type ClientRepo struct {
	db *database.Queries
}

func NewClientRepo(db *database.Queries) *ClientRepo {

	return &ClientRepo{
		db: db,
	}

}

func (repo *ClientRepo) GetClients(businessId uuid.UUID) ([]database.Client, error) {
	ctx := context.Background()

	clientsList, err := repo.db.GetClientsByBusinessId(ctx, businessId)

	if isErrNoRows(err) {
		return clientsList, nil
	}

	if err != nil {
		return clientsList, err
	}

	return clientsList, nil

}

func (repo *ClientRepo) CreateClient(input *database.CreateClientParams) (*database.Client, error) {
	ctx := context.Background()

	client, err := repo.db.CreateClient(ctx, *input)

	if err != nil {
		return nil, err
	}

	return &client, nil

}

func (repo *ClientRepo) FindBusinessClientById(id, businessId uuid.UUID) (*database.Client, error) {
	ctx := context.Background()

	client, err := repo.db.FindBusinessClientByID(ctx, database.FindBusinessClientByIDParams{
		ID:         id,
		BusinessID: businessId,
	})

	if err != nil {
		return nil, err
	}

	return &client, nil
}
