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

  http.get(`${BASE_URL}/inventories/:id/inventory-products`, () => {
    return HttpResponse.json([
      {
        id: 'inv-prod-1',
        inventory_id: 'inv-1',
        product_variant_id: 'var-1',
        canonical_product_name: 'Milk',
        brand_name: 'Tesco',
        variant_name: '1L',
        quantity: 2,
        unit: 'L',
        last_updated: new Date().toISOString(),
      },
    ])
  }),

  http.post(`${BASE_URL}/inventories/:id/transactions`, async ({ request }) => {
    const body = (await request.json()) as any
    return HttpResponse.json(
      {
        id: `tx-${Date.now()}`,
        ...body,
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

  http.get(`${BASE_URL}/inventories/:id/members`, ({ params }) => {
    return HttpResponse.json([
      {
        id: 'mem-1',
        inventory_id: params.id as string,
        user_id: '1',
        role: 'admin',
        invited_at: new Date().toISOString(),
      },
      {
        id: 'mem-2',
        inventory_id: params.id as string,
        user_id: '2',
        role: 'viewer',
        invited_at: new Date().toISOString(),
      },
    ])
  }),

  http.get(`${BASE_URL}/inventories/:id/canonical-products`, ({ params }) => {
    return HttpResponse.json([
      {
        id: 'prod-1',
        inventory_id: params.id as string,
        name: 'Milk',
        description: 'Cow milk',
        created_at: new Date().toISOString(),
      },
    ])
  }),

  http.post(
    `${BASE_URL}/inventories/:id/canonical-products`,
    async ({ request }) => {
      const body = (await request.json()) as {
        name: string
        description?: string
        category_id?: string
      }
      return HttpResponse.json(
        {
          id: `prod-${Date.now()}`,
          ...body,
          created_at: new Date().toISOString(),
        },
        { status: 201 }
      )
    }
  ),

  http.get(`${BASE_URL}/canonical-products/:id`, ({ params }) => {
    return HttpResponse.json({
      id: params.id as string,
      inventory_id: 'inv-1',
      name: 'Milk',
      description: 'Cow milk',
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

  http.post(`${BASE_URL}/inventories/:id/products`, async ({ request }) => {
    const body = (await request.json()) as {
      name: string
      brand?: string
      description?: string
      category_id?: string
      canonical_product_id?: string
    }
    return HttpResponse.json(
      {
        id: `brand-${Date.now()}`,
        ...body,
        created_at: new Date().toISOString(),
      },
      { status: 201 }
    )
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

  http.post(`${BASE_URL}/products/:id/variants`, async ({ request }) => {
    const body = (await request.json()) as {
      variant_name: string
      sku?: string
      unit?: string
      size?: number
    }
    return HttpResponse.json(
      {
        id: `var-${Date.now()}`,
        ...body,
      },
      { status: 201 }
    )
  }),
]
