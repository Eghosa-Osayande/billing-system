package repos

import (
	"blanq_invoice/database"
	"blanq_invoice/util"
	"context"

	
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
	

	for index := range clients {
		clientList = append(clientList, clients[index].Client)
		total = int(clients[index].TotalCount)
		
	}

	return util.NewPagedResult[database.Client](clientList, total, ), nil

}

func (repo *ClientRepo) CreateClient(input *database.CreateClientParams) (*database.Client, error) {
	ctx := context.Background()

	client, err := repo.db.CreateClient(ctx, *input)

	if err != nil {
		return nil, err
	}

	return &client, nil

}
