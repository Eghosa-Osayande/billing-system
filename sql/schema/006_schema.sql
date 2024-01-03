-- +goose Up

CREATE TABLE IF NOT EXISTS invoiceitem (
    id uuid primary key DEFAULT uuid_generate_v4(),
    created_at timestamp DEFAULT timezone('utc', now()),
    invoice_id uuid NOT NULL,
    FOREIGN KEY (invoice_id) REFERENCES invoice(id) ON DELETE CASCADE,
    title varchar(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) NULL,
    discount_type varchar(255) NULL
);



