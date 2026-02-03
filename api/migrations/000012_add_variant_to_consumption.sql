-- +goose Up
ALTER TABLE consumption_events ADD COLUMN product_variant_id UUID REFERENCES product_variants(id);

-- +goose Down
ALTER TABLE consumption_events DROP COLUMN product_variant_id;
