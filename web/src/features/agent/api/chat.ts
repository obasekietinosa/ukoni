import { api } from '@/lib/api'

export type ChatResponse = {
  response: string
  error?: string
}

export type ChatRequest = {
  prompt: string
  inventory_id: string
}

export const sendMessage = (data: ChatRequest): Promise<ChatResponse> => {
  return api('/agent/chat', {
    method: 'POST',
    json: data,
  })
}
