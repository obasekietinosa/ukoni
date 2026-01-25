-- +goose Up
ALTER TABLE shopping_lists ADD COLUMN created_by UUID REFERENCES users(id);
ALTER TABLE shopping_list_items ADD COLUMN notes TEXT;
ALTER TABLE transaction_items ADD COLUMN shopping_list_item_id UUID REFERENCES shopping_list_items(id);

-- +goose Down
ALTER TABLE transaction_items DROP COLUMN shopping_list_item_id;
ALTER TABLE shopping_list_items DROP COLUMN notes;
ALTER TABLE shopping_lists DROP COLUMN created_by;
