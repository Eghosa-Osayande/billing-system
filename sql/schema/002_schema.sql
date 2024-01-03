-- +goose Up
ALTER TABLE users DROP COLUMN phone;

-- +goose Down
ALTER TABLE users ADD COLUMN phone varchar(255) NOT NULL;

