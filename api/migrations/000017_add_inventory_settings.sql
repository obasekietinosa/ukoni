-- +goose Up
CREATE TABLE IF NOT EXISTS inventory_settings (
    inventory_id UUID PRIMARY KEY REFERENCES inventories(id) ON DELETE CASCADE,
    llm_provider VARCHAR(50),
    llm_api_key TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS inventory_settings;
