-- +goose Up
-- Create plans table
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID NOT NULL REFERENCES inventories(id),
    parent_plan_id UUID REFERENCES plans(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_plans_inventory_id ON plans(inventory_id);
CREATE INDEX idx_plans_parent_plan_id ON plans(parent_plan_id);

-- Create plan_items table
CREATE TABLE plan_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id),
    target_type VARCHAR(50) NOT NULL CHECK (target_type IN ('canonical_product', 'product', 'product_variant')),
    target_id UUID NOT NULL,
    quantity DECIMAL,
    unit VARCHAR(50),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_plan_items_plan_id ON plan_items(plan_id);
CREATE INDEX idx_plan_items_target ON plan_items(target_type, target_id);

-- Create plan_shopping_lists table for many-to-many relationship
CREATE TABLE plan_shopping_lists (
    plan_id UUID NOT NULL REFERENCES plans(id),
    shopping_list_id UUID NOT NULL REFERENCES shopping_lists(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (plan_id, shopping_list_id)
);

CREATE INDEX idx_plan_shopping_lists_shopping_list_id ON plan_shopping_lists(shopping_list_id);

-- +goose Down
DROP TABLE IF EXISTS plan_shopping_lists;
DROP TABLE IF EXISTS plan_items;
DROP TABLE IF EXISTS plans;
