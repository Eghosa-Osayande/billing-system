package repos

import "github.com/jackc/pgx/v5"

type ApiRepos struct {
	ClientRepo   *ClientRepo
	BusinessRepo *BusinessRepo
	AuthRepo     *AuthRepo
	InvoiceRepo  *InvoiceRepo
	*UserRepo
}

func NewApiRepos(params ApiReposParams) *ApiRepos {
	return &ApiRepos{
		InvoiceRepo:  params.InvoiceRepo,
		ClientRepo:   params.ClientRepo,
		BusinessRepo: params.BusinessRepo,
		AuthRepo:     params.AuthRepo,
		UserRepo:     params.UserRepo,
	}
}

type ApiReposParams struct {
	ClientRepo   *ClientRepo
	BusinessRepo *BusinessRepo
	AuthRepo     *AuthRepo
	InvoiceRepo  *InvoiceRepo
	*UserRepo
}

func isErrNoRows (err error) bool {
	return err==pgx.ErrNoRows;
}