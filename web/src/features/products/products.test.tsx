import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { ProductsPage } from './routes/products-page'
import { ProductDetailsPage } from './routes/product-details-page'
import { MemoryRouter, Route, Routes } from 'react-router-dom'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

describe('Products Features', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
    useInventoryStore.setState({ activeInventoryId: 'inv-1' })
    queryClient.clear()
  })

  it('renders product page and lists products', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
            <ProductsPage />
        </MemoryRouter>
      </QueryClientProvider>
    )

    expect(screen.getByText('Product Catalog')).toBeInTheDocument()
    expect(screen.getByText('Add New Product')).toBeInTheDocument()

    await waitFor(() => {
        expect(screen.getByText('Milk')).toBeInTheDocument()
    })
  })

  it('filters products based on search', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
            <ProductsPage />
        </MemoryRouter>
      </QueryClientProvider>
    )

    await waitFor(() => {
        expect(screen.getByText('Milk')).toBeInTheDocument()
    })

    const searchInput = screen.getByPlaceholderText('Search products...')
    fireEvent.change(searchInput, { target: { value: 'Bread' } }) // Milk should disappear, but we don't have Bread in mocks

    await waitFor(() => {
        expect(screen.queryByText('Milk')).not.toBeInTheDocument()
        expect(screen.getByText('No products match your search.')).toBeInTheDocument()
    })

    fireEvent.change(searchInput, { target: { value: 'Milk' } })
    await waitFor(() => {
        expect(screen.getByText('Milk')).toBeInTheDocument()
    })
  })

  it('renders product details page', async () => {
    render(
        <QueryClientProvider client={queryClient}>
          <MemoryRouter initialEntries={['/products/prod-1']}>
            <Routes>
                <Route path="/products/:id" element={<ProductDetailsPage />} />
            </Routes>
          </MemoryRouter>
        </QueryClientProvider>
      )

      await waitFor(() => {
          expect(screen.getByRole('heading', { name: 'Milk' })).toBeInTheDocument()
          expect(screen.getByText('Cow milk')).toBeInTheDocument()
      })

      await waitFor(() => {
        expect(screen.getByText('Tesco Whole Milk')).toBeInTheDocument()
      })
  })

  it('opens edit dialog for canonical product', async () => {
    render(
        <QueryClientProvider client={queryClient}>
          <MemoryRouter initialEntries={['/products/prod-1']}>
            <Routes>
                <Route path="/products/:id" element={<ProductDetailsPage />} />
            </Routes>
          </MemoryRouter>
        </QueryClientProvider>
      )

      await waitFor(() => {
          expect(screen.getByRole('heading', { name: 'Milk' })).toBeInTheDocument()
      })

      fireEvent.click(screen.getByText('Edit'))

      expect(screen.getByRole('dialog')).toBeInTheDocument()
      const dialog = screen.getByRole('dialog')
      expect(within(dialog).getByLabelText('Product Name')).toHaveValue('Milk')
  })
})
