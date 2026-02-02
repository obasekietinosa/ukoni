-- +goose Up
ALTER TABLE shopping_list_items ADD COLUMN quantity DECIMAL(10, 2) DEFAULT 1;

-- +goose Down
ALTER TABLE shopping_list_items DROP COLUMN quantity;
