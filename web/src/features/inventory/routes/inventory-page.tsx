import { useQuery } from '@tanstack/react-query'
import { getInventoryProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'

export function InventoryPage() {
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

  if (isLoading) {
    return (
      <div className="p-8 text-center text-gray-500">Loading inventory...</div>
    )
  }

  if (error) {
    return (
      <div className="p-8 text-center text-red-500">
        Error loading inventory: {(error as Error).message}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Inventory</h1>
        <p className="text-gray-500">Manage your current stock levels.</p>
      </div>

      <div className="overflow-hidden rounded-lg border bg-white shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-b bg-gray-50 font-medium text-gray-700">
              <tr>
                <th className="px-4 py-3">Product</th>
                <th className="px-4 py-3">Brand</th>
                <th className="px-4 py-3">Variant</th>
                <th className="px-4 py-3 text-right">Quantity</th>
                <th className="px-4 py-3">Unit</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {!products || products.length === 0 ? (
                <tr>
                  <td
                    colSpan={5}
                    className="px-4 py-8 text-center text-gray-500"
                  >
                    No items in inventory.
                  </td>
                </tr>
              ) : (
                products?.map((item) => (
                  <tr
                    key={item.id}
                    className="transition-colors hover:bg-gray-50"
                  >
                    <td className="px-4 py-3 font-medium text-gray-900">
                      {item.product_name}
                    </td>
                    <td className="px-4 py-3 text-gray-600">
                      {item.brand || '-'}
                    </td>
                    <td className="px-4 py-3 text-gray-600">
                      {item.variant_name}
                      {item.size ? ` (${item.size} ${item.product_unit})` : ''}
                    </td>
                    <td className="px-4 py-3 text-right font-medium text-gray-900">
                      {item.quantity}
                    </td>
                    <td className="px-4 py-3 text-gray-600">
                      {item.unit || item.product_unit || '-'}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
