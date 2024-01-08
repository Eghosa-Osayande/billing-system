package database

import (
	// "database/sql/driver"
	// "encoding/json"
	// "database/sql"

	"github.com/jackc/pgx/v5"
)

func (q Queries) GetSqlDB() *pgx.Conn {
	return q.db.(*pgx.Conn)
}

type InvoiceWithItemsT[T any] struct  {
	Invoice 
	Items T `json:"items"`
	Clients *Client `json:"client"`
}

type InvoiceWithItems InvoiceWithItemsT[[]Invoiceitem]
type InvoiceWithItemsAny InvoiceWithItemsT[any]


// type InvoiceItem struct {
// 	Name     string   `json:"name" validate:"required" db:"name"`
// 	Price    float64  `json:"price" validate:"required" db:"price"`
// 	Quantity int      `json:"quantity" validate:"required" db:"quantity"`
// 	Discount *float64 `json:"discount" db:"discount"`
// }

// type InvoiceItemList []InvoiceItem

// func (h *InvoiceItemList) MarshalJSON() ([]byte, error) {
// 	type Alias InvoiceItemList
// 	return json.Marshal(&struct {
// 		*Alias
// 	}{
// 		Alias: (*Alias)(h),
// 	})
// }

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

// func (h *Invoice) MarshalJSON() ([]byte, error) {
// 	type Alias Invoice
// 	return json.Marshal(&struct {
// 		*Alias
// 		Items interface{} `json:"items"`
// 	}{
// 		Alias: (*Alias)(h),
// 		Items: json.RawMessage(h.Items),
// 	})
// }
