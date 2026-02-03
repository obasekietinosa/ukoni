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
    // We use getAllByText because the item is rendered twice (desktop table and mobile card)
    await waitFor(() => {
      const items = screen.getAllByText('Tesco Whole Milk')
      expect(items.length).toBeGreaterThan(0)
      expect(items[0]).toBeInTheDocument()
    })

    // Check other details (also likely duplicated)
    expect(screen.getAllByText('Tesco')[0]).toBeInTheDocument()
    expect(screen.getAllByText(/1L/)[0]).toBeInTheDocument()
    expect(screen.getAllByText('2')[0]).toBeInTheDocument()
  })
})
