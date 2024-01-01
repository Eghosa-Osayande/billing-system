package repos

type ApiRepos struct {
	ClientRepo   *ClientRepo
	BusinessRepo *BusinessRepo
	AuthRepo     *AuthRepo
}

func NewApiRepos(params ApiReposParams) *ApiRepos {
	return &ApiRepos{

		ClientRepo:   params.ClientRepo,
		BusinessRepo: params.BusinessRepo,
		AuthRepo:     params.AuthRepo,
	}
}

type ApiReposParams struct {
	ClientRepo   *ClientRepo
	BusinessRepo *BusinessRepo
	AuthRepo     *AuthRepo
}
