package repos

import (
	"blanq_invoice/database"
	"blanq_invoice/util"
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

func (repo *ClientRepo) GetClients(input *database.GetClientsWhereParams) (*util.PagedResult[database.Client], error) {
	ctx := context.Background()

	clients, err := repo.db.GetClientsWhere(ctx, *input)
	
	if database.IsErrNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	clientList := []database.Client{}
	total := 0
	remaining := 0

	for index := range clients {
		clientList = append(clientList, clients[index].Client)
		total = int(clients[index].TotalCount)
		remaining = int(clients[index].RemainingCount)
	}

	return util.NewPagedResult[database.Client](clientList, total, remaining), nil

}

func (repo *ClientRepo) UpdateBusiness(input *database.UpdateBusinessParams) (*database.Business, error) {
	ctx := context.Background()

	newbusiness, err := repo.db.UpdateBusiness(ctx, *input)
	if err != nil {
		return nil, err
	}

	return &newbusiness, nil

}

func (repo *ClientRepo) FindBusinessByUserID(id uuid.UUID) (*database.Business, error) {
	ctx := context.Background()

	business, err := repo.db.FindBusinessByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &business, nil

}
