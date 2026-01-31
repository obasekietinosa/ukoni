import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useCurrentUserRole } from './use-current-user-role'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
)

describe('useCurrentUserRole', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
  })

  it('returns the role for the current user', async () => {
    const { result } = renderHook(() => useCurrentUserRole('inv-1'), { wrapper })

    await waitFor(() => {
      expect(result.current.role).toBe('admin')
    })
    expect(result.current.isLoading).toBe(false)
  })

  it('returns undefined if user is not found in members', async () => {
    useAuthStore.setState({
      user: { id: '3', name: 'Other User', email: 'other@example.com' },
      token: 'fake-token',
    })

    const { result } = renderHook(() => useCurrentUserRole('inv-1'), { wrapper })

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.role).toBeUndefined()
  })
})
