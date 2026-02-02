import type {
  CanonicalProduct,
  Product,
  ProductVariant,
} from '../products/types'
import type { Outlet } from '../sellers/types'

export interface ShoppingList {
  id: string
  inventory_id: string
  name: string
  created_by: string
  created_at: string
  last_updated_at: string
  deleted_at?: string
}

export interface ShoppingListItem {
  id: string
  shopping_list_id: string
  target_type: 'canonical_product' | 'product_variant'
  target_id: string
  preferred_outlet_id?: string
  notes?: string
  quantity?: number
  unit?: string
  created_at: string
  deleted_at?: string

  // Joined fields
  canonical_product?: CanonicalProduct
  product_variant?: ProductVariant
  product?: Product // If product_variant is set, this might be populated
  preferred_outlet?: Outlet
}
