-- +goose Up
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