package handlers

import "github.com/jackc/pgx/v5"


func isErrNoRows (err error) bool {
	return err==pgx.ErrNoRows;
}