-- name: CreateInvoiceItem :one
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
RETURNING *;

-- name: DeleteInvoiceItemByID :exec
DELETE FROM
    invoiceitem
WHERE
    id = $1;

-- name: FindInvoiceItemsByInvoiceId :many
SELECT
    *
FROM
    invoiceitem
WHERE
    (invoiceitem.invoice_id = $1);