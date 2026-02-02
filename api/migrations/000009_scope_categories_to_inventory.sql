-- +goose Up
DELETE FROM shopping_list_items;
DELETE FROM transaction_items;
DELETE FROM inventory_products;
DELETE FROM consumption_events;
DELETE FROM product_variants;
DELETE FROM products;
DELETE FROM canonical_products;
DELETE FROM product_categories;

ALTER TABLE product_categories ADD COLUMN inventory_id UUID NOT NULL REFERENCES inventories(id);

-- +goose Down
ALTER TABLE product_categories DROP COLUMN inventory_id;
