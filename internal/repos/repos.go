package repos

type ApiRepos struct {
	ClientRepo   *ClientRepo
	BusinessRepo *BusinessRepo
	AuthRepo     *AuthRepo
	InvoiceRepo *InvoiceRepo
}

func NewApiRepos(params ApiReposParams) *ApiRepos {
	return &ApiRepos{
		InvoiceRepo:  params.InvoiceRepo,
		ClientRepo:   params.ClientRepo,
		BusinessRepo: params.BusinessRepo,
		AuthRepo:     params.AuthRepo,
	}
}

type ApiReposParams struct {
	ClientRepo   *ClientRepo
	BusinessRepo *BusinessRepo
	AuthRepo     *AuthRepo
	InvoiceRepo *InvoiceRepo
}
