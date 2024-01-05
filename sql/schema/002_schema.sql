-- +goose Up
alter table invoice add CONSTRAINT userclients CHECK(
	client.business_id=business_id
)


-- +goose Down
