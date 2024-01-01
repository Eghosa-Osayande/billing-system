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
    currency = case
        when $2 is null then currency
        else $2
    end,
    payment_due_date = case
        when $3 is null then payment_due_date
        else $3
    end,
    date_of_issue = case
        when $4 is null then date_of_issue
        else $4
    end,
    notes = case
        when $5 is null then notes
        else $5
    end,
    payment_method = case
        when $6 is null then payment_method
        else $6
    end,
    payment_status = case
        when $7 is null then payment_status
        else $7
    end,
    client_id = case
        when $8 is null then client_id
        else $8
    end,
    shipping_fee_type = case
        when $9 is null then shipping_fee_type
        else $9
    end,
    shipping_fee = case
        when $10 is null then shipping_fee
        else $10
    end
WHERE
    id = $1 RETURNING *;

-- name: DeleteInvoiceById :exec
Delete From
    invoice
WHERE
    id = $1;

-- name: GetInvoiceWhere :many

SELECT
    COUNT(*) OVER () AS total_count,
    COUNT(*) OVER (ORDER BY created_at ASC RANGE BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING) AS remaining_count,
    *
FROM
    invoice
WHERE
    (
        $1 is null
        or business_id = $1
    )
    and (
        $2 is null
        or id = $2
    )
    and (
        $3 is null
        or client_id = $3
    )
    and (
        $4 is null
        or payment_status = $4
    )
    and (
        $5 is null
        or payment_method = $5
    )
    and (
        $6 is null
        or shipping_fee_type = $6
    )
    and (
        $7 is null
        or currency = $7
    )
    and (
        $8 is null
        or payment_status = $8
    )
    and (
        $9 is null
        or payment_method = $9
    )

ORDER BY
    created_at ASC
LIMIT 
    $10
OFFSET 
    $11;