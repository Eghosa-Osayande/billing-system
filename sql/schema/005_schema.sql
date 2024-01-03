-- +goose Up


ALTER TABLE invoice
ALTER COLUMN created_at SET DEFAULT timezone('utc', now());


-- +goose Down
ALTER TABLE invoice
ALTER COLUMN created_at DROP DEFAULT;