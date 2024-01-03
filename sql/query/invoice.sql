-- name: CreateInvoice :one
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
    ) RETURNING *;

-- name: UpdateInvoice :one
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
    id = $1 RETURNING *;

-- name: DeleteInvoiceById :exec
Delete From
    invoice
WHERE
    id = $1;

-- name: FindInvoiceById :one
SELECT
    sqlc.embed(invoice),
    sqlc.embed(invoiceitem)
FROM
    invoice
JOIN invoiceitem ON invoiceitem.invoice_id = invoice.id
WHERE
    (invoice.id = $1);