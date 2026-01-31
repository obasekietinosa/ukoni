import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { InventoryPage } from '../routes/inventory-page'
import { MemoryRouter } from 'react-router-dom'
import { http, HttpResponse } from 'msw'
import { server } from '@/test/setup'
import { BASE_URL } from '@/lib/api'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

describe('Inventory Page', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
    useInventoryStore.setState({ activeInventoryId: 'inv-1' })
  })

  it('renders inventory list', async () => {
    // Override handler to ensure specific data
    server.use(
        http.get(`${BASE_URL}/inventories/:id/inventory-products`, () => {
            return HttpResponse.json([
              {
                id: 'inv-prod-1',
                inventory_id: 'inv-1',
                product_variant_id: 'var-1',
                canonical_product_name: 'Milk',
                brand_name: 'Tesco',
                variant_name: '1L',
                quantity: 2,
                unit: 'L',
                last_updated: new Date().toISOString(),
              },
            ])
        }),
    )

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
            <InventoryPage />
        </MemoryRouter>
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Milk')).toBeInTheDocument()
    })
    expect(screen.getByText('Tesco')).toBeInTheDocument()
    expect(screen.getByText('1L')).toBeInTheDocument()
    expect(screen.getByText('2 L')).toBeInTheDocument()
  })

  it('renders empty state', async () => {
    server.use(
        http.get(`${BASE_URL}/inventories/:id/inventory-products`, () => {
            return HttpResponse.json([])
        }),
    )

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
            <InventoryPage />
        </MemoryRouter>
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText(/Your inventory is empty/)).toBeInTheDocument()
    })
  })
})
