import { api } from '@/lib/api'
import type { Outlet, Seller } from './types'

// Sellers
export const getSellers = async () => {
  return api<Seller[]>('/sellers')
}

export const getSeller = async (id: string) => {
  return api<Seller>(`/sellers/${id}`)
}

export const createSeller = async (data: { name: string; type: Seller['type'] }) => {
  return api<Seller>('/sellers', {
    method: 'POST',
    json: data,
  })
}

export const updateSeller = async (
  id: string,
  data: { name: string; type: Seller['type'] }
) => {
  return api<Seller>(`/sellers/${id}`, {
    method: 'PUT',
    json: data,
  })
}

export const deleteSeller = async (id: string) => {
  return api<void>(`/sellers/${id}`, {
    method: 'DELETE',
  })
}

// Outlets
export const getOutlets = async (sellerId: string) => {
  return api<Outlet[]>(`/sellers/${sellerId}/outlets`)
}

export const getOutlet = async (id: string) => {
    return api<Outlet>(`/outlets/${id}`)
}

export const createOutlet = async (
  sellerId: string,
  data: {
    name: string
    channel: Outlet['channel']
    address?: string
    website_url?: string
  }
) => {
  return api<Outlet>(`/sellers/${sellerId}/outlets`, {
    method: 'POST',
    json: data,
  })
}

export const updateOutlet = async (
  id: string,
  data: {
    name: string
    channel: Outlet['channel']
    address?: string
    website_url?: string
  }
) => {
  return api<Outlet>(`/outlets/${id}`, {
    method: 'PUT',
    json: data,
  })
}

export const deleteOutlet = async (id: string) => {
  return api<void>(`/outlets/${id}`, {
    method: 'DELETE',
  })
}
