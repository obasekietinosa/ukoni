import { http, HttpResponse } from 'msw'
import { BASE_URL } from '@/lib/api'
import type {
  Transaction,
  CreateTransactionRequest,
} from '@/features/transactions/types'

const transactions: Transaction[] = []

// Initialize with some data
const inventories = [
  {
    id: 'inv-1',
    name: 'My Household',
    owner_user_id: '1',
    created_at: new Date().toISOString(),
  },
]

const canonicalProducts = [
  {
    id: 'prod-1',
    inventory_id: 'inv-1',
    name: 'Milk',
    description: 'Cow milk',
    created_at: new Date().toISOString(),
  },
]

const products = [
  {
    id: 'brand-1',
    inventory_id: 'inv-1',
    canonical_product_id: 'prod-1',
    brand: 'Tesco',
    name: 'Tesco Whole Milk',
    created_at: new Date().toISOString(),
  },
]

const productVariants = [
  {
    id: 'var-1',
    product_id: 'brand-1',
    variant_name: '1L',
    size: 1,
    unit: 'L',
  },
]

const sellers = [
  {
    id: 'seller-1',
    name: 'Tesco',
    type: 'chain',
    created_at: new Date().toISOString(),
  },
]

const outlets = [
  {
    id: 'outlet-1',
    seller_id: 'seller-1',
    name: 'Tesco Extra',
    channel: 'physical',
    address: '123 High St',
    created_at: new Date().toISOString(),
  },
]

