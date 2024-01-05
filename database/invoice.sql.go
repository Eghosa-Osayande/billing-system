// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: invoice.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const createInvoice = `-- name: CreateInvoice :one
INSERT INTO
    invoice (
        business_id,
        currency,
        payment_due_date,
        date_of_issue,
        notes,
        payment_method,
        payment_status,
        client_id,
        shipping_fee_type,
        shipping_fee
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10
    ) RETURNING id, created_at, updated_at, deleted_at, business_id, currency, payment_due_date, date_of_issue, notes, payment_method, payment_status, client_id, shipping_fee_type, shipping_fee
`

type CreateInvoiceParams struct {
	BusinessID      uuid.UUID        `db:"business_id" json:"business_id"`
	Currency        *string          `db:"currency" json:"currency"`
	PaymentDueDate  *time.Time       `db:"payment_due_date" json:"payment_due_date"`
	DateOfIssue     *time.Time       `db:"date_of_issue" json:"date_of_issue"`
	Notes           *string          `db:"notes" json:"notes"`
	PaymentMethod   *string          `db:"payment_method" json:"payment_method"`
	PaymentStatus   *string          `db:"payment_status" json:"payment_status"`
	ClientID        *uuid.UUID       `db:"client_id" json:"client_id"`
	ShippingFeeType *string          `db:"shipping_fee_type" json:"shipping_fee_type"`
	ShippingFee     *decimal.Decimal `db:"shipping_fee" json:"shipping_fee"`
}

func (q *Queries) CreateInvoice(ctx context.Context, arg CreateInvoiceParams) (Invoice, error) {
	row := q.db.QueryRow(ctx, createInvoice,
		arg.BusinessID,
		arg.Currency,
		arg.PaymentDueDate,
		arg.DateOfIssue,
		arg.Notes,
		arg.PaymentMethod,
		arg.PaymentStatus,
		arg.ClientID,
		arg.ShippingFeeType,
		arg.ShippingFee,
	)
	var i Invoice
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.BusinessID,
		&i.Currency,
		&i.PaymentDueDate,
		&i.DateOfIssue,
		&i.Notes,
		&i.PaymentMethod,
		&i.PaymentStatus,
		&i.ClientID,
		&i.ShippingFeeType,
		&i.ShippingFee,
	)
	return i, err
}

const deleteInvoiceById = `-- name: DeleteInvoiceById :exec
Delete From
    invoice
WHERE
    id = $1
`

func (q *Queries) DeleteInvoiceById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteInvoiceById, id)
	return err
}

const findAllBusinessInvoices = `-- name: FindAllBusinessInvoices :many
SELECT
    id, created_at, updated_at, deleted_at, business_id, currency, payment_due_date, date_of_issue, notes, payment_method, payment_status, client_id, shipping_fee_type, shipping_fee
FROM
    invoice
WHERE
    business_id = $1
ORDER BY
    created_at DESC
`

func (q *Queries) FindAllBusinessInvoices(ctx context.Context, businessID uuid.UUID) ([]Invoice, error) {
	rows, err := q.db.Query(ctx, findAllBusinessInvoices, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Invoice
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.BusinessID,
			&i.Currency,
			&i.PaymentDueDate,
			&i.DateOfIssue,
			&i.Notes,
			&i.PaymentMethod,
			&i.PaymentStatus,
			&i.ClientID,
			&i.ShippingFeeType,
			&i.ShippingFee,
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

const findInvoiceById = `-- name: FindInvoiceById :one
SELECT
    invoice.id, invoice.created_at, invoice.updated_at, invoice.deleted_at, invoice.business_id, invoice.currency, invoice.payment_due_date, invoice.date_of_issue, invoice.notes, invoice.payment_method, invoice.payment_status, invoice.client_id, invoice.shipping_fee_type, invoice.shipping_fee,
    invoiceitem.id, invoiceitem.created_at, invoiceitem.invoice_id, invoiceitem.title, invoiceitem.price, invoiceitem.quantity, invoiceitem.discount, invoiceitem.discount_type
FROM
    invoice
    JOIN invoiceitem ON invoiceitem.invoice_id = invoice.id
WHERE
    (invoice.id = $1)
`

type FindInvoiceByIdRow struct {
	Invoice     Invoice     `db:"invoice" json:"invoice"`
	Invoiceitem Invoiceitem `db:"invoiceitem" json:"invoiceitem"`
}

func (q *Queries) FindInvoiceById(ctx context.Context, id uuid.UUID) (FindInvoiceByIdRow, error) {
	row := q.db.QueryRow(ctx, findInvoiceById, id)
	var i FindInvoiceByIdRow
	err := row.Scan(
		&i.Invoice.ID,
		&i.Invoice.CreatedAt,
		&i.Invoice.UpdatedAt,
		&i.Invoice.DeletedAt,
		&i.Invoice.BusinessID,
		&i.Invoice.Currency,
		&i.Invoice.PaymentDueDate,
		&i.Invoice.DateOfIssue,
		&i.Invoice.Notes,
		&i.Invoice.PaymentMethod,
		&i.Invoice.PaymentStatus,
		&i.Invoice.ClientID,
		&i.Invoice.ShippingFeeType,
		&i.Invoice.ShippingFee,
		&i.Invoiceitem.ID,
		&i.Invoiceitem.CreatedAt,
		&i.Invoiceitem.InvoiceID,
		&i.Invoiceitem.Title,
		&i.Invoiceitem.Price,
		&i.Invoiceitem.Quantity,
		&i.Invoiceitem.Discount,
		&i.Invoiceitem.DiscountType,
	)
	return i, err
}

const findInvoicesWhere = `-- name: FindInvoicesWhere :many
SELECT
    invoice.id, invoice.created_at, invoice.updated_at, invoice.deleted_at, invoice.business_id, invoice.currency, invoice.payment_due_date, invoice.date_of_issue, invoice.notes, invoice.payment_method, invoice.payment_status, invoice.client_id, invoice.shipping_fee_type, invoice.shipping_fee,
    JSON_AGG(
        jsonb_build_object(
            'id',
            invoiceitem.id,
            'created_at',
            invoiceitem.created_at,
            'invoice_id',
            invoiceitem.invoice_id,
            'title',
            invoiceitem.title,
            'price',
            invoiceitem.price,
            'quantity',
            invoiceitem.quantity,
            'discount',
            invoiceitem.discount,
            'discount_type',
            invoiceitem.discount_type
        )
    ) as items
