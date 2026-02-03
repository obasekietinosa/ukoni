export interface ConsumptionEvent {
  id: string
  inventory_id: string
  canonical_product_id?: string
  product_variant_id?: string
  created_by_user_id?: string
  quantity?: number
  unit?: string
  note?: string
  source: string
  consumed_at: string
  deleted_at?: string
}

export interface ConsumptionEventDetail extends ConsumptionEvent {
  canonical_product_name?: string
  variant_name?: string
  product_name?: string
  brand?: string
  size?: number
  product_unit?: string
}

export interface CreateConsumptionInput {
  canonical_product_id?: string
  product_variant_id?: string
  quantity?: number
  unit?: string
  note?: string
  source?: string
  consumed_at?: string
}
