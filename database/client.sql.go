// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: client.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createClient = `-- name: CreateClient :one
INSERT INTO
    client (
        business_id,
        fullname,
        email,
        phone
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4
    ) RETURNING count_id, id, created_at, updated_at, deleted_at, business_id, fullname, email, phone
`

type CreateClientParams struct {
	BusinessID uuid.UUID `db:"business_id" json:"business_id"`
	Fullname   string    `db:"fullname" json:"fullname"`
	Email      *string   `db:"email" json:"email"`
	Phone      *string   `db:"phone" json:"phone"`
}

func (q *Queries) CreateClient(ctx context.Context, arg CreateClientParams) (Client, error) {
	row := q.db.QueryRow(ctx, createClient,
		arg.BusinessID,
		arg.Fullname,
		arg.Email,
		arg.Phone,
	)
	var i Client
	err := row.Scan(
		&i.CountID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.BusinessID,
		&i.Fullname,
		&i.Email,
		&i.Phone,
	)
	return i, err
}

const findBusinessClientByID = `-- name: FindBusinessClientByID :one
SELECT
    count_id, id, created_at, updated_at, deleted_at, business_id, fullname, email, phone
FROM
    client
WHERE
    id = $1
    AND business_id = $2
`

type FindBusinessClientByIDParams struct {
	ID         uuid.UUID `db:"id" json:"id"`
	BusinessID uuid.UUID `db:"business_id" json:"business_id"`
}

func (q *Queries) FindBusinessClientByID(ctx context.Context, arg FindBusinessClientByIDParams) (Client, error) {
	row := q.db.QueryRow(ctx, findBusinessClientByID, arg.ID, arg.BusinessID)
	var i Client
	err := row.Scan(
		&i.CountID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.BusinessID,
		&i.Fullname,
		&i.Email,
		&i.Phone,
	)
	return i, err
}

const findClientsWhere = `-- name: FindClientsWhere :many
SELECT
    count_id, id, created_at, updated_at, deleted_at, business_id, fullname, email, phone
FROM
    client
WHERE
    (
        client.id = $1 or $1 IS NULL
    )
    AND (
        client.business_id = $2 or $2 IS NULL
    )
    AND (
        client.fullname ilike $3 or $3 IS NULL
    )
    AND (
        client.email ilike $4
        or $4 IS NULL
    )
    AND (
        client.phone ilike $5
        or $5 IS NULL
    )
    AND (
        client.created_at <= $6
        or $6 IS NULL
    )
    AND (
        client.count_id < $7
        or $7 IS NULL
    )
ORDER BY
    client.created_at DESC,
    client.count_id DESC
LIMIT
    $8
`

type FindClientsWhereParams struct {
	ID         *uuid.UUID         `db:"id" json:"id"`
	BusinessID *uuid.UUID         `db:"business_id" json:"business_id"`
	Fullname   *string            `db:"fullname" json:"fullname"`
	Email      *string            `db:"email" json:"email"`
	Phone      *string            `db:"phone" json:"phone"`
	CursorTime pgtype.Timestamptz `db:"cursor_time" json:"cursor_time"`
	CursorID   *int64             `db:"cursor_id" json:"cursor_id"`
	Limit      *int32             `db:"limit" json:"limit"`
}

func (q *Queries) FindClientsWhere(ctx context.Context, arg FindClientsWhereParams) ([]Client, error) {
	rows, err := q.db.Query(ctx, findClientsWhere,
		arg.ID,
		arg.BusinessID,
		arg.Fullname,
		arg.Email,
		arg.Phone,
		arg.CursorTime,
		arg.CursorID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Client
	for rows.Next() {
		var i Client
		if err := rows.Scan(
			&i.CountID,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.BusinessID,
			&i.Fullname,
			&i.Email,
			&i.Phone,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateClient = `-- name: UpdateClient :one
UPDATE
    client
SET
    updated_at = timezone('utc', now()),
    fullname = COALESCE($2, fullname),
    email = COALESCE($3, email),
    phone = COALESCE($4, phone)
WHERE
    id = $1 RETURNING count_id, id, created_at, updated_at, deleted_at, business_id, fullname, email, phone
`

type UpdateClientParams struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Fullname *string   `db:"fullname" json:"fullname"`
	Email    *string   `db:"email" json:"email"`
	Phone    *string   `db:"phone" json:"phone"`
}

func (q *Queries) UpdateClient(ctx context.Context, arg UpdateClientParams) (Client, error) {
	row := q.db.QueryRow(ctx, updateClient,
		arg.ID,
		arg.Fullname,
		arg.Email,
		arg.Phone,
	)
	var i Client
	err := row.Scan(
		&i.CountID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.BusinessID,
		&i.Fullname,
		&i.Email,
		&i.Phone,
	)
	return i, err
}
