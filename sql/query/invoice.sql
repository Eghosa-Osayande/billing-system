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

-- name: FindAllBusinessInvoices :many
SELECT
    *
FROM
    invoice
WHERE
    business_id = $1
ORDER BY
    created_at DESC;

-- name: FindInvoicesWhere :many
SELECT
    sqlc.embed(invoice),
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
    invoice.business_id = COALESCE(sqlc.narg('business_id'), invoice.business_id)
    AND invoice.client_id = COALESCE(sqlc.narg('client_id'), invoice.client_id)
    AND invoice.id = COALESCE(sqlc.narg('invoice_id'), invoice.id)
    AND (
        sqlc.narg('cursor_time')::timestamp IS NULL
        OR invoice.created_at <= sqlc.narg('cursor_time')
    )
    AND (
        sqlc.narg('cursor_id')::uuid IS NULL
        OR invoice.id < sqlc.narg('cursor_id')
    )
GROUP BY
    invoice.id
ORDER BY
    invoice.created_at DESC, invoice.id DESC
LIMIT
    COALESCE(sqlc.narg('limit'), 1);