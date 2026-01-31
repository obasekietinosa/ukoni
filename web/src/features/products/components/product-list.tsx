import { useQuery } from '@tanstack/react-query'
import { getProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'
import type { Product } from '../types'

interface Props {
  canonicalProductId: string
  children?: (product: Product) => React.ReactNode
}

export function ProductList({ canonicalProductId, children }: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const {
    data: allProducts,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['products', activeInventoryId],
    queryFn: () => getProducts(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  // Filter client-side
  const products = allProducts?.filter(
    (p) => p.canonical_product_id === canonicalProductId
  )

  if (isLoading) return <div>Loading brands...</div>
  if (error) return <div className="text-red-500">Failed to load brands</div>

  if (!products || products.length === 0) {
    return <div className="text-gray-500">No brands found. Add one below.</div>
  }

  return (
    <div className="space-y-4">
      {products.map((product) => (
        <div key={product.id} className="rounded-lg border p-4 bg-white">
          <div className="flex items-center justify-between mb-2">
            <div>
              <div className="font-medium">{product.name}</div>
              {product.brand && (
                <div className="text-sm text-gray-500">{product.brand}</div>
              )}
            </div>
          </div>
          {children && children(product)}
        </div>
      ))}
    </div>
  )
}
