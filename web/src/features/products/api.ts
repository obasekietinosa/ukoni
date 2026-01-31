import { api } from '@/lib/api'
import type { CanonicalProduct, Product, ProductVariant } from './types'

export const getCanonicalProducts = async (
  inventoryId: string
): Promise<CanonicalProduct[]> => {
  return api<CanonicalProduct[]>(
    `/inventories/${inventoryId}/canonical-products`
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

export const getProducts = async (inventoryId: string): Promise<Product[]> => {
  return api<Product[]>(`/inventories/${inventoryId}/products`)
}

export const createProduct = async (
  inventoryId: string,
  data: {
    name: string
    canonical_product_id?: string
    brand?: string
    description?: string
    category_id?: string
  }
): Promise<Product> => {
  return api<Product>(`/inventories/${inventoryId}/products`, {
    method: 'POST',
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
