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

