// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: auth.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrUpdateUserEmailVerification = `-- name: CreateOrUpdateUserEmailVerification :exec

INSERT INTO
	user_email_verifications (email, code, expires_at)
VALUES
	($1, $2, $3) ON CONFLICT (email) DO
UPDATE
SET
	code = $2,
	expires_at = $3,
	created_at = timezone('utc', now())
`

type CreateOrUpdateUserEmailVerificationParams struct {
	Email     string             `db:"email" json:"email"`
	Code      string             `db:"code" json:"code"`
	ExpiresAt pgtype.Timestamptz `db:"expires_at" json:"expires_at"`
}

func (q *Queries) CreateOrUpdateUserEmailVerification(ctx context.Context, arg CreateOrUpdateUserEmailVerificationParams) error {
	_, err := q.db.Exec(ctx, createOrUpdateUserEmailVerification, arg.Email, arg.Code, arg.ExpiresAt)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO
	users (
		fullname,
		email,
		password,
		email_verified
	)
VALUES
	($1, $2, $3, $4)
RETURNING id, created_at, updated_at, deleted_at, fullname, email, password, email_verified
`

type CreateUserParams struct {
	Fullname      string `db:"fullname" json:"fullname"`
	Email         string `db:"email" json:"email"`
	Password      string `db:"password" json:"password"`
	EmailVerified bool   `db:"email_verified" json:"email_verified"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Fullname,
		arg.Email,
		arg.Password,
		arg.EmailVerified,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Fullname,
		&i.Email,
		&i.Password,
		&i.EmailVerified,
	)
	return i, err
}

const deleteUserEmailVerificationByEmail = `-- name: DeleteUserEmailVerificationByEmail :exec
DELETE FROM
	user_email_verifications
WHERE
	email = $1
`

func (q *Queries) DeleteUserEmailVerificationByEmail(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, deleteUserEmailVerificationByEmail, email)
	return err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT
	id, created_at, updated_at, deleted_at, fullname, email, password, email_verified
FROM
	users
WHERE
	email = $1
LIMIT 1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Fullname,
		&i.Email,
		&i.Password,
		&i.EmailVerified,
	)
	return i, err
}

const findUserEmailVerificationByEmail = `-- name: FindUserEmailVerificationByEmail :one
SELECT
	email, created_at, code, expires_at
FROM
	user_email_verifications
WHERE
	email = $1
`

func (q *Queries) FindUserEmailVerificationByEmail(ctx context.Context, email string) (UserEmailVerification, error) {
	row := q.db.QueryRow(ctx, findUserEmailVerificationByEmail, email)
	var i UserEmailVerification
	err := row.Scan(
		&i.Email,
		&i.CreatedAt,
		&i.Code,
		&i.ExpiresAt,
	)
	return i, err
}

const updateUserEmailVerifiedByEmail = `-- name: UpdateUserEmailVerifiedByEmail :one
UPDATE
	users
SET
	email_verified = $1
WHERE
	email = $2
RETURNING id, created_at, updated_at, deleted_at, fullname, email, password, email_verified
`

type UpdateUserEmailVerifiedByEmailParams struct {
	EmailVerified bool   `db:"email_verified" json:"email_verified"`
	Email         string `db:"email" json:"email"`
}

func (q *Queries) UpdateUserEmailVerifiedByEmail(ctx context.Context, arg UpdateUserEmailVerifiedByEmailParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserEmailVerifiedByEmail, arg.EmailVerified, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Fullname,
		&i.Email,
		&i.Password,
		&i.EmailVerified,
	)
	return i, err
}
