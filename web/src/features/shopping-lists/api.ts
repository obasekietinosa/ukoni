import { api } from '@/lib/api'
import type { ShoppingList, ShoppingListItem } from './types'

export const getShoppingLists = async (
  inventoryId: string
): Promise<ShoppingList[]> => {
  return api<ShoppingList[]>(`/inventories/${inventoryId}/shopping-lists`)
}

export const createShoppingList = async (
  inventoryId: string,
  name: string
): Promise<ShoppingList> => {
  return api<ShoppingList>(`/inventories/${inventoryId}/shopping-lists`, {
    method: 'POST',
    json: { name },
  })
}

export const getShoppingList = async (
  listId: string
): Promise<ShoppingList> => {
  return api<ShoppingList>(`/shopping-lists/${listId}`)
}

export const updateShoppingList = async (
  listId: string,
  name: string
): Promise<ShoppingList> => {
  return api<ShoppingList>(`/shopping-lists/${listId}`, {
    method: 'PUT',
    json: { name },
  })
}

export const deleteShoppingList = async (listId: string): Promise<void> => {
  return api<void>(`/shopping-lists/${listId}`, {
    method: 'DELETE',
  })
}

export const getShoppingListItems = async (
  listId: string
): Promise<ShoppingListItem[]> => {
  return api<ShoppingListItem[]>(`/shopping-lists/${listId}/items`)
}

export const addShoppingListItem = async (
  listId: string,
  data: {
    target_type: 'canonical_product' | 'product_variant'
    target_id: string
    preferred_outlet_id?: string
    notes?: string
    quantity?: number
    unit?: string
  }
): Promise<ShoppingListItem> => {
  return api<ShoppingListItem>(`/shopping-lists/${listId}/items`, {
    method: 'POST',
    json: data,
  })
}

export const updateShoppingListItem = async (
  itemId: string,
  data: {
    notes?: string
    preferred_outlet_id?: string
    quantity?: number
    unit?: string
  }
): Promise<ShoppingListItem> => {
  return api<ShoppingListItem>(`/shopping-list-items/${itemId}`, {
    method: 'PUT',
    json: data,
  })
}

export const deleteShoppingListItem = async (itemId: string): Promise<void> => {
  return api<void>(`/shopping-list-items/${itemId}`, {
    method: 'DELETE',
  })
}
