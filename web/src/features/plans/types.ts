import type {
  CanonicalProduct,
  Product,
  ProductVariant,
} from '@/features/products/types'

export interface Plan {
  id: string
  inventory_id: string
  title: string
  description?: string
  created_at: string
  updated_at: string
  deleted_at?: string

  // Fields from PlanWithDetails
  items?: PlanItem[]
  shopping_lists?: string[] // Array of Shopping List IDs
}

export interface PlanGroup {
  id: string
  inventory_id: string
  title: string
  description?: string
  created_at: string
  updated_at: string
  deleted_at?: string

  // Fields from PlanGroupWithDetails
  plans?: Plan[]
  shopping_lists?: string[]
}

export interface PlanItem {
  id: string
  plan_id: string
  target_type: 'canonical_product' | 'product' | 'product_variant'
  target_id: string
  quantity?: number
  unit?: string
  note?: string
  created_at: string
  updated_at: string
  deleted_at?: string

  // Expanded fields (if available from join or manual fetch)
  canonical_product?: CanonicalProduct
  product?: Product
  product_variant?: ProductVariant
}
