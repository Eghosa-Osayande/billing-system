-- +goose Up
ALTER TABLE
    invoice
ALTER COLUMN
    created_at
SET
    DEFAULT timezone('utc', now());

ALTER TABLE
    invoice
ALTER COLUMN
    shipping_fee TYPE DECIMAL(10, 2);


-- +goose Down
ALTER TABLE
    invoice
ALTER COLUMN
    created_at DROP DEFAULT;