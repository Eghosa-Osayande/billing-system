-- +goose Up
Create table if not exists users (
	id uuid primary key  NOT NULL,
	created_at timestamp  NOT NULL,
	updated_at timestamp,
	deleted_at timestamp,
	fullname varchar(255)  NOT NULL,
	email varchar(255) NOT NULL UNIQUE,
	phone varchar(255) NOT NULL,
	password varchar(255) NOT NULL,
	email_verified boolean  NOT NULL
);

Create table if not exists user_email_verifications (
	email varchar(255) primary key  NOT NULL UNIQUE,
	created_at timestamp  NOT NULL,
	code varchar(255)  NOT NULL,
	expires_at timestamp  NOT NULL
);

Create table if not exists business (
	id uuid primary key NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp,
	deleted_at timestamp,
	business_name varchar(255) NOT NULL,
	business_avatar varchar(255),
	owner_id uuid NOT NULL,
	FOREIGN KEY (owner_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE IF EXISTS business;
DROP TABLE IF EXISTS user_email_verifications;
DROP TABLE IF EXISTS users;

