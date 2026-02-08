-- +goose Up
-- Inventories
CREATE INDEX idx_inventories_owner_user_id ON inventories(owner_user_id);

-- Inventory Memberships
CREATE INDEX idx_inventory_memberships_inventory_id ON inventory_memberships(inventory_id);
CREATE INDEX idx_inventory_memberships_user_id ON inventory_memberships(user_id);

-- Invitations
CREATE INDEX idx_invitations_inventory_id ON invitations(inventory_id);
CREATE INDEX idx_invitations_invited_by_user_id ON invitations(invited_by_user_id);

-- Canonical Products
CREATE INDEX idx_canonical_products_inventory_id ON canonical_products(inventory_id);

-- Product Categories
CREATE INDEX idx_product_categories_inventory_id ON product_categories(inventory_id);
CREATE INDEX idx_product_categories_parent_category_id ON product_categories(parent_category_id);

-- Products
CREATE INDEX idx_products_canonical_product_id ON products(canonical_product_id);
CREATE INDEX idx_products_inventory_id ON products(inventory_id);
CREATE INDEX idx_products_category_id ON products(category_id);

-- Product Variants
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);

-- Outlets
CREATE INDEX idx_outlets_seller_id ON outlets(seller_id);

-- Inventory Products
CREATE INDEX idx_inventory_products_inventory_id ON inventory_products(inventory_id);
CREATE INDEX idx_inventory_products_product_variant_id ON inventory_products(product_variant_id);

-- Transactions
CREATE INDEX idx_transactions_inventory_id ON transactions(inventory_id);
CREATE INDEX idx_transactions_outlet_id ON transactions(outlet_id);
CREATE INDEX idx_transactions_created_by_user_id ON transactions(created_by_user_id);

-- Transaction Items (from 000004)
CREATE INDEX idx_transaction_items_shopping_list_item_id ON transaction_items(shopping_list_item_id);

-- Consumption Events
CREATE INDEX idx_consumption_events_inventory_id ON consumption_events(inventory_id);
CREATE INDEX idx_consumption_events_canonical_product_id ON consumption_events(canonical_product_id);
CREATE INDEX idx_consumption_events_created_by_user_id ON consumption_events(created_by_user_id);
CREATE INDEX idx_consumption_events_product_variant_id ON consumption_events(product_variant_id);

-- Activity Logs
CREATE INDEX idx_activity_logs_inventory_id ON activity_logs(inventory_id);
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);

-- Shopping Lists
CREATE INDEX idx_shopping_lists_inventory_id ON shopping_lists(inventory_id);
CREATE INDEX idx_shopping_lists_created_by ON shopping_lists(created_by);

-- Shopping List Items
CREATE INDEX idx_shopping_list_items_shopping_list_id ON shopping_list_items(shopping_list_id);
CREATE INDEX idx_shopping_list_items_preferred_outlet_id ON shopping_list_items(preferred_outlet_id);


-- +goose Down
DROP INDEX IF EXISTS idx_shopping_list_items_preferred_outlet_id;
DROP INDEX IF EXISTS idx_shopping_list_items_shopping_list_id;
DROP INDEX IF EXISTS idx_shopping_lists_created_by;
DROP INDEX IF EXISTS idx_shopping_lists_inventory_id;
DROP INDEX IF EXISTS idx_activity_logs_user_id;
DROP INDEX IF EXISTS idx_activity_logs_inventory_id;
DROP INDEX IF EXISTS idx_consumption_events_product_variant_id;
DROP INDEX IF EXISTS idx_consumption_events_created_by_user_id;
DROP INDEX IF EXISTS idx_consumption_events_canonical_product_id;
DROP INDEX IF EXISTS idx_consumption_events_inventory_id;
DROP INDEX IF EXISTS idx_transaction_items_shopping_list_item_id;
DROP INDEX IF EXISTS idx_transactions_created_by_user_id;
DROP INDEX IF EXISTS idx_transactions_outlet_id;
DROP INDEX IF EXISTS idx_transactions_inventory_id;
DROP INDEX IF EXISTS idx_inventory_products_product_variant_id;
DROP INDEX IF EXISTS idx_inventory_products_inventory_id;
DROP INDEX IF EXISTS idx_outlets_seller_id;
DROP INDEX IF EXISTS idx_product_variants_product_id;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_inventory_id;
DROP INDEX IF EXISTS idx_products_canonical_product_id;
DROP INDEX IF EXISTS idx_product_categories_parent_category_id;
DROP INDEX IF EXISTS idx_product_categories_inventory_id;
DROP INDEX IF EXISTS idx_canonical_products_inventory_id;
DROP INDEX IF EXISTS idx_invitations_invited_by_user_id;
DROP INDEX IF EXISTS idx_invitations_inventory_id;
DROP INDEX IF EXISTS idx_inventory_memberships_user_id;
DROP INDEX IF EXISTS idx_inventory_memberships_inventory_id;
DROP INDEX IF EXISTS idx_inventories_owner_user_id;
