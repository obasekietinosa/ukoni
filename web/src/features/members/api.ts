import { api } from '@/lib/api'
import type { Invitation, InventoryMembership } from './types'

export const listMembers = (inventoryId: string) => {
  return api<InventoryMembership[]>(`/inventories/${inventoryId}/members`)
}

export const inviteUser = (
  inventoryId: string,
  data: { email: string; role: string }
) => {
  return api<Invitation>(`/inventories/${inventoryId}/invitations`, {
    method: 'POST',
    json: data,
  })
}

export const removeMember = (inventoryId: string, userId: string) => {
  return api<void>(`/inventories/${inventoryId}/members/${userId}`, {
    method: 'DELETE',
  })
}

export const updateMemberRole = (
  inventoryId: string,
  userId: string,
  role: string
) => {
  return api<void>(`/inventories/${inventoryId}/members/${userId}`, {
    method: 'PUT',
    json: { role },
  })
}

export const acceptInvite = (inviteId: string, token: string) => {
  return api<{ status: string }>(`/invitations/${inviteId}/accept`, {
    method: 'POST',
    json: { token },
  })
}
