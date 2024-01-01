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
UPDATE
    client
SET
    updated_at = timezone('utc', now()),
    fullname = case
        when $2 is null then fullname
        else $2
    end,
    email = case
        when $3 is null then email
        else $3
    end,
    phone = case
        when $4 is null then phone
        else $4
    end
WHERE
    id = $1;

-- name: DeleteClient :exec
DELETE FROM
    client
WHERE
    id = $1;

-- name: FindClientsWhere
SELECT
    COUNT(*) OVER () AS total_count,
    COUNT(*) OVER (
        ORDER BY
            created_at ASC RANGE BETWEEN CURRENT ROW
            AND UNBOUNDED FOLLOWING
    ) AS remaining_count,
    *
FROM
    client
WHERE
    (
        $1 is null
        or business_id = $1
    )
    and (
        $2 is null
        or fullname ilike $2
    )
    and (
        $3 is null
        or email ilike $3
    )
    and (
        $4 is null
        or phone ilike $4
    )
ORDER BY
    created_at ASC
LIMIT
    $5 OFFSET $6;