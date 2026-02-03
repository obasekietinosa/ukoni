import { api } from '@/lib/api'
import type {
  ConsumptionEvent,
  ConsumptionEventDetail,
  CreateConsumptionInput,
} from './types'

export const createConsumptionEvent = async (
  inventoryId: string,
  data: CreateConsumptionInput
): Promise<ConsumptionEvent> => {
  return api<ConsumptionEvent>(
    `/inventories/${inventoryId}/consumption-events`,
    {
      method: 'POST',
      json: data,
    }
  )
}

export const getConsumptionEvents = async (
  inventoryId: string,
  options: {
    limit?: number
    offset?: number
  } = {}
): Promise<ConsumptionEventDetail[]> => {
  const params = new URLSearchParams()
  if (options.limit) {
    params.append('limit', options.limit.toString())
  }
  if (options.offset) {
    params.append('offset', options.offset.toString())
  }
  return api<ConsumptionEventDetail[]>(
    `/inventories/${inventoryId}/consumption-events?${params.toString()}`
  )
}
