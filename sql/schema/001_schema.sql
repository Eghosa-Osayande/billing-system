-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SET TIMEZONE TO 'UTC';


Create table if not exists users (
	count_id BIGSERIAL  NOT NULL,
	id uuid primary key DEFAULT uuid_generate_v1mc() NOT NULL,
	created_at TIMESTAMPTZ DEFAULT timezone('utc', now()) NOT NULL,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	fullname varchar(255) NOT NULL,
	email varchar(255) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	email_verified boolean NOT NULL
);

Create table if not exists user_email_verifications (
	count_id BIGSERIAL  NOT NULL,
	email varchar(255) primary key NOT NULL UNIQUE,
	created_at TIMESTAMPTZ DEFAULT timezone('utc', now()) NOT NULL,
	code varchar(255) NOT NULL,
	expires_at TIMESTAMPTZ NOT NULL
);

Create table if not exists business (
	count_id BIGSERIAL  NOT NULL,
	id uuid primary key DEFAULT uuid_generate_v1mc() NOT NULL,
	created_at TIMESTAMPTZ DEFAULT timezone('utc', now()) NOT NULL,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	business_name varchar(255) NOT NULL,
	business_avatar varchar(255),
	owner_id uuid NOT NULL,
	FOREIGN KEY (owner_id) REFERENCES users(id),
	invoice_count DECIMAL(10,0) NOT NULL DEFAULT 0
);

Create table if not exists client (
	count_id BIGSERIAL  NOT NULL,
	id uuid primary key DEFAULT uuid_generate_v1mc() NOT NULL,
	created_at TIMESTAMPTZ DEFAULT timezone('utc', now()) NOT NULL,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	business_id uuid NOT NULL,
	FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE,
	fullname varchar(255) NOT NULL,
	email varchar(255) NULL,
	phone varchar(255) NULL
);




Create table if not exists invoice (
	count_id BIGSERIAL  NOT NULL,
	id uuid primary key DEFAULT uuid_generate_v1mc() NOT NULL,
	created_at TIMESTAMPTZ DEFAULT timezone('utc', now()) NOT NULL,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	business_id uuid NOT NULL,
	FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE,
	currency varchar(255) NULL,
	currency_symbol varchar(255) NULL,
	payment_due_date TIMESTAMPTZ NULL,
	date_of_issue TIMESTAMPTZ NULL,
	notes varchar(255) NULL,
	payment_method varchar(255) NULL,
	payment_status varchar(255) DEFAULT 'unpaid' NOT NULL,
	client_id uuid NULL,
	FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE SET NULL,
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
	),
	tax DECIMAL(10, 2) NULL,
	invoice_number VARCHAR(16) NOT NULL DEFAULT '-',
	total DECIMAL(10, 2) NULL
);

Create INDEX idx_invoice_pagination ON invoice (created_at, id);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_invoice_number()
RETURNS TRIGGER AS $$
DECLARE count INTEGER;
BEGIN
	count= (SELECT invoice_count FROM business WHERE id=NEW.business_id);

    UPDATE business SET invoice_count = count + 1 WHERE id = NEW.business_id;

    NEW.invoice_number = CONCAT('IN'::text, RIGHT( CONCAT('00000'::text , 
             to_hex(count)),5));
    RETURN NEW;

END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_invoice
BEFORE INSERT ON invoice
FOR EACH ROW
EXECUTE FUNCTION update_invoice_number();
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS invoiceitem (
	id uuid primary key DEFAULT uuid_generate_v1mc(),
	created_at TIMESTAMPTZ DEFAULT timezone('utc', now()),
	invoice_id uuid NOT NULL,
	FOREIGN KEY (invoice_id) REFERENCES invoice(id) ON DELETE CASCADE,
	title varchar(255) NOT NULL,
	price DECIMAL(10, 2) NOT NULL,
	quantity DECIMAL(10, 2) NOT NULL,
	discount DECIMAL(10, 2) NULL,
	discount_type varchar(255) NULL,
	CONSTRAINT check_discounttype_is_not_null_when_discount_is_not_null CHECK (
		(
			discount_type IS NULL
			AND discount IS NULL
		)
		OR (
			discount_type IS NOT NULL
			AND discount IS NOT NULL
		)
	)
);

-- +goose Down

DROP TABLE IF EXISTS invoiceitem;

DROP TABLE IF EXISTS invoice;

DROP TABLE IF EXISTS client;

DROP TABLE IF EXISTS business;

DROP TABLE IF EXISTS user_email_verifications;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";

DROP INDEX IF EXISTS idx_invoice_pagination;

Drop TYPE if exists invoice_payment_status;

DROP  TRIGGER IF EXISTS  before_insert_invoice ON invoice;

DROP FUNCTION update_invoice_number();


