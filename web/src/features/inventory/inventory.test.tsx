import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { InventorySelectionRoute } from './routes/inventory-selection'
import { InventoryGuard } from './components/inventory-guard'
import { MainLayout } from '@/components/layout/main-layout'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

describe('Inventory Flow', () => {
  beforeEach(() => {
    // Authenticated user
    useAuthStore.setState({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-token',
    })
    // No inventory selected initially
    useInventoryStore.setState({ activeInventoryId: null })
  })

  it('redirects to selection page if no inventory selected', async () => {
    const router = createMemoryRouter(
      [
        {
          path: '/select-inventory',
          element: <InventorySelectionRoute />,
        },
        {
          element: <InventoryGuard />,
          children: [
            {
              path: '/',
              element: <div>Dashboard</div>,
            },
          ],
        },
      ],
      {
        initialEntries: ['/'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Select Household')).toBeInTheDocument()
    })
  })

  it('allows selecting an existing inventory', async () => {
    // handlers.ts returns one inventory "My Household" with ID 'inv-1'
    const router = createMemoryRouter(
      [
        {
          path: '/select-inventory',
          element: <InventorySelectionRoute />,
        },
        {
          element: <InventoryGuard />,
          children: [
            {
              element: <MainLayout />, // To test layout display
              children: [
                {
                  path: '/',
                  element: <div>Dashboard</div>,
                },
              ],
            },
          ],
        },
      ],
      {
        initialEntries: ['/select-inventory'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    const inventoryButton = await screen.findByText('My Household')
    await userEvent.click(inventoryButton)

    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
    })

    // Check if store is updated
    expect(useInventoryStore.getState().activeInventoryId).toBe('inv-1')
  })

  it('allows creating a new inventory', async () => {
    const router = createMemoryRouter(
      [
        {
          path: '/select-inventory',
          element: <InventorySelectionRoute />,
        },
        {
          element: <InventoryGuard />,
          children: [
            {
              path: '/',
              element: <div>Dashboard</div>,
            },
          ],
        },
      ],
      {
        initialEntries: ['/select-inventory'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    const nameInput = screen.getByLabelText(/Household Name/i)
    const createButton = screen.getByRole('button', {
      name: /Create Household/i,
    })

    await userEvent.type(nameInput, 'New Place')
    await userEvent.click(createButton)

    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
    })

    expect(useInventoryStore.getState().activeInventoryId).toContain('inv-')
  })
})
