-- +goose Up
Create table if not exists client (
	id uuid primary key NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp,
	deleted_at timestamp,
	business_id uuid NOT NULL,
	FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE,
	fullname varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	phone varchar(255) NOT NULL
);

Create table if not exists invoice (
	id uuid primary key NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp,
	deleted_at timestamp,
	business_id uuid NOT NULL,
	FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE,
	currency varchar(255) NULL,
	payment_due_date timestamp NULL,
	date_of_issue timestamp NULL,
	notes varchar(255) NULL,
	payment_method varchar(255) NULL,
	payment_status varchar(255) NULL,
	items jsonb NULL,
	CHECK (items IS NULL OR (
            items -> 'name' IS NOT NULL AND
            items -> 'price' IS NOT NULL
			AND
            items -> 'quantity' IS NOT NULL
			AND
            items ? 'discount' 
        )),
	client_id uuid NULL,
	FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE SET NULL,
	shipping_fee_type varchar(255) NULL,
	shipping_fee numeric NULL,
	CONSTRAINT check_shippingfeetype_is_not_null_when_shippingfee_is_not_null CHECK (
		(
			shipping_fee_type IS NULL
			AND shipping_fee IS NULL
		)
		OR (
			shipping_fee_type IS NOT NULL
			AND shipping_fee IS NOT NULL
		)
	)
);


-- +goose Down

DROP TABLE IF EXISTS invoice_item;
DROP TABLE IF EXISTS invoice;
DROP TABLE IF EXISTS client;