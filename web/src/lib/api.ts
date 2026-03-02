import { useAuthStore } from '@/store/auth'

export const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
export const AGENT_BASE_URL =
  import.meta.env.VITE_AGENT_API_URL || 'http://localhost:8081'

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

  // Default timeout of 15 seconds
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 15000)

  const config: RequestInit = {
    ...options,
    signal: options.signal || controller.signal,
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(json ? { 'Content-Type': 'application/json' } : {}),
      ...headers,
    },
    body: json ? JSON.stringify(json) : options.body,
  }

  let url = `${BASE_URL}${endpoint}`
  if (endpoint.startsWith('/agent')) {
    url = `${AGENT_BASE_URL}${endpoint.replace(/^\/agent/, '')}`
  }

  let response
  try {
    response = await fetch(url, config)
  } finally {
    clearTimeout(timeoutId)
  }

  if (!response.ok) {
    if (response.status === 401 && useAuthStore.getState().token) {
      useAuthStore.getState().setSessionExpired(true)
    }

    let errorMessage = 'An error occurred'
    let errorData
    try {
      errorData = await response.json()
      // adjust based on your API error structure
      if (
        typeof errorData === 'object' &&
        errorData !== null &&
        'error' in errorData
      ) {
        errorMessage = (errorData as { error: string }).error
      } else if (
        typeof errorData === 'object' &&
        errorData !== null &&
        'message' in errorData
      ) {
        errorMessage = (errorData as { message: string }).message
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
