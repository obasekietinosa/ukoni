import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { ShoppingListsPage } from './routes/shopping-lists-page'
import { ShoppingListDetailsPage } from './routes/shopping-list-details-page'
import { InventoryGuard } from '@/features/inventory/components/inventory-guard'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

describe('Shopping Lists', () => {
  beforeEach(() => {
    // Authenticated user
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
    // Selected inventory
    useInventoryStore.setState({ activeInventoryId: 'inv-1' })
  })

  it('renders list of shopping lists', async () => {
    const router = createMemoryRouter(
      [
        {
          element: <InventoryGuard />,
          children: [
            {
              path: '/shopping-lists',
              element: <ShoppingListsPage />,
            },
          ],
        },
      ],
      {
        initialEntries: ['/shopping-lists'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Weekly Shop')).toBeInTheDocument()
    })
  })

  it('allows creating a new shopping list', async () => {
    const router = createMemoryRouter(
      [
         {
          element: <InventoryGuard />,
          children: [
            {
              path: '/shopping-lists',
              element: <ShoppingListsPage />,
            },
          ],
        },
      ],
      {
        initialEntries: ['/shopping-lists'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    const createButton = await screen.findByRole('button', { name: /New List/i })
    await userEvent.click(createButton)

    const dialog = screen.getByRole('dialog')
    const nameInput = within(dialog).getByLabelText(/Name/i)
    await userEvent.type(nameInput, 'Party List')

    const submitButton = within(dialog).getByRole('button', { name: /Create/i })
    await userEvent.click(submitButton)
  })

  it('renders shopping list details and items', async () => {
    const router = createMemoryRouter(
      [
         {
          element: <InventoryGuard />,
          children: [
            {
              path: '/shopping-lists/:id',
              element: <ShoppingListDetailsPage />,
            },
          ],
        },
      ],
      {
        initialEntries: ['/shopping-lists/list-1'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Weekly Shop')).toBeInTheDocument()
    })

    expect(screen.getByText('Milk')).toBeInTheDocument()
    expect(screen.getByText('"Get fresh one"')).toBeInTheDocument()
  })

  it('opens add item dialog and shows preferred outlet selection', async () => {
      const router = createMemoryRouter(
      [
         {
          element: <InventoryGuard />,
          children: [
            {
              path: '/shopping-lists/:id',
              element: <ShoppingListDetailsPage />,
            },
          ],
        },
      ],
      {
        initialEntries: ['/shopping-lists/list-1'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    const addButton = await screen.findByRole('button', { name: /Add Item/i })
    await userEvent.click(addButton)

    const dialog = screen.getByRole('dialog')
    const searchInput = within(dialog).getByPlaceholderText(/Search products/i)

    // Simulate searching
    await userEvent.type(searchInput, 'Milk')

    // Click on the result
    // Note: The mock data in handlers.ts for Canonical Products is 'Milk'
    const milkButton = await screen.findByText('Milk', { selector: 'div' })
    await userEvent.click(milkButton)

    // Now in Details step
    const outletSelect = screen.getByLabelText(/Preferred Outlet/i)
    expect(outletSelect).toBeInTheDocument()

    // Check if outlets are populated (Tesco Extra from MSW)
    await waitFor(() => {
       expect(screen.getByText(/Tesco - Tesco Extra/i)).toBeInTheDocument()
    })
  })
})
