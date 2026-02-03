import { InventoryMembership } from '@/features/inventory/types'

export type { InventoryMembership }

export interface Invitation {
  id: string
  inventory_id: string
  email: string
  role: string
  invited_by_user_id: string
  status: string
  token: string
  created_at: string
  expires_at?: string
}
