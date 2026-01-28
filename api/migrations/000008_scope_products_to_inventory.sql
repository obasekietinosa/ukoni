-- +goose Up
DELETE FROM shopping_list_items;
DELETE FROM transaction_items;
DELETE FROM inventory_products;
DELETE FROM consumption_events;
DELETE FROM product_variants;
DELETE FROM products;
DELETE FROM canonical_products;

ALTER TABLE canonical_products ADD COLUMN inventory_id UUID NOT NULL REFERENCES inventories(id);
ALTER TABLE products ADD COLUMN inventory_id UUID NOT NULL REFERENCES inventories(id);

-- +goose Down
ALTER TABLE products DROP COLUMN inventory_id;
ALTER TABLE canonical_products DROP COLUMN inventory_id;
