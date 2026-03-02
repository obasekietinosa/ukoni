import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getCanonicalProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'

interface Props {
  searchQuery?: string
}

export function CanonicalProductList({ searchQuery = '' }: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const [debouncedSearch, setDebouncedSearch] = useState('')

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchQuery)
    }, 500)
    return () => clearTimeout(timer)
  }, [searchQuery])

  const {
    data: products,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['canonical-products', activeInventoryId, debouncedSearch],
    queryFn: () =>
      getCanonicalProducts(activeInventoryId!, {
        search: debouncedSearch,
        limit: 100,
      }),
    enabled: !!activeInventoryId,
  })

  if (isLoading) return <div>Loading products...</div>
  if (error) return <div className="text-red-500">Failed to load products</div>

  if (!products || products.length === 0) {
    return (
      <div className="text-gray-500">
        {searchQuery
          ? 'No products match your search.'
          : 'No products found. Start by adding one.'}
      </div>
    )
  }

  return (
    <div className="space-y-2">
      {products.map((product) => (
        <Link key={product.id} to={`/products/${product.id}`} className="block">
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
  )
}
