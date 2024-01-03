// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: invoice_item.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const createInvoiceItem = `-- name: CreateInvoiceItem :exec
INSERT INTO
    invoiceitem (
        invoice_id,
        title,
        price,
        quantity,
        discount,
        discount_type
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
`

type CreateInvoiceItemParams struct {
	InvoiceID    uuid.UUID        `db:"invoice_id" json:"invoice_id"`
	Title        string           `db:"title" json:"title"`
	Price        decimal.Decimal  `db:"price" json:"price"`
	Quantity     decimal.Decimal  `db:"quantity" json:"quantity"`
	Discount     *decimal.Decimal `db:"discount" json:"discount"`
	DiscountType *string          `db:"discount_type" json:"discount_type"`
}

func (q *Queries) CreateInvoiceItem(ctx context.Context, arg CreateInvoiceItemParams) error {
	_, err := q.db.Exec(ctx, createInvoiceItem,
		arg.InvoiceID,
		arg.Title,
		arg.Price,
		arg.Quantity,
		arg.Discount,
		arg.DiscountType,
	)
	return err
}

const deleteInvoiceItemByID = `-- name: DeleteInvoiceItemByID :exec
DELETE FROM
    invoiceitem
WHERE
    id = $1
`

func (q *Queries) DeleteInvoiceItemByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteInvoiceItemByID, id)
	return err
}

const findInvoiceItemsByInvoiceID = `-- name: FindInvoiceItemsByInvoiceID :many
SELECT
    id, created_at, invoice_id, title, price, quantity, discount, discount_type
FROM
    invoiceitem
WHERE
    invoice_id = $1
`

func (q *Queries) FindInvoiceItemsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]Invoiceitem, error) {
	rows, err := q.db.Query(ctx, findInvoiceItemsByInvoiceID, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Invoiceitem
	for rows.Next() {
		var i Invoiceitem
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.InvoiceID,
			&i.Title,
			&i.Price,
			&i.Quantity,
			&i.Discount,
			&i.DiscountType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
