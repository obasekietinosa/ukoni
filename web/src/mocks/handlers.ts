import { http, HttpResponse } from 'msw'
import { BASE_URL } from '@/lib/api'

const inventories = [
  {
    id: 'inv-1',
    name: 'My Household',
    owner_user_id: '1',
    created_at: new Date().toISOString(),
  },
]

export const handlers = [
  http.post(`${BASE_URL}/login`, async ({ request }) => {
    const { email, password } = (await request.json()) as {
      email: string
      password: string
    }

    if (email === 'test@example.com' && password === 'password') {
      return HttpResponse.json({
        user: { id: '1', name: 'Test User', email: 'test@example.com' },
        token: 'fake-jwt-token',
      })
    }

    return HttpResponse.json(
      { message: 'Invalid credentials' },
      { status: 401 }
    )
  }),

  http.post(`${BASE_URL}/signup`, async () => {
    return HttpResponse.json(
      {
        user: { id: '1', name: 'Test User', email: 'test@example.com' },
        token: 'fake-jwt-token',
      },
      { status: 201 }
    )
  }),

  http.get(`${BASE_URL}/inventories`, () => {
    return HttpResponse.json(inventories)
  }),

  http.post(`${BASE_URL}/inventories`, async ({ request }) => {
    const { name } = (await request.json()) as { name: string }
    const newInventory = {
      id: `inv-${Date.now()}`,
      name,
      owner_user_id: '1',
      created_at: new Date().toISOString(),
    }
    inventories.push(newInventory)
    return HttpResponse.json(newInventory, { status: 201 })
  }),

  http.get(`${BASE_URL}/inventories/:id`, ({ params }) => {
    const inventory = inventories.find((i) => i.id === params.id)
    if (!inventory) {
      return new HttpResponse(null, { status: 404 })
    }
    return HttpResponse.json(inventory)
  }),
]
