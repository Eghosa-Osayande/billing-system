package repos

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

