import { api } from '@/lib/api'
import type {
  Inventory,
  InventoryMembership,
  InventoryProductDetail,
  InventorySettings,
} from './types'

export const getInventories = async (): Promise<Inventory[]> => {
  const result = await api<Inventory[]>('/inventories')
  return Array.isArray(result) ? result : []
}

export const getInventory = async (id: string): Promise<Inventory> => {
  return api<Inventory>(`/inventories/${id}`)
}

export const createInventory = async (name: string): Promise<Inventory> => {
  return api<Inventory>('/inventories', {
    method: 'POST',
    json: { name },
  })
}

export const getMembers = async (
  inventoryId: string
): Promise<InventoryMembership[]> => {
  return api<InventoryMembership[]>(`/inventories/${inventoryId}/members`)
}

export const getInventoryProducts = async (
  inventoryId: string
): Promise<InventoryProductDetail[]> => {
  return api<InventoryProductDetail[]>(
    `/inventories/${inventoryId}/inventory-products`
  )
}

export const getInventorySettings = async (
  inventoryId: string
): Promise<InventorySettings> => {
  return api<InventorySettings>(`/inventories/${inventoryId}/settings`)
}

export const updateInventorySettings = async (
  inventoryId: string,
  settings: Partial<InventorySettings>
): Promise<InventorySettings> => {
  return api<InventorySettings>(`/inventories/${inventoryId}/settings`, {
    method: 'PUT',
    json: settings,
  })
}
