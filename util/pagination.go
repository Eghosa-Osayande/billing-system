package util

import "strconv"

type PagedResult[k any] struct {
	Items []k `json:"items"`
	Total int `json:"total"`
}

func NewPagedResult[k any](items []k, total int) *PagedResult[k] {

	return &PagedResult[k]{
		Items: items,
		Total: total,
	}
}

func GetPaginationFromQueries(queries map[string]string) (int, int) {
	limit := 10
	offset := 0
	if perpage, ok := queries["perpage"]; ok {
		perpagex, err := strconv.Atoi(perpage)
		if err == nil {

			if perpagex > 0 {
				limit = perpagex
			}
		}
	}
	
	if page, ok := queries["page"]; ok {
		pagex, err := strconv.Atoi(page)
		if err == nil {
			if offset = (pagex - 1)*limit; offset < 0 {
				offset = 0
			}
		}
	}
	return limit, offset
}
