export interface CanonicalProduct {
  id: string
  inventory_id: string
  name: string
  description?: string
  category_id?: string
  created_at: string
  updated_at?: string
  deleted_at?: string
}

export interface Product {
  id: string
  inventory_id: string
  canonical_product_id?: string
  brand?: string
  name: string
  description?: string
  category_id?: string
  created_at: string
  deleted_at?: string
}

export interface ProductVariant {
  id: string
  product_id: string
  variant_name: string
  sku?: string
  unit?: string
  size?: number
  deleted_at?: string
}
