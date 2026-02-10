import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface User {
  id: string
  name: string
  email: string
}

interface AuthState {
  user: User | null
  token: string | null
  sessionExpired: boolean
  setAuth: (user: User, token: string) => void
  clearAuth: () => void
  setSessionExpired: (expired: boolean) => void
  isAuthenticated: () => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      sessionExpired: false,
      setAuth: (user, token) => set({ user, token, sessionExpired: false }),
      clearAuth: () => set({ user: null, token: null, sessionExpired: false }),
      setSessionExpired: (expired) => set({ sessionExpired: expired }),
      isAuthenticated: () => !!get().token,
    }),
    {
      name: 'auth-storage',
    }
  )
)
