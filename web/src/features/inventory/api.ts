import { api } from '@/lib/api'
import type { Inventory } from './types'

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
