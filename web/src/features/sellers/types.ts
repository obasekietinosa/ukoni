export interface Seller {
  id: string
  name: string
  type: 'chain' | 'independent' | 'online'
  created_at: string
  deleted_at?: string
}

export interface Outlet {
  id: string
  seller_id: string
  name: string
  channel: 'physical' | 'online'
  address?: string
  website_url?: string
  created_at: string
  deleted_at?: string
}
