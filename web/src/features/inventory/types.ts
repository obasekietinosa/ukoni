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
  quantity: number
  unit?: string
  created_at: string
  last_updated: string
  product_name: string
  brand?: string
  variant_name: string
  size?: number
  product_unit?: string
  canonical_product_id?: string
}

export interface InventorySettings {
  inventory_id: string
  llm_provider?: string
  llm_api_key?: string
  updated_at: string
}
