-- +goose Up
ALTER TABLE shopping_list_items ADD COLUMN unit VARCHAR(100);
ALTER TABLE shopping_list_items ALTER COLUMN quantity TYPE DECIMAL;

-- +goose Down
ALTER TABLE shopping_list_items ALTER COLUMN quantity TYPE DECIMAL(10, 2);
ALTER TABLE shopping_list_items DROP COLUMN unit;
