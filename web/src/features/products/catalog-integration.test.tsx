import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach, afterEach, afterAll, beforeAll } from 'vitest'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { ProductsPage } from './routes/products-page'
import { ProductDetailsPage } from './routes/product-details-page'
import { BASE_URL } from '@/lib/api'

// Define handlers
const handlers = [
  http.get(`${BASE_URL}/inventories/:id/canonical-products`, ({ request }) => {
    const url = new URL(request.url)
    const search = url.searchParams.get('search')
    if (search === 'Bread') {
        return HttpResponse.json([
            {
              id: 'prod-bread',
              inventory_id: 'inv-1',
              name: 'Bread',
              description: 'Fresh bread',
              created_at: new Date().toISOString(),
            },
          ])
    }
    return HttpResponse.json([
      {
        id: 'prod-1',
        inventory_id: 'inv-1',
        name: 'Milk',
        description: 'Cow milk',
        created_at: new Date().toISOString(),
      },
    ])
  }),

  http.get(`${BASE_URL}/canonical-products/:id`, () => {
    return HttpResponse.json({
      id: 'prod-1',
      inventory_id: 'inv-1',
      name: 'Milk',
      description: 'Cow milk',
      created_at: new Date().toISOString(),
    })
  }),

  http.put(`${BASE_URL}/canonical-products/:id`, async ({ request }) => {
    const body = (await request.json()) as any
    return HttpResponse.json({
      id: 'prod-1',
      inventory_id: 'inv-1',
      ...body,
      created_at: new Date().toISOString(),
    })
  }),

  http.get(`${BASE_URL}/inventories/:id/products`, () => {
    return HttpResponse.json([
        {
          id: 'brand-1',
          inventory_id: 'inv-1',
          canonical_product_id: 'prod-1',
          brand: 'Tesco',
          name: 'Tesco Whole Milk',
          created_at: new Date().toISOString(),
        },
      ])
  }),

  http.put(`${BASE_URL}/products/:id`, async ({ request }) => {
      const body = (await request.json()) as any
      return HttpResponse.json({
        id: 'brand-1',
        inventory_id: 'inv-1',
        ...body,
        created_at: new Date().toISOString(),
      })
  }),

  http.get(`${BASE_URL}/products/:id/variants`, () => {
    return HttpResponse.json([
        {
          id: 'var-1',
          product_id: 'brand-1',
          variant_name: '1L',
          size: 1,
          unit: 'L',
        },
      ])
  }),

  http.put(`${BASE_URL}/product-variants/:id`, async ({ request }) => {
      const body = (await request.json()) as any
      return HttpResponse.json({
        id: 'var-1',
        product_id: 'brand-1',
        ...body,
      })
  }),
]

const server = setupServer(...handlers)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('Product Catalog Integration', () => {
    beforeEach(() => {
        useAuthStore.setState({
            user: { id: '1', name: 'Test User', email: 'test@example.com' },
            token: 'fake-token',
        })
        useInventoryStore.setState({ activeInventoryId: 'inv-1' })
    })

    it('searches for canonical products', async () => {
        const user = userEvent.setup()
        const queryClient = new QueryClient({
            defaultOptions: { queries: { retry: false } },
        })

        render(
            <QueryClientProvider client={queryClient}>
                <MemoryRouter initialEntries={['/products']}>
                    <Routes>
                        <Route path="/products" element={<ProductsPage />} />
                    </Routes>
                </MemoryRouter>
            </QueryClientProvider>
        )

        await waitFor(() => {
            expect(screen.getByText('Milk')).toBeInTheDocument()
        })

        const searchInput = screen.getByPlaceholderText('Search products...')
        await user.type(searchInput, 'Bread')

        await waitFor(() => {
            expect(screen.getByText('Bread')).toBeInTheDocument()
        }, { timeout: 3000 })

        expect(screen.queryByText('Milk')).not.toBeInTheDocument()
    })

    it('navigates to details and edits canonical product', async () => {
        const user = userEvent.setup()
        const queryClient = new QueryClient({
             defaultOptions: { queries: { retry: false } }
        })

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

        await user.click(screen.getByRole('button', { name: 'Edit' }))

        await screen.findByText('Edit Product', { selector: 'h3' })

        const nameInput = screen.getByDisplayValue('Milk')
        await user.clear(nameInput)
        await user.type(nameInput, 'Soy Milk')

        const saveButton = screen.getByRole('button', { name: 'Save Changes' })
        await user.click(saveButton)

        await waitFor(() => {
            expect(screen.queryByText('Edit Product', { selector: 'h3' })).not.toBeInTheDocument()
        })
    })
})
