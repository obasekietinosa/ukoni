-- +goose Up
ALTER TABLE shopping_list_items ADD COLUMN manual_order INT;

-- +goose Down
ALTER TABLE shopping_list_items DROP COLUMN manual_order;