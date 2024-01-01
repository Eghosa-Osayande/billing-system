-- name: CreateInvoice :one
INSERT INTO
    invoice (
        id,
        created_at,
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
        timezone('utc', now()),
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11
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

-- name: GetInvoiceWhere :many
SELECT
    sqlc.embed(invoice),
    COUNT(*) OVER () AS total_count,
    COUNT(*) OVER (
        ORDER BY
            created_at ASC RANGE BETWEEN CURRENT ROW
            AND UNBOUNDED FOLLOWING
    ) AS remaining_count
FROM
    invoice
WHERE
    (business_id = COALESCE($1, business_id))
    AND (id = COALESCE($2, id))
    AND (client_id = COALESCE($3, client_id))
    AND (payment_status = COALESCE($4, payment_status))
    AND (payment_method = COALESCE($5, payment_method))
    AND (
        shipping_fee_type = COALESCE($6, shipping_fee_type)
    )
    AND (currency = COALESCE($7, currency))
    AND (payment_status = COALESCE($8, payment_status))
    AND (payment_method = COALESCE($9, payment_method))
ORDER BY
    created_at ASC
LIMIT
    $10 OFFSET $11;