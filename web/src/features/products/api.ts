import { api } from '@/lib/api'
import type { CanonicalProduct } from './types'

export const getCanonicalProducts = async (
  inventoryId: string
): Promise<CanonicalProduct[]> => {
  return api<CanonicalProduct[]>(
    `/inventories/${inventoryId}/canonical-products`
  )
}

export const createCanonicalProduct = async (
  inventoryId: string,
  data: { name: string; description?: string; category_id?: string }
): Promise<CanonicalProduct> => {
  return api<CanonicalProduct>(
    `/inventories/${inventoryId}/canonical-products`,
    {
      method: 'POST',
      json: data,
    }
  )
}
