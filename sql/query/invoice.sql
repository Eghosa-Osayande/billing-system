-- name: CreateInvoice :one
INSERT INTO
    invoice (
        business_id,
        currency,
        currency_symbol,
        payment_due_date,
        date_of_issue,
        notes,
        payment_method,
        client_id,
        shipping_fee_type,
        shipping_fee,
        total,
        tax,
        payment_status
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
        $10,
        $11,
        $12,
        COALESCE(sqlc.narg('payment_status'),'unpaid')::varchar
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
    client_id = COALESCE($7, client_id),
    shipping_fee_type = COALESCE($8, shipping_fee_type),
    shipping_fee = COALESCE($9, shipping_fee),
    total = COALESCE($10, total),
    payment_status = COALESCE(sqlc.narg('payment_status'), payment_status)
WHERE
    id = $1 RETURNING *;

-- name: DeleteInvoiceById :exec
Delete From
    invoice
WHERE
    id = $1;

-- name: FindInvoiceById :one
SELECT
    *
FROM
    invoice
WHERE
    (invoice.id = $1);

-- name: FindInvoicesWhere :many
SELECT
    sqlc.embed(invoice),
    JSON_AGG(client.*) as client,
    JSON_AGG(invoiceitem.*) as items
FROM
    invoice
    LEFT JOIN invoiceitem ON invoice.id = invoiceitem.invoice_id
    LEFT JOIN client ON invoice.client_id = client.id
WHERE
    (
        (
            (
                sqlc.narg('business_id') :: uuid is not null
                and invoice.business_id = sqlc.narg('business_id')
            )
            or sqlc.narg('business_id') is null
        )
        and (
            (
                sqlc.narg('client_id') :: uuid is not null
                and invoice.client_id = sqlc.narg('client_id')
            )
            or sqlc.narg('client_id') is null
        ) 
        and (
            (
                sqlc.narg('invoice_id') :: uuid is not null
                and invoice.id = sqlc.narg('invoice_id')
            )
            or sqlc.narg('invoice_id') is null
        )
        and (
            (
                invoice.created_at <= sqlc.narg('cursor_time') and
                sqlc.narg('cursor_time')  is not null
                
            )
            or sqlc.narg('cursor_time') is null
        )
        and (
            (
                invoice.count_id < sqlc.narg('cursor_id') and
                sqlc.narg('cursor_id')  is not null
            
            )
            or sqlc.narg('cursor_id')  is null
        )

    )
    
GROUP BY
    invoice.id,
    client.id
ORDER BY
    invoice.created_at DESC,
    invoice.count_id DESC
LIMIT
    sqlc.narg('limit')::INTEGER;