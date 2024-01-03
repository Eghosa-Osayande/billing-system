-- name: CreateInvoiceItem :exec
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
    );

-- name: DeleteInvoiceItemByID :exec
DELETE FROM
    invoiceitem
WHERE
    id = $1;

-- name: FindInvoiceItemsByInvoiceID :many
SELECT
    sqlc.embed(invoice),sqlc.embed(invoiceitem)
FROM
    invoiceitem
JOIN invoiceitem ON invoiceitem.invoice_id = invoice.id
WHERE
    invoice_id = $1;