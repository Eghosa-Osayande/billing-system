-- name: GetInvoiceCounts :one
SELECT
    SUM(CASE WHEN payment_status = 'paid' THEN 1 ELSE 0 END)::text AS paid_count,
    SUM(CASE WHEN payment_status = 'unpaid' THEN 1 ELSE 0 END)::text AS unpaid_count,
    SUM(CASE WHEN payment_status = 'partial_paid' THEN 1 ELSE 0 END)::text AS partial_paid_count,
    SUM(CASE WHEN payment_status = 'overdue' THEN 1 ELSE 0 END)::text AS overdue_count
FROM
    invoice
WHERE
    business_id = $1;