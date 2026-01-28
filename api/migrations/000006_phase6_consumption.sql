-- +goose Up
ALTER TABLE consumption_events ADD COLUMN source VARCHAR(50) NOT NULL DEFAULT 'manual';
ALTER TABLE consumption_events RENAME COLUMN quantity_consumed TO quantity;

-- +goose Down
ALTER TABLE consumption_events RENAME COLUMN quantity TO quantity_consumed;
ALTER TABLE consumption_events DROP COLUMN source;
