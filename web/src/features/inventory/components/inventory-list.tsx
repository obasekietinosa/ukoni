import { useQuery } from '@tanstack/react-query'
import { getInventoryProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'

export function InventoryList() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const {
    data: products,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['inventory-products', activeInventoryId],
    queryFn: () => getInventoryProducts(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  if (isLoading) return <div>Loading inventory...</div>
  if (error) return <div className="text-red-500">Failed to load inventory</div>

  if (!products || products.length === 0) {
    return (
      <div className="text-gray-500">
        Your inventory is empty. Start by adding items from the Product Catalog.
      </div>
    )
  }

  return (
    <div className="rounded-lg border bg-white shadow-sm overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Product
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Brand / Variant
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Quantity
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {products.map((item) => (
            <tr key={item.id}>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {item.canonical_product_name}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {item.brand_name && (
                    <span className="block font-medium">{item.brand_name}</span>
                )}
                <span>{item.variant_name}</span>
                {item.sku && <span className="ml-2 text-xs text-gray-400">({item.sku})</span>}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-medium">
                {item.quantity} {item.unit || ''}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
