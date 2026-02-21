import { api } from '@/lib/api'

export type ActionResult = {
  tool_name: string
  data: unknown
}

export type ChatResponse = {
  response: string
  error?: string
  actions?: ActionResult[]
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
