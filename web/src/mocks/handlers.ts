import { http, HttpResponse } from 'msw'
import { BASE_URL } from '@/lib/api'

// Initialize with some data
let inventories = [
  {
    id: 'inv-1',
    name: 'My Household',
    owner_user_id: '1',
    created_at: new Date().toISOString(),
  },
]

let canonicalProducts = [
  {
    id: 'prod-1',
    inventory_id: 'inv-1',
    name: 'Milk',
    description: 'Cow milk',
    created_at: new Date().toISOString(),
  },
]

let products = [
  {
    id: 'brand-1',
    inventory_id: 'inv-1',
    canonical_product_id: 'prod-1',
    brand: 'Tesco',
    name: 'Tesco Whole Milk',
    created_at: new Date().toISOString(),
  },
]

let productVariants = [
  {
    id: 'var-1',
    product_id: 'brand-1',
    variant_name: '1L',
    size: 1,
    unit: 'L',
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

  // Canonical Products
  http.get(`${BASE_URL}/inventories/:id/canonical-products`, ({ params }) => {
    return HttpResponse.json(
      canonicalProducts.filter((p) => p.inventory_id === params.id)
    )
  }),

  http.post(
    `${BASE_URL}/inventories/:id/canonical-products`,
    async ({ request, params }) => {
      const body = (await request.json()) as {
        name: string
        description?: string
        category_id?: string
      }
      const newProduct = {
        id: `prod-${Date.now()}`,
        inventory_id: params.id as string,
        description: body.description || '',
        ...body,
        created_at: new Date().toISOString(),
      }
      canonicalProducts.push(newProduct)
      return HttpResponse.json(newProduct, { status: 201 })
    }
  ),

  http.get(`${BASE_URL}/canonical-products/:id`, ({ params }) => {
    const product = canonicalProducts.find((p) => p.id === params.id)
    if (!product) return new HttpResponse(null, { status: 404 })
    return HttpResponse.json(product)
  }),

  http.put(
    `${BASE_URL}/canonical-products/:id`,
    async ({ request, params }) => {
      const body = (await request.json()) as Partial<typeof canonicalProducts[0]>
      const index = canonicalProducts.findIndex((p) => p.id === params.id)
      if (index === -1) return new HttpResponse(null, { status: 404 })

      canonicalProducts[index] = { ...canonicalProducts[index], ...body }
      return HttpResponse.json(canonicalProducts[index])
    }
  ),

  http.delete(`${BASE_URL}/canonical-products/:id`, ({ params }) => {
    const index = canonicalProducts.findIndex((p) => p.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    canonicalProducts.splice(index, 1)
    return new HttpResponse(null, { status: 204 })
  }),


  // Products (Brands)
  http.get(`${BASE_URL}/inventories/:id/products`, ({ params }) => {
    // Basic filter simulation
    return HttpResponse.json(
      products.filter(p => p.inventory_id === params.id)
    )
  }),

  http.post(`${BASE_URL}/inventories/:id/products`, async ({ request, params }) => {
    const body = (await request.json()) as any
    const newProduct = {
        id: `brand-${Date.now()}`,
        inventory_id: params.id as string,
        ...body,
        created_at: new Date().toISOString(),
    }
    products.push(newProduct)
    return HttpResponse.json(newProduct, { status: 201 })
  }),

  http.put(`${BASE_URL}/products/:id`, async ({ request, params }) => {
     const body = (await request.json()) as any
     const index = products.findIndex(p => p.id === params.id)
     if (index === -1) return new HttpResponse(null, { status: 404 })
     products[index] = { ...products[index], ...body }
     return HttpResponse.json(products[index])
  }),

  http.delete(`${BASE_URL}/products/:id`, ({ params }) => {
    const index = products.findIndex(p => p.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    products.splice(index, 1)
    return new HttpResponse(null, { status: 204 })
  }),


  // Variants
  http.get(`${BASE_URL}/products/:id/variants`, ({ params }) => {
    return HttpResponse.json(productVariants.filter(v => v.product_id === params.id))
  }),

  http.post(`${BASE_URL}/products/:id/variants`, async ({ request, params }) => {
    const body = (await request.json()) as any
    const newVariant = {
        id: `var-${Date.now()}`,
        product_id: params.id as string,
        ...body,
    }
    productVariants.push(newVariant)
    return HttpResponse.json(newVariant, { status: 201 })
  }),

  http.get(`${BASE_URL}/inventories/:id/inventory-products`, ({ params }) => {
    return HttpResponse.json([
      {
        id: 'ip-1',
        inventory_id: params.id as string,
        product_variant_id: 'var-1',
        quantity: 2,
        unit: 'L',
        created_at: new Date().toISOString(),
        last_updated: new Date().toISOString(),
        product_name: 'Tesco Whole Milk',
        brand: 'Tesco',
        variant_name: '1L',
        size: 1,
        product_unit: 'L',
      },
    ])
  }),
]
