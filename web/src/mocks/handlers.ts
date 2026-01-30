import { http, HttpResponse } from 'msw'

const BASE_URL = 'http://localhost:8080'

export const handlers = [
  http.post(`${BASE_URL}/login`, async ({ request }) => {
    const { email, password } = (await request.json()) as any

    if (email === 'test@example.com' && password === 'password') {
      return HttpResponse.json({
        user: { id: '1', name: 'Test User', email: 'test@example.com' },
        token: 'fake-jwt-token',
      })
    }

    return HttpResponse.json({ message: 'Invalid credentials' }, { status: 401 })
  }),

  http.post(`${BASE_URL}/signup`, async () => {
    return HttpResponse.json({
      user: { id: '1', name: 'Test User', email: 'test@example.com' },
      token: 'fake-jwt-token',
    }, { status: 201 })
  }),
]
