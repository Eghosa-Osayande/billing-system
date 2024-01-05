
-- name: FindUserById :one
Select * from users where id=$1 LIMIT 1;