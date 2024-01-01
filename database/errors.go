package database


import (
	"github.com/jackc/pgx/v5")

func IsErrNoRows (err error) bool {
	return err==pgx.ErrNoRows;
}