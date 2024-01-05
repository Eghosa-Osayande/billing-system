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

	clients, err := repo.db.GetClientsByBusinessId(ctx, businessId)

	if database.IsErrNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	clientList := []database.Client{}

	clientList = append(clientList, clients...)

	return clientList, nil

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
