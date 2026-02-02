-- +goose Up
UPDATE products SET category_id = NULL;
UPDATE canonical_products SET category_id = NULL;
DELETE FROM product_categories;

ALTER TABLE product_categories ADD COLUMN inventory_id UUID NOT NULL REFERENCES inventories(id);

-- +goose Down
ALTER TABLE product_categories DROP COLUMN inventory_id;
