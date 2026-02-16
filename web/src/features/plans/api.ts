import { api } from '@/lib/api'
import type { Plan, PlanItem } from './types'
import type { ShoppingList } from '../shopping-lists/types'

export const getPlans = async (
  inventoryId: string,
  params?: {
    parent_plan_id?: string
    limit?: number
    offset?: number
  }
): Promise<Plan[]> => {
  const query = new URLSearchParams()
  if (params?.parent_plan_id)
    query.append('parent_plan_id', params.parent_plan_id)
  if (params?.limit) query.append('limit', params.limit.toString())
  if (params?.offset) query.append('offset', params.offset.toString())

  return api<Plan[]>(`/inventories/${inventoryId}/plans?${query.toString()}`)
}

export const getPlan = async (id: string): Promise<Plan> => {
  return api<Plan>(`/plans/${id}`)
}

export const createPlan = async (
  inventoryId: string,
  data: {
    title: string
    description?: string
    parent_plan_id?: string
  }
): Promise<Plan> => {
  return api<Plan>(`/inventories/${inventoryId}/plans`, {
    method: 'POST',
    json: data,
  })
}

export const updatePlan = async (
  id: string,
  data: {
    title: string
    description?: string
    parent_plan_id?: string
  }
): Promise<Plan> => {
  return api<Plan>(`/plans/${id}`, {
    method: 'PUT',
    json: data,
  })
}

export const deletePlan = async (id: string): Promise<void> => {
  return api<void>(`/plans/${id}`, {
    method: 'DELETE',
  })
}

export const addPlanItem = async (
  planId: string,
  data: {
    target_type: 'canonical_product' | 'product' | 'product_variant'
    target_id: string
    quantity?: number
    unit?: string
    note?: string
  }
): Promise<PlanItem> => {
  return api<PlanItem>(`/plans/${planId}/items`, {
    method: 'POST',
    json: data,
  })
}

export const updatePlanItem = async (
  itemId: string,
  data: {
    quantity?: number
    unit?: string
    note?: string
  }
): Promise<PlanItem> => {
  return api<PlanItem>(`/plan-items/${itemId}`, {
    method: 'PUT',
    json: data,
  })
}

export const removePlanItem = async (itemId: string): Promise<void> => {
  return api<void>(`/plan-items/${itemId}`, {
    method: 'DELETE',
  })
}

export const linkShoppingList = async (
  planId: string,
  shoppingListId: string
): Promise<void> => {
  return api<void>(`/plans/${planId}/shopping-lists`, {
    method: 'POST',
    json: { shopping_list_id: shoppingListId },
  })
}

export const unlinkShoppingList = async (
  planId: string,
  shoppingListId: string
): Promise<void> => {
  return api<void>(`/plans/${planId}/shopping-lists/${shoppingListId}`, {
    method: 'DELETE',
  })
}

export const createShoppingListFromPlan = async (
  planId: string,
  name?: string
): Promise<ShoppingList> => {
  return api<ShoppingList>(`/plans/${planId}/shopping-list`, {
    method: 'POST',
    json: { name },
  })
}
