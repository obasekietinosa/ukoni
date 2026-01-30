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
