import { api } from '@/lib/api'
import { Invitation } from './types'
import { getMembers } from '@/features/inventory/api'

export { getMembers }

export const inviteUser = async (
  inventoryId: string,
  email: string,
  role: string
): Promise<Invitation> => {
  return api<Invitation>(`/inventories/${inventoryId}/invitations`, {
    method: 'POST',
    json: { email, role },
  })
}

export const removeMember = async (
  inventoryId: string,
  userId: string
): Promise<void> => {
  return api<void>(`/inventories/${inventoryId}/members/${userId}`, {
    method: 'DELETE',
  })
}

export const updateMemberRole = async (
  inventoryId: string,
  userId: string,
  role: string
): Promise<void> => {
  return api<void>(`/inventories/${inventoryId}/members/${userId}`, {
    method: 'PUT',
    json: { role },
  })
}

export const acceptInvite = async (
  inviteId: string,
  token: string
): Promise<void> => {
  return api<void>(`/invitations/${inviteId}/accept`, {
    method: 'POST',
    json: { token },
  })
}
