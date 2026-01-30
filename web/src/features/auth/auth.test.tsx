import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { LoginRoute } from './routes/login'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'

const queryClient = new QueryClient()

describe('Login Flow', () => {
  beforeEach(() => {
    useAuthStore.getState().clearAuth()
  })

  it('allows user to login successfully', async () => {
    const router = createMemoryRouter(
      [
        {
          path: '/login',
          element: <LoginRoute />,
        },
        {
          path: '/',
          element: <div>Home Page</div>,
        },
      ],
      {
        initialEntries: ['/login'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    const emailInput = screen.getByLabelText(/email/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /login/i })

    await userEvent.type(emailInput, 'test@example.com')
    await userEvent.type(passwordInput, 'password')
    await userEvent.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Home Page')).toBeInTheDocument()
    })

    expect(useAuthStore.getState().user?.email).toBe('test@example.com')
  })

  it('shows error on invalid credentials', async () => {
    const router = createMemoryRouter(
      [
        {
          path: '/login',
          element: <LoginRoute />,
        },
      ],
      {
        initialEntries: ['/login'],
      }
    )

    render(
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    )

    const emailInput = screen.getByLabelText(/email/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /login/i })

    await userEvent.type(emailInput, 'wrong@example.com')
    await userEvent.type(passwordInput, 'wrong')
    await userEvent.click(submitButton)

    await waitFor(() => {
      // "Invalid credentials" is what the mock returns and ApiError parsing extracts
      expect(screen.getByText(/Invalid credentials/i)).toBeInTheDocument()
    })
  })
})