FROM
    invoice
    JOIN invoiceitem ON invoiceitem.invoice_id = invoice.id
WHERE
    invoice.business_id = COALESCE($1, invoice.business_id)
    AND invoice.client_id = COALESCE($2, invoice.client_id)
    AND invoice.id = COALESCE($3, invoice.id)
    AND (
        $4::timestamp IS NULL
        OR invoice.created_at <= $4
    )
    AND (
        $5::uuid IS NULL
        OR invoice.id < $5
    )
GROUP BY
    invoice.id
ORDER BY
    invoice.created_at DESC, invoice.id DESC
LIMIT
    COALESCE($6, 1)
`

type FindInvoicesWhereParams struct {
	BusinessID *uuid.UUID  `db:"business_id" json:"business_id"`
	ClientID   *uuid.UUID  `db:"client_id" json:"client_id"`
	InvoiceID  *uuid.UUID  `db:"invoice_id" json:"invoice_id"`
	CursorTime *time.Time  `db:"cursor_time" json:"cursor_time"`
	CursorID   *uuid.UUID  `db:"cursor_id" json:"cursor_id"`
	Limit      interface{} `db:"limit" json:"limit"`
}

type FindInvoicesWhereRow struct {
	Invoice Invoice `db:"invoice" json:"invoice"`
	Items   []byte  `db:"items" json:"items"`
}

func (q *Queries) FindInvoicesWhere(ctx context.Context, arg FindInvoicesWhereParams) ([]FindInvoicesWhereRow, error) {
	rows, err := q.db.Query(ctx, findInvoicesWhere,
		arg.BusinessID,
		arg.ClientID,
		arg.InvoiceID,
		arg.CursorTime,
		arg.CursorID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindInvoicesWhereRow
	for rows.Next() {
		var i FindInvoicesWhereRow
		if err := rows.Scan(
			&i.Invoice.ID,
			&i.Invoice.CreatedAt,
			&i.Invoice.UpdatedAt,
			&i.Invoice.DeletedAt,
			&i.Invoice.BusinessID,
			&i.Invoice.Currency,
			&i.Invoice.PaymentDueDate,
			&i.Invoice.DateOfIssue,
			&i.Invoice.Notes,
			&i.Invoice.PaymentMethod,
			&i.Invoice.PaymentStatus,
			&i.Invoice.ClientID,
			&i.Invoice.ShippingFeeType,
			&i.Invoice.ShippingFee,
			&i.Items,
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

const updateInvoice = `-- name: UpdateInvoice :one
Update
    invoice
SET
    updated_at = timezone('utc', now()),
    currency = COALESCE($2, currency),
    payment_due_date = COALESCE($3, payment_due_date),
    date_of_issue = COALESCE($4, date_of_issue),
    notes = COALESCE($5, notes),
    payment_method = COALESCE($6, payment_method),
    payment_status = COALESCE($7, payment_status),
    client_id = COALESCE($8, client_id),
    shipping_fee_type = COALESCE($9, shipping_fee_type),
    shipping_fee = COALESCE($10, shipping_fee)
WHERE
    id = $1 RETURNING id, created_at, updated_at, deleted_at, business_id, currency, payment_due_date, date_of_issue, notes, payment_method, payment_status, client_id, shipping_fee_type, shipping_fee
`

type UpdateInvoiceParams struct {
	ID              uuid.UUID        `db:"id" json:"id"`
	Currency        *string          `db:"currency" json:"currency"`
	PaymentDueDate  *time.Time       `db:"payment_due_date" json:"payment_due_date"`
	DateOfIssue     *time.Time       `db:"date_of_issue" json:"date_of_issue"`
	Notes           *string          `db:"notes" json:"notes"`
	PaymentMethod   *string          `db:"payment_method" json:"payment_method"`
	PaymentStatus   *string          `db:"payment_status" json:"payment_status"`
	ClientID        *uuid.UUID       `db:"client_id" json:"client_id"`
	ShippingFeeType *string          `db:"shipping_fee_type" json:"shipping_fee_type"`
	ShippingFee     *decimal.Decimal `db:"shipping_fee" json:"shipping_fee"`
}

func (q *Queries) UpdateInvoice(ctx context.Context, arg UpdateInvoiceParams) (Invoice, error) {
	row := q.db.QueryRow(ctx, updateInvoice,
		arg.ID,
		arg.Currency,
		arg.PaymentDueDate,
		arg.DateOfIssue,
		arg.Notes,
		arg.PaymentMethod,
		arg.PaymentStatus,
		arg.ClientID,
		arg.ShippingFeeType,
		arg.ShippingFee,
	)
	var i Invoice
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.BusinessID,
		&i.Currency,
		&i.PaymentDueDate,
		&i.DateOfIssue,
		&i.Notes,
		&i.PaymentMethod,
		&i.PaymentStatus,
		&i.ClientID,
		&i.ShippingFeeType,
		&i.ShippingFee,
	)
	return i, err
}
