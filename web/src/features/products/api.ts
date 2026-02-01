import { api } from '@/lib/api'
import type { CanonicalProduct, Product, ProductVariant } from './types'

export const getCanonicalProducts = async (
  inventoryId: string,
  options: {
    categoryId?: string
    search?: string
    limit?: number
    offset?: number
  } = {}
): Promise<CanonicalProduct[]> => {
  const params = new URLSearchParams()
  if (options.categoryId) {
    params.append('category_id', options.categoryId)
  }
  if (options.search) {
    params.append('search', options.search)
  }
  if (options.limit) {
    params.append('limit', options.limit.toString())
  }
  if (options.offset) {
    params.append('offset', options.offset.toString())
  }
  return api<CanonicalProduct[]>(
    `/inventories/${inventoryId}/canonical-products?${params.toString()}`
  )
}

export const getCanonicalProduct = async (
  id: string
): Promise<CanonicalProduct> => {
  return api<CanonicalProduct>(`/canonical-products/${id}`)
}

export const createCanonicalProduct = async (
  inventoryId: string,
  data: { name: string; description?: string; category_id?: string }
): Promise<CanonicalProduct> => {
  return api<CanonicalProduct>(
    `/inventories/${inventoryId}/canonical-products`,
    {
      method: 'POST',
      json: data,
    }
  )
}

export const updateCanonicalProduct = async (
  id: string,
  data: { name: string; description?: string; category_id?: string }
): Promise<CanonicalProduct> => {
  return api<CanonicalProduct>(`/canonical-products/${id}`, {
    method: 'PUT',
    json: data,
  })
}

export const deleteCanonicalProduct = async (id: string): Promise<void> => {
  return api<void>(`/canonical-products/${id}`, {
    method: 'DELETE',
  })
}

export const getProducts = async (
  inventoryId: string,
  canonicalProductId?: string
): Promise<Product[]> => {
  const params = new URLSearchParams()
  if (canonicalProductId) {
    params.append('canonical_product_id', canonicalProductId)
  }
  return api<Product[]>(
    `/inventories/${inventoryId}/products?${params.toString()}`
  )
}

export const createProduct = async (
  inventoryId: string,
  data: {
    name: string
    brand?: string
    description?: string
    category_id?: string
    canonical_product_id?: string
  }
): Promise<Product> => {
  return api<Product>(`/inventories/${inventoryId}/products`, {
    method: 'POST',
    json: data,
  })
}

export const updateProduct = async (
  id: string,
  data: {
    name: string
    brand?: string
    description?: string
    category_id?: string
    canonical_product_id?: string
  }
): Promise<Product> => {
  return api<Product>(`/products/${id}`, {
    method: 'PUT',
    json: data,
  })
}

export const deleteProduct = async (id: string): Promise<void> => {
  return api<void>(`/products/${id}`, {
    method: 'DELETE',
  })
}

export const getVariants = async (
  productId: string
): Promise<ProductVariant[]> => {
  return api<ProductVariant[]>(`/products/${productId}/variants`)
}

export const createVariant = async (
  productId: string,
  data: {
    variant_name: string
    sku?: string
    unit?: string
    size?: number
  }
): Promise<ProductVariant> => {
  return api<ProductVariant>(`/products/${productId}/variants`, {
    method: 'POST',
    json: data,
  })
}

export const updateVariant = async (
  id: string,
  data: {
    variant_name: string
    sku?: string
    unit?: string
    size?: number
  }
): Promise<ProductVariant> => {
  return api<ProductVariant>(`/product-variants/${id}`, {
    method: 'PUT',
    json: data,
  })
}

export const deleteVariant = async (id: string): Promise<void> => {
  return api<void>(`/product-variants/${id}`, {
    method: 'DELETE',
  })
}
