-- +goose Up
-- +goose StatementBegin
BEGIN;
ALTER TABLE couriers ALTER COLUMN status SET DEFAULT 'available';
ALTER TABLE couriers ADD COLUMN transport_type TEXT NOT NULL DEFAULT 'on_foot';
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;
ALTER TABLE couriers ALTER COLUMN status DROP DEFAULT;
ALTER TABLE couriers DROP COLUMN transport_type;
COMMIT;
-- +goose StatementEnd
