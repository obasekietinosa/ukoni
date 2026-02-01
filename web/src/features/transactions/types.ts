export interface TransactionItem {
  id: string
  transaction_id: string
  product_variant_id: string
  quantity: number
  price_per_unit?: number
  shopping_list_item_id?: string
}

export interface Transaction {
  id: string
  inventory_id: string
  outlet_id?: string
  created_by_user_id: string
  transaction_date: string
  total_amount?: number
  created_at: string
  items?: TransactionItem[]
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
