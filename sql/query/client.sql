-- name: CreateClient :one
INSERT INTO
    client (
        id,
        created_at,
        business_id,
        fullname,
        email,
        phone
    )
VALUES
    (
        $1,
        timezone('utc', now()),
        $2,
        $3,
        $4,
        $5
    ) RETURNING *;

-- name: UpdateClient :one


UPDATE client
SET
    updated_at = timezone('utc', now()),
    fullname = COALESCE($2, fullname),
    email = COALESCE($3, email),
    phone = COALESCE($4, phone)
WHERE
    id = $1
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM
    client
WHERE
    id = $1;

-- name: GetClientsWhere :many

SELECT
    sqlc.embed(client),
    COUNT(client) OVER () AS total_count,
    COUNT(client) OVER (ORDER BY created_at ASC RANGE BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING) AS remaining_count
FROM
    client


WHERE
    (   
        business_id = $1
        or $1 is null
    )
    and (
        fullname ilike $2
        or $2 is null
    )
    and (
        email ilike $3
        or $3 is null
    )
    and (
        phone ilike $4
        or $4 is null
    )
ORDER BY
    created_at DESC
LIMIT
    $5 OFFSET $6;