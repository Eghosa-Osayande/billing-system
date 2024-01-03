-- name: CreateOrUpdateUserEmailVerification :exec

INSERT INTO
	user_email_verifications (email, code, expires_at)
VALUES
	($1, $2, $3) ON CONFLICT (email) DO
UPDATE
SET
	code = $2,
	expires_at = $3,
	created_at = timezone('utc', now());

-- name: CreateUser :one
INSERT INTO
	users (
		fullname,
		email,
		password,
		email_verified
	)
VALUES
	($1, $2, $3, $4)
RETURNING *;

-- name: FindUserEmailVerificationByEmail :one
SELECT
	*
FROM
	user_email_verifications
WHERE
	email = $1;

-- name: UpdateUserEmailVerifiedByEmail :one
UPDATE
	users
SET
	email_verified = $1
WHERE
	email = $2
RETURNING *; 

-- name: DeleteUserEmailVerificationByEmail :exec
DELETE FROM
	user_email_verifications
WHERE
	email = $1;

-- name: FindUserByEmail :one
SELECT
	*
FROM
	users
WHERE
	email = $1
LIMIT 1;