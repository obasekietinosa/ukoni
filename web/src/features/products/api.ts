import { api } from '@/lib/api'
import type { CanonicalProduct, Product, ProductVariant } from './types'

export const getCanonicalProducts = async (
  inventoryId: string,
  search?: string
): Promise<CanonicalProduct[]> => {
  const params = new URLSearchParams()
  if (search) {
    params.append('search', search)
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
