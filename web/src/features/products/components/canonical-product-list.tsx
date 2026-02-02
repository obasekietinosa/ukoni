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

  const {
    data: products,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['canonical-products', activeInventoryId],
    queryFn: () => getCanonicalProducts(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  if (isLoading) return <div>Loading products...</div>
  if (error) return <div className="text-red-500">Failed to load products</div>

  if (!products || products.length === 0) {
    return (
      <div className="text-gray-500">
        No products found. Start by adding one.
      </div>
    )
  }

  const filteredProducts = products.filter((product) => {
    const query = searchQuery.toLowerCase()
    return (
      product.name.toLowerCase().includes(query) ||
      product.description?.toLowerCase().includes(query)
    )
  })

  if (filteredProducts.length === 0) {
    return <div className="text-gray-500">No products match your search.</div>
  }

  return (
    <div className="space-y-2">
      {filteredProducts.map((product) => (
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
