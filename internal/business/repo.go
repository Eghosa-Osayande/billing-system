package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)





type BusinessRepo struct {	
	db *Queries
}


func NewBusinessRepo(conn *pgx.Conn) (*BusinessRepo) {
	
	db := New(conn)

	return &BusinessRepo{
		db: db,
	}

}



func (repo *BusinessRepo) CreateBusiness(input *CreateBusinessParams) (*Business, error) {
	ctx := context.Background()
	
	newbusiness,err:= repo.db.CreateBusiness(ctx,*input)
	if err != nil {
		return nil,err
	}

	return &newbusiness,nil;
	
}

func (repo *BusinessRepo) UpdateBusiness(input *UpdateBusinessParams) (*Business, error) {
	ctx := context.Background()
	
	newbusiness,err:= repo.db.UpdateBusiness(ctx,*input)
	if err != nil {
		return nil,err
	}

	return &newbusiness,nil;
	
}

func (repo *BusinessRepo) FindBusinessByUserID(id uuid.UUID) (*Business, error) {
	ctx := context.Background()
	
	business,err:= repo.db.FindBusinessByUserID(ctx,id)
	if err != nil {
		return nil,err
	}

	return &business,nil;
	
}