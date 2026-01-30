import { useAuthStore } from '@/store/auth'

export const BASE_URL = import.meta.env.VITE_API_URL || 'https://api.ukoni.app'

type FetchOptions = RequestInit & {
  json?: unknown
}

export class ApiError extends Error {
  public status: number
  public data?: unknown

  constructor(status: number, message: string, data?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

export async function api<T = unknown>(
  endpoint: string,
  { json, headers, ...options }: FetchOptions = {}
): Promise<T> {
  const token = useAuthStore.getState().token

  const config: RequestInit = {
    ...options,
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(json ? { 'Content-Type': 'application/json' } : {}),
      ...headers,
    },
    body: json ? JSON.stringify(json) : options.body,
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, config)

  if (!response.ok) {
    let errorMessage = 'An error occurred'
    let errorData
    try {
      errorData = await response.json()
      // adjust based on your API error structure
      if (typeof errorData === 'object' && errorData !== null && 'error' in errorData) {
        errorMessage = (errorData as any).error
      } else if (typeof errorData === 'object' && errorData !== null && 'message' in errorData) {
        errorMessage = (errorData as any).message
      }
    } catch {
      // ignore JSON parse error
    }
    throw new ApiError(response.status, errorMessage, errorData)
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return {} as T
  }

  try {
    return await response.json()
  } catch {
    return {} as T
  }
}
