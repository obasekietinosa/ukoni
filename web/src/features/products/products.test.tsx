import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { ProductsPage } from './routes/products-page'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

describe('Products Page', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
    useInventoryStore.setState({ activeInventoryId: 'inv-1' })
  })

  it('renders product page', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <ProductsPage />
      </QueryClientProvider>
    )

    expect(screen.getByText('Product Catalog')).toBeInTheDocument()
    expect(screen.getByText('Add New Product')).toBeInTheDocument()
  })
})
