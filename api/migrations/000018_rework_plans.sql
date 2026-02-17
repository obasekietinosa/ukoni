-- +goose Up
-- Create plan_groups table
CREATE TABLE plan_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID NOT NULL REFERENCES inventories(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_plan_groups_inventory_id ON plan_groups(inventory_id);

-- Create plan_group_plans table for many-to-many relationship
CREATE TABLE plan_group_plans (
    plan_group_id UUID NOT NULL REFERENCES plan_groups(id),
    plan_id UUID NOT NULL REFERENCES plans(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (plan_group_id, plan_id)
);

CREATE INDEX idx_plan_group_plans_plan_id ON plan_group_plans(plan_id);

-- Create plan_group_shopping_lists table for many-to-many relationship
CREATE TABLE plan_group_shopping_lists (
    plan_group_id UUID NOT NULL REFERENCES plan_groups(id),
    shopping_list_id UUID NOT NULL REFERENCES shopping_lists(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (plan_group_id, shopping_list_id)
);

CREATE INDEX idx_plan_group_shopping_lists_shopping_list_id ON plan_group_shopping_lists(shopping_list_id);

-- Drop parent_plan_id from plans
ALTER TABLE plans DROP COLUMN parent_plan_id;

-- +goose Down
-- Re-add parent_plan_id to plans
ALTER TABLE plans ADD COLUMN parent_plan_id UUID REFERENCES plans(id);
CREATE INDEX idx_plans_parent_plan_id ON plans(parent_plan_id);

DROP TABLE IF EXISTS plan_group_shopping_lists;
DROP TABLE IF EXISTS plan_group_plans;
DROP TABLE IF EXISTS plan_groups;
