-- name: CreateBusiness :one
INSERT INTO
	business (
		id,
		created_at,
		business_name,
		business_avatar,
		owner_id
	)
VALUES
	(
		$1,
		timezone('utc', now()),
		$2,
		$3,
		$4
	) RETURNING *;

-- name: UpdateBusiness :one
UPDATE
	business
SET
	updated_at = timezone('utc', now()),
	business_name = $2,
	business_avatar = $3
WHERE
	owner_id = $1 
RETURNING *;

-- name: FindBusinessByUserID :one
SELECT
	*
FROM
	business
WHERE
	owner_id = $1
LIMIT
	1;