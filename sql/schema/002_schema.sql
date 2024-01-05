-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION check_invoice_business_id()
RETURNS TRIGGER AS 
$$
BEGIN
	IF (
        NEW.client_id IS NOT NULL AND
        NEW.business_id IS NOT NULL AND
        NEW.business_id != (
            SELECT business_id
            FROM client
            WHERE id = NEW.client_id
        )
    ) THEN
        RAISE EXCEPTION 'Business_id in invoice must match business_id in client';
    END IF;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER check_invoice_business_id_trigger
BEFORE INSERT OR UPDATE
ON invoice
FOR EACH ROW
EXECUTE FUNCTION check_invoice_business_id();

-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS check_invoice_business_id_trigger ON invoice;
DROP FUNCTION IF EXISTS check_invoice_business_id();
