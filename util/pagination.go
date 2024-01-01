package util


type PagedResult[k any] struct {
	Items []k `json:"items"`
	Total int `json:"total"`
	Remaining int `json:"remaining"`
}

func NewPagedResult[k any](items []k, total int, remaining int) *PagedResult[k] {
	
	return &PagedResult[k]{
		Items: items,
		Total: total,
		Remaining: remaining,
	}
}