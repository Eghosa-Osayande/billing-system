-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


Create table if not exists users (
	id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamp DEFAULT timezone('utc', now()) NOT NULL,
	updated_at timestamp,
	deleted_at timestamp,
	fullname varchar(255) NOT NULL,
	email varchar(255) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	email_verified boolean NOT NULL
);

Create table if not exists user_email_verifications (
	email varchar(255) primary key NOT NULL UNIQUE,
	created_at timestamp DEFAULT timezone('utc', now()) NOT NULL,
	code varchar(255) NOT NULL,
	expires_at timestamp NOT NULL
);

Create table if not exists business (
	id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamp DEFAULT timezone('utc', now()) NOT NULL,
	updated_at timestamp,
	deleted_at timestamp,
	business_name varchar(255) NOT NULL,
	business_avatar varchar(255),
	owner_id uuid NOT NULL,
	FOREIGN KEY (owner_id) REFERENCES users(id)
);

Create table if not exists client (
	id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamp DEFAULT timezone('utc', now()) NOT NULL,
	updated_at timestamp,
	deleted_at timestamp,
	business_id uuid NOT NULL,
	FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE,
	fullname varchar(255) NOT NULL,
	email varchar(255) NULL,
	phone varchar(255) NULL
);

Create table if not exists invoice (
	id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamp DEFAULT timezone('utc', now()) NOT NULL,
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
	client_id uuid NULL,
	FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE
	SET
		NULL,
		shipping_fee_type varchar(255) NULL,
		shipping_fee DECIMAL(10, 2) NULL,
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

Create INDEX idx_invoice_pagination ON invoice (created_at, id);

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

-- +goose Down
DROP TABLE IF EXISTS invoiceitem;

DROP TABLE IF EXISTS invoice;

DROP TABLE IF EXISTS client;

DROP TABLE IF EXISTS business;

DROP TABLE IF EXISTS user_email_verifications;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";