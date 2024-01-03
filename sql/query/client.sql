-- name: CreateClient :one
INSERT INTO
    client (
       
        business_id,
        fullname,
        email,
        phone
    )
VALUES
    (
        sqlc.arg('business_id'),
        sqlc.arg('fullname'),
        sqlc.narg('email'),
        sqlc.narg('phone')
    ) RETURNING *;

-- name: UpdateClient :one
UPDATE
    client
SET
    updated_at = timezone('utc', now()),
    fullname = COALESCE($2, fullname),
    email = COALESCE($3, email),
    phone = COALESCE($4, phone)
WHERE
    id = $1 RETURNING *;

-- name: DeleteClient :exec
DELETE FROM
    client
WHERE
    id = $1;

-- name: GetClientsByBusinessId :many
SELECT
    *
FROM
    client
WHERE
    business_id = sqlc.arg('business_id')
ORDER BY
    created_at DESC;