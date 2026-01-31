import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { InventoryPage } from './routes/inventory-page'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

describe('InventoryPage', () => {
  beforeEach(() => {
    // Authenticated user
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
    // Set active inventory
    useInventoryStore.setState({ activeInventoryId: 'inv-1' })
  })

  it('renders inventory items', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <InventoryPage />
      </QueryClientProvider>
    )

    // Check loading state
    expect(screen.getByText(/loading inventory/i)).toBeInTheDocument()

    // Check content from handler (Tesco Whole Milk, 1L, qty 2)
    await waitFor(() => {
      expect(screen.getByText('Tesco Whole Milk')).toBeInTheDocument()
    })

    expect(screen.getByText('Tesco')).toBeInTheDocument()
    expect(screen.getByText(/1L/)).toBeInTheDocument()
    expect(screen.getByText('2')).toBeInTheDocument()
  })
})
