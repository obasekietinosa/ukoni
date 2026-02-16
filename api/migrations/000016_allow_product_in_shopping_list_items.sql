-- +goose Up
ALTER TABLE shopping_list_items DROP CONSTRAINT shopping_list_items_target_type_check;
ALTER TABLE shopping_list_items ADD CONSTRAINT shopping_list_items_target_type_check CHECK (target_type IN ('canonical_product', 'product', 'product_variant'));

-- +goose Down
-- Note: This might fail if there are 'product' items in the table
ALTER TABLE shopping_list_items DROP CONSTRAINT shopping_list_items_target_type_check;
ALTER TABLE shopping_list_items ADD CONSTRAINT shopping_list_items_target_type_check CHECK (target_type IN ('canonical_product', 'product_variant'));
