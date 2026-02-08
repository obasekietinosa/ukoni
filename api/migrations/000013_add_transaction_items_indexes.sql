-- +goose Up
CREATE INDEX idx_transaction_items_transaction_id ON transaction_items(transaction_id);
CREATE INDEX idx_transaction_items_product_variant_id ON transaction_items(product_variant_id);

-- +goose Down
DROP INDEX IF EXISTS idx_transaction_items_product_variant_id;
DROP INDEX IF EXISTS idx_transaction_items_transaction_id;
