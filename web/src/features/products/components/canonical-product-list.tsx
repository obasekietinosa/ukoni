import { useQuery } from '@tanstack/react-query'
import { getCanonicalProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'

export function CanonicalProductList() {
  const activeInventoryId = useInventoryStore((state) => state.activeInventoryId)

  const { data: products, isLoading, error } = useQuery({
    queryKey: ['canonical-products', activeInventoryId],
    queryFn: () => getCanonicalProducts(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  if (isLoading) return <div>Loading products...</div>
  if (error) return <div className="text-red-500">Failed to load products</div>

  if (!products || products.length === 0) {
      return <div className="text-gray-500">No products found. Start by adding one.</div>
  }

  return (
    <div className="space-y-2">
      {products.map((product) => (
        <div key={product.id} className="flex items-center justify-between rounded-lg border p-3 bg-white">
          <div>
            <div className="font-medium">{product.name}</div>
            {product.description && <div className="text-sm text-gray-500">{product.description}</div>}
          </div>
        </div>
      ))}
    </div>
  )
}
