import { api } from '@/lib/api'
import type { CreateTransactionRequest, Transaction } from './types'

export const createTransaction = async (
  inventoryId: string,
  data: CreateTransactionRequest
): Promise<Transaction> => {
  return api<Transaction>(`/inventories/${inventoryId}/transactions`, {
    method: 'POST',
    json: data,
  })
}

export const getTransactions = async (
  inventoryId: string,
  limit?: number,
  offset?: number
): Promise<Transaction[]> => {
  const params = new URLSearchParams()
  if (limit) params.append('limit', limit.toString())
  if (offset) params.append('offset', offset.toString())

  return api<Transaction[]>(
    `/inventories/${inventoryId}/transactions?${params.toString()}`
  )
}

export const getTransaction = async (id: string): Promise<Transaction> => {
  return api<Transaction>(`/transactions/${id}`)
}
