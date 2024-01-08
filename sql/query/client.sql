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
    fullname = COALESCE(sqlc.narg(fullname), fullname),
    email = COALESCE(sqlc.narg(email), email),
    phone = COALESCE(sqlc.narg(phone), phone)
WHERE
    id = $1 RETURNING *;

-- name: FindBusinessClientByID :one
SELECT
    *
FROM
    client
WHERE
    id = $1
    AND business_id = $2;

-- name: FindClientsWhere :many
SELECT
    *
FROM
    client
WHERE
    (
        client.id = sqlc.narg('id') or sqlc.narg('id') IS NULL
    )
    AND (
        client.business_id = sqlc.narg('business_id') or sqlc.narg('business_id') IS NULL
    )
    AND (
        client.fullname ilike sqlc.narg('fullname') or sqlc.narg('fullname') IS NULL
    )
    AND (
        client.email ilike sqlc.narg('email')
        or sqlc.narg('email') IS NULL
    )
    AND (
        client.phone ilike sqlc.narg('phone')
        or sqlc.narg('phone') IS NULL
    )
    AND (
        client.created_at <= sqlc.narg('cursor_time')
        or sqlc.narg('cursor_time') IS NULL
    )
    AND (
        client.count_id < sqlc.narg('cursor_id')
        or sqlc.narg('cursor_id') IS NULL
    )
ORDER BY
    client.created_at DESC,
    client.count_id DESC
LIMIT
    sqlc.narg('limit');