import { useQuery } from '@tanstack/react-query'
import { getProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { VariantList } from './variant-list'
import { CreateVariantForm } from './create-variant-form'

interface Props {
  canonicalProductId: string
}

export function ProductList({ canonicalProductId }: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const {
    data: products,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['products', activeInventoryId, canonicalProductId],
    queryFn: () => getProducts(activeInventoryId!, canonicalProductId),
    enabled: !!activeInventoryId && !!canonicalProductId,
  })

  if (isLoading) return <div>Loading brands...</div>
  if (error) return <div className="text-red-500">Failed to load brands</div>

  if (!products || products.length === 0) {
    return (
      <div className="text-gray-500">No brands found for this product.</div>
    )
  }

  return (
    <div className="space-y-4">
      {products.map((product) => (
        <div key={product.id} className="rounded-lg border p-4 bg-white">
          <div className="flex items-center justify-between mb-4">
            <div>
              <div className="font-medium text-lg">{product.name}</div>
              {product.brand && (
                <div className="text-sm text-gray-500">
                  Brand: {product.brand}
                </div>
              )}
            </div>
          </div>

          <div className="pl-4 border-l-2 border-gray-100">
            <div className="text-sm text-gray-500 mb-2">Variants:</div>
            <VariantList productId={product.id} />
            <CreateVariantForm productId={product.id} />
          </div>
        </div>
      ))}
    </div>
  )
}
