-- +goose Up
ALTER TABLE canonical_products ADD COLUMN category_id UUID REFERENCES product_categories(id);
ALTER TABLE canonical_products ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE canonical_products DROP COLUMN updated_at;
ALTER TABLE canonical_products DROP COLUMN category_id;
