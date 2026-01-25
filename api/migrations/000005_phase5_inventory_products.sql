-- +goose Up
ALTER TABLE inventory_products ADD COLUMN unit VARCHAR(100);
ALTER TABLE inventory_products ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
CREATE UNIQUE INDEX uk_inventory_products_active ON inventory_products (inventory_id, product_variant_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS uk_inventory_products_active;
ALTER TABLE inventory_products DROP COLUMN created_at;
ALTER TABLE inventory_products DROP COLUMN unit;
