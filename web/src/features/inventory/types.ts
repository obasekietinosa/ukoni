export interface Inventory {
  id: string
  name: string
  owner_user_id: string
  created_at: string
  deleted_at?: string
}

export interface InventoryMembership {
  id: string
  inventory_id: string
  user_id: string
  role: string
  invited_at: string
  deleted_at?: string
}

export interface InventoryProductDetail {
  id: string
  inventory_id: string
  product_variant_id: string
  canonical_product_name: string
  brand_name?: string
  variant_name: string
  quantity: number
  unit?: string
  sku?: string
  last_updated: string
}

export interface Transaction {
  id: string
  inventory_id: string
  outlet_id?: string
  created_by_user_id: string
  transaction_date: string
  total_amount?: number
  deleted_at?: string
}

export interface CreateTransactionItemRequest {
  product_variant_id: string
  quantity: number
  price_per_unit?: number
  shopping_list_item_id?: string
}

export interface CreateTransactionRequest {
  outlet_id?: string
  transaction_date: string
  items: CreateTransactionItemRequest[]
}
