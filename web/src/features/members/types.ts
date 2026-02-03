export type InventoryMembership = {
  id: string
  inventory_id: string
  user_id: string
  user_email?: string
  user_name?: string
  role: string
  invited_at: string
  deleted_at?: string
}

export type Invitation = {
  id: string
  inventory_id: string
  email: string
  role: string
  invited_by_user_id: string
  status: string
  token: string
  created_at: string
  accepted_at?: string
  expires_at?: string
}
