import { useState } from 'react'
import { CanonicalProductList } from '../components/canonical-product-list'
import { CreateCanonicalProductForm } from '../components/create-canonical-product-form'
import { Input } from '@/components/ui/input'

export function ProductsPage() {
  const [searchQuery, setSearchQuery] = useState('')

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Product Catalog</h1>
        <p className="text-gray-500">Manage your canonical products here.</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <div>
          <CreateCanonicalProductForm />
        </div>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold">All Products</h2>
          </div>
          <Input
            type="search"
            placeholder="Search products..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="max-w-md"
          />
          <CanonicalProductList searchQuery={searchQuery} />
        </div>
      </div>
    </div>
  )
}
