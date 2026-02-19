import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { describe, it, expect, beforeEach } from 'vitest'
import { useAuthStore } from '@/store/auth'
import { SessionExpiredDialog } from './session-expired-dialog'

const queryClient = new QueryClient()

describe('SessionExpiredDialog', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      token: 'expired-token',
      sessionExpired: true,
    })
  })

  it('renders login form and logout button when session is expired', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <SessionExpiredDialog />
        </BrowserRouter>
      </QueryClientProvider>
    )

    expect(screen.getByText(/session has expired/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    // The login button text depends on loading state but starts as 'Login'
    expect(screen.getByRole('button', { name: /^login$/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /log out/i })).toBeInTheDocument()
  })

  it('closes dialog and updates auth state on successful login', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <SessionExpiredDialog />
        </BrowserRouter>
      </QueryClientProvider>
    )

    const emailInput = screen.getByLabelText(/email/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const loginButton = screen.getByRole('button', { name: /^login$/i })

    await userEvent.type(emailInput, 'test@example.com')
    await userEvent.type(passwordInput, 'password')
    await userEvent.click(loginButton)

    await waitFor(() => {
      const dialog = screen.getByText(/session has expired/i).closest('dialog')
      expect(dialog).not.toHaveAttribute('open')
    })

    const state = useAuthStore.getState()
    expect(state.sessionExpired).toBe(false)
    expect(state.token).toBe('fake-jwt-token')
    expect(state.user?.email).toBe('test@example.com')
  })

  it('logs out and closes dialog when logout button is clicked', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <SessionExpiredDialog />
        </BrowserRouter>
      </QueryClientProvider>
    )

    const logoutButton = screen.getByRole('button', { name: /log out/i })
    await userEvent.click(logoutButton)

    await waitFor(() => {
      const dialog = screen.getByText(/session has expired/i).closest('dialog')
      expect(dialog).not.toHaveAttribute('open')
    })

    const state = useAuthStore.getState()
    expect(state.sessionExpired).toBe(false)
    expect(state.token).toBeNull()
    expect(state.user).toBeNull()
  })
})
