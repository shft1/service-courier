-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_delivery_courier_id ON delivery(courier_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_delivery_courier_id;
-- +goose StatementEnd
