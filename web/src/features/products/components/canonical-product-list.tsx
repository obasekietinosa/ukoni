import { Link } from 'react-router-dom'
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getCanonicalProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { Input } from '@/components/ui/input'
import { useDebounce } from '@/hooks/use-debounce'

export function CanonicalProductList() {
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 500)
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const {
    data: products,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['canonical-products', activeInventoryId, debouncedSearch],
    queryFn: () => getCanonicalProducts(activeInventoryId!, debouncedSearch),
    enabled: !!activeInventoryId,
  })

  if (isLoading) return <div>Loading products...</div>
  if (error) return <div className="text-red-500">Failed to load products</div>

  return (
    <div className="space-y-4">
      <div className="flex items-center space-x-2">
        <Input
          placeholder="Search products..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
      </div>

      {!products || products.length === 0 ? (
        <div className="text-gray-500">
          {search
            ? 'No products found matching your search.'
            : 'No products found. Start by adding one.'}
        </div>
      ) : (
        <div className="space-y-2">
          {products.map((product) => (
            <Link
              key={product.id}
              to={`/products/${product.id}`}
              className="block"
            >
              <div className="flex items-center justify-between rounded-lg border p-3 bg-white hover:bg-gray-50 transition-colors">
                <div>
                  <div className="font-medium">{product.name}</div>
                  {product.description && (
                    <div className="text-sm text-gray-500">
                      {product.description}
                    </div>
                  )}
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
