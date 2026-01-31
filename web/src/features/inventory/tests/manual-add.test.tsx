import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { ProductDetailsPage } from '@/features/products/routes/product-details-page'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
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

describe('Manual Add to Inventory', () => {
  beforeEach(() => {
    queryClient.clear()
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
    useInventoryStore.setState({ activeInventoryId: 'inv-1' })
  })

  it('adds an item to inventory from product details', async () => {
    // Mock Product, Brand, Variant
    server.use(
        http.get(`${BASE_URL}/canonical-products/prod-1`, () => {
            return HttpResponse.json({
                id: 'prod-1',
                inventory_id: 'inv-1',
                name: 'Milk',
                description: 'Cow milk',
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
                },
            ])
        }),
        http.get(`${BASE_URL}/products/brand-1/variants`, () => {
            return HttpResponse.json([
                {
                    id: 'var-1',
                    product_id: 'brand-1',
                    variant_name: '1L',
                    size: 1,
                    unit: 'L',
                },
            ])
        })
    )

    let transactionPayload: any = null
    server.use(
        http.post(`${BASE_URL}/inventories/:id/transactions`, async ({ request }) => {
            transactionPayload = await request.json()
            return HttpResponse.json({ id: 'tx-1' }, { status: 201 })
        })
    )

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={['/products/prod-1']}>
            <Routes>
                <Route path="/products/:id" element={<ProductDetailsPage />} />
            </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    )

    // Wait for variant to load
    await waitFor(() => {
        expect(screen.getByText('1L')).toBeInTheDocument()
    })

    // Click Add button
    const addButton = screen.getByRole('button', { name: /\+ Add/i })
    await userEvent.click(addButton)

    // Verify Dialog opens
    expect(screen.getByText('Add 1L to Inventory')).toBeInTheDocument()

    // Enter quantity
    const quantityInput = screen.getByLabelText(/Quantity/i)
    await userEvent.clear(quantityInput)
    await userEvent.type(quantityInput, '5')

    // Submit
    const submitButton = screen.getByRole('button', { name: 'Add to Inventory' })
    await userEvent.click(submitButton)

    // Verify transaction created
    await waitFor(() => {
        expect(transactionPayload).not.toBeNull()
    })
    expect(transactionPayload.items).toHaveLength(1)
    expect(transactionPayload.items[0].product_variant_id).toBe('var-1')
    expect(transactionPayload.items[0].quantity).toBe(5)
  })
})