const inventoryProducts = [
  {
    id: 'ip-1',
    inventory_id: 'inv-1',
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
      const body = (await request.json()) as Partial<
        (typeof canonicalProducts)[0]
      >
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
      products.filter((p) => p.inventory_id === params.id)
    )
  }),

  http.post(
    `${BASE_URL}/inventories/:id/products`,
    async ({ request, params }) => {
      const body = (await request.json()) as Omit<
        (typeof products)[number],
        'id' | 'inventory_id' | 'created_at'
      >
      const newProduct = {
        id: `brand-${Date.now()}`,
        inventory_id: params.id as string,
        ...body,
        created_at: new Date().toISOString(),
      }
      products.push(newProduct)
      return HttpResponse.json(newProduct, { status: 201 })
    }
  ),

  http.put(`${BASE_URL}/products/:id`, async ({ request, params }) => {
    const body = (await request.json()) as Partial<(typeof products)[number]>
    const index = products.findIndex((p) => p.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    products[index] = { ...products[index], ...body }
    return HttpResponse.json(products[index])
  }),

  http.delete(`${BASE_URL}/products/:id`, ({ params }) => {
    const index = products.findIndex((p) => p.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    products.splice(index, 1)
    return new HttpResponse(null, { status: 204 })
  }),

  // Variants
  http.get(`${BASE_URL}/products/:id/variants`, ({ params }) => {
    return HttpResponse.json(
      productVariants.filter((v) => v.product_id === params.id)
    )
  }),

  http.post(
    `${BASE_URL}/products/:id/variants`,
    async ({ request, params }) => {
      const body = (await request.json()) as Omit<
        (typeof productVariants)[number],
        'id' | 'product_id'
      >
      const newVariant = {
        id: `var-${Date.now()}`,
        product_id: params.id as string,
        ...body,
      }
      productVariants.push(newVariant)
      return HttpResponse.json(newVariant, { status: 201 })
    }
  ),

  http.get(`${BASE_URL}/inventories/:id/inventory-products`, ({ params }) => {
    return HttpResponse.json(
      inventoryProducts.filter((ip) => ip.inventory_id === params.id)
    )
  }),

  // Sellers
  http.get(`${BASE_URL}/sellers`, () => {
    return HttpResponse.json(sellers)
  }),

  http.post(`${BASE_URL}/sellers`, async ({ request }) => {
    const body = (await request.json()) as Omit<
      (typeof sellers)[number],
      'id' | 'created_at'
    >
    const newSeller = {
      id: `seller-${Date.now()}`,
      ...body,
      created_at: new Date().toISOString(),
    }
    sellers.push(newSeller)
    return HttpResponse.json(newSeller, { status: 201 })
  }),

  http.get(`${BASE_URL}/sellers/:id`, ({ params }) => {
    const seller = sellers.find((s) => s.id === params.id)
    if (!seller) return new HttpResponse(null, { status: 404 })
    return HttpResponse.json(seller)
  }),

  http.put(`${BASE_URL}/sellers/:id`, async ({ request, params }) => {
    const body = (await request.json()) as Partial<(typeof sellers)[number]>
    const index = sellers.findIndex((s) => s.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    sellers[index] = { ...sellers[index], ...body }
    return HttpResponse.json(sellers[index])
  }),

  http.delete(`${BASE_URL}/sellers/:id`, ({ params }) => {
    const index = sellers.findIndex((s) => s.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    sellers.splice(index, 1)
    return new HttpResponse(null, { status: 204 })
  }),

  // Outlets
  http.get(`${BASE_URL}/sellers/:id/outlets`, ({ params }) => {
    return HttpResponse.json(outlets.filter((o) => o.seller_id === params.id))
  }),

  http.post(`${BASE_URL}/sellers/:id/outlets`, async ({ request, params }) => {
    const body = (await request.json()) as Omit<
      (typeof outlets)[number],
      'id' | 'seller_id' | 'created_at'
    >
    const newOutlet = {
      id: `outlet-${Date.now()}`,
      seller_id: params.id as string,
      ...body,
      created_at: new Date().toISOString(),
    }
    outlets.push(newOutlet)
    return HttpResponse.json(newOutlet, { status: 201 })
  }),

  http.get(`${BASE_URL}/outlets/:id`, ({ params }) => {
    const outlet = outlets.find((o) => o.id === params.id)
    if (!outlet) return new HttpResponse(null, { status: 404 })
    return HttpResponse.json(outlet)
  }),

  http.put(`${BASE_URL}/outlets/:id`, async ({ request, params }) => {
    const body = (await request.json()) as Partial<(typeof outlets)[number]>
    const index = outlets.findIndex((o) => o.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    outlets[index] = { ...outlets[index], ...body }
    return HttpResponse.json(outlets[index])
  }),

  http.delete(`${BASE_URL}/outlets/:id`, ({ params }) => {
    const index = outlets.findIndex((o) => o.id === params.id)
    if (index === -1) return new HttpResponse(null, { status: 404 })
    outlets.splice(index, 1)
    return new HttpResponse(null, { status: 204 })
  }),

  // Shopping Lists
  http.get(`${BASE_URL}/inventories/:id/shopping-lists`, ({ params }) => {
    return HttpResponse.json([
      {
        id: 'list-1',
        inventory_id: params.id as string,
        name: 'Weekly Shop',
        created_at: new Date().toISOString(),
        last_updated_at: new Date().toISOString(),
      },
    ])
  }),

  http.post(
    `${BASE_URL}/inventories/:id/shopping-lists`,
    async ({ request, params }) => {
      const { name } = (await request.json()) as { name: string }
      return HttpResponse.json(
        {
          id: `list-${Date.now()}`,
          inventory_id: params.id as string,
          name,
          created_at: new Date().toISOString(),
          last_updated_at: new Date().toISOString(),
        },
        { status: 201 }
      )
    }
  ),

  http.get(`${BASE_URL}/shopping-lists/:id`, ({ params }) => {
    if (params.id === 'list-1') {
      return HttpResponse.json({
        id: 'list-1',
        inventory_id: 'inv-1',
        name: 'Weekly Shop',
        created_at: new Date().toISOString(),
        last_updated_at: new Date().toISOString(),
      })
    }
    return new HttpResponse(null, { status: 404 })
  }),

  http.get(`${BASE_URL}/shopping-lists/:id/items`, ({ params }) => {
    if (params.id === 'list-1') {
      return HttpResponse.json([
        {
          id: 'item-1',
          shopping_list_id: 'list-1',
          target_type: 'canonical_product',
          target_id: 'prod-1',
          quantity: 2,
          unit: 'L',
          notes: 'Get fresh one',
          created_at: new Date().toISOString(),
          canonical_product: {
            id: 'prod-1',
            name: 'Milk',
          },
        },
      ])
    }
    return HttpResponse.json([])
  }),

  http.post(`${BASE_URL}/shopping-lists/:id/items`, async ({ request }) => {
    const body = await request.json()
    return HttpResponse.json(
      {
        id: `item-${Date.now()}`,
        ...(body as object),
        created_at: new Date().toISOString(),
      },
      { status: 201 }
    )
  }),

  // Transactions
  http.get(`${BASE_URL}/inventories/:id/transactions`, ({ params }) => {
    return HttpResponse.json(
      transactions.filter((t) => t.inventory_id === params.id)
    )
  }),

  http.post(
    `${BASE_URL}/inventories/:id/transactions`,
    async ({ request, params }) => {
      const body = (await request.json()) as CreateTransactionRequest
      // Mock full transaction object based on request
      const newTransaction: Transaction = {
        id: `tx-${Date.now()}`,
        inventory_id: params.id as string,
        outlet_id: body.outlet_id,
        created_by_user_id: '1', // Mock user
        transaction_date: body.transaction_date,
        created_at: new Date().toISOString(),
        items: body.items.map((item, index) => ({
          id: `tx-item-${Date.now()}-${index}`,
          transaction_id: `tx-${Date.now()}`,
          product_variant_id: item.product_variant_id,
          quantity: item.quantity,
          price_per_unit: item.price_per_unit,
          shopping_list_item_id: item.shopping_list_item_id,
        })),
      }

      transactions.push(newTransaction)

      if (body.items && Array.isArray(body.items)) {
        body.items.forEach((item) => {
          const invItem = inventoryProducts.find(
            (ip) =>
              ip.inventory_id === params.id &&
              ip.product_variant_id === item.product_variant_id
          )
          if (invItem) {
            invItem.quantity += item.quantity
          }
        })
      }

      return HttpResponse.json(newTransaction, { status: 201 })
    }
  ),
]
