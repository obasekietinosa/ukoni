import { api } from '@/lib/api'
import type { Plan, PlanGroup, PlanItem } from './types'
import type { ShoppingList } from '../shopping-lists/types'

// Plan Functions

export const getPlans = async (
  inventoryId: string,
  params?: {
    limit?: number
    offset?: number
  }
): Promise<Plan[]> => {
  const query = new URLSearchParams()
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

// Plan Group Functions

export const getPlanGroups = async (
  inventoryId: string,
  params?: {
    limit?: number
    offset?: number
  }
): Promise<PlanGroup[]> => {
  const query = new URLSearchParams()
  if (params?.limit) query.append('limit', params.limit.toString())
  if (params?.offset) query.append('offset', params.offset.toString())

  return api<PlanGroup[]>(`/inventories/${inventoryId}/plan-groups?${query.toString()}`)
}

export const getPlanGroup = async (id: string): Promise<PlanGroup> => {
  return api<PlanGroup>(`/plan-groups/${id}`)
}

export const createPlanGroup = async (
  inventoryId: string,
  data: {
    title: string
    description?: string
  }
): Promise<PlanGroup> => {
  return api<PlanGroup>(`/inventories/${inventoryId}/plan-groups`, {
    method: 'POST',
    json: data,
  })
}

export const updatePlanGroup = async (
  id: string,
  data: {
    title: string
    description?: string
  }
): Promise<PlanGroup> => {
  return api<PlanGroup>(`/plan-groups/${id}`, {
    method: 'PUT',
    json: data,
  })
}

export const deletePlanGroup = async (id: string): Promise<void> => {
  return api<void>(`/plan-groups/${id}`, {
    method: 'DELETE',
  })
}

export const addPlanToGroup = async (
  groupId: string,
  planId: string
): Promise<void> => {
  return api<void>(`/plan-groups/${groupId}/plans`, {
    method: 'POST',
    json: { plan_id: planId },
  })
}

export const removePlanFromGroup = async (
  groupId: string,
  planId: string
): Promise<void> => {
  return api<void>(`/plan-groups/${groupId}/plans/${planId}`, {
    method: 'DELETE',
  })
}

export const createShoppingListFromGroup = async (
  groupId: string,
  name?: string
): Promise<ShoppingList> => {
  return api<ShoppingList>(`/plan-groups/${groupId}/shopping-list`, {
    method: 'POST',
    json: { name },
  })
}

export const linkShoppingListToGroup = async (
  groupId: string,
  shoppingListId: string
): Promise<void> => {
  return api<void>(`/plan-groups/${groupId}/shopping-lists`, {
    method: 'POST',
    json: { shopping_list_id: shoppingListId },
  })
}

export const unlinkShoppingListFromGroup = async (
  groupId: string,
  shoppingListId: string
): Promise<void> => {
  return api<void>(`/plan-groups/${groupId}/shopping-lists/${shoppingListId}`, {
    method: 'DELETE',
  })
}
