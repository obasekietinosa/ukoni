import { api } from '@/lib/api'
import type {
  Inventory,
  InventoryMembership,
  InventoryProductDetail,
} from './types'

export const getInventories = async (): Promise<Inventory[]> => {
  return api<Inventory[]>('/inventories')
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
