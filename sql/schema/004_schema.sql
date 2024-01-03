-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE invoice
ALTER COLUMN id SET DEFAULT uuid_generate_v4();


-- +goose Down
ALTER TABLE invoice
ALTER COLUMN id DROP DEFAULT;