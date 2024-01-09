package database

import (
	// "database/sql/driver"
	// "encoding/json"
	// "database/sql"

	"encoding/json"

	"github.com/jackc/pgx/v5"
)

func (q Queries) GetSqlDB() *pgx.Conn {
	return q.db.(*pgx.Conn)
}

type FullInvoice struct {
	Invoice
	Items   any `json:"items"`
	Clients any `json:"client"`
}

func removeNilValues(input []interface{}) []interface{} {
	var result = make([]interface{}, 0)

	for _, value := range input {
		if value != nil {
			result = append(result, value)
		}
	}

	return result
}

func (i *FindInvoicesWhereRow) ToFullInvoice() (*FullInvoice, error) {
	var client []any
	err := json.Unmarshal(i.Client, &client)
	if err != nil {
		return nil, err
	}

	var jsonItems []any

	err = json.Unmarshal(i.Items, &jsonItems)

	if err != nil {
		return nil, err
	}
	return &FullInvoice{
		Invoice: i.Invoice,
		Items:   removeNilValues(jsonItems),
		Clients: client[0],
	}, nil
}

type FullUserProfile struct {
	 GetUserProfileWhereRow
	Business any `json:"business"`
}

func (i *GetUserProfileWhereRow) ToFullUser() (*FullUserProfile, error) {
	var business []any
	err := json.Unmarshal(i.Business, &business)
	if err != nil {
		return nil, err
	}

	return &FullUserProfile{
		GetUserProfileWhereRow: *i,
		Business: business[0],
	}, nil
}

// func (h *InvoiceItemList) Value() (driver.Value, error) {
// 	println("value")

// 	if h != nil {
// 		return json.Marshal(h)
// 	}
// 	return nil, nil
// }

// func (h *InvoiceItemList) Scan(value interface{}) error {
// 	return json.Unmarshal(value.([]byte), h)
// }

// func (h *GetUserProfileWhereRow) MarshalJSON() ([]byte, error) {
// 	type Alias GetUserProfileWhereRow
// 	return json.Marshal(&struct {
// 		*Alias
// 		Items interface{} `json:"items"`
// 	}{
// 		Alias: (*Alias)(h),
		
// 	})
// }

// func (h *GetUserProfileWhereRow) UnmarshalJSON(data []byte) error {
// 	if err := json.Unmarshal(data, h); err != nil {
// 		return err
// 	}
// 	return nil
// }
