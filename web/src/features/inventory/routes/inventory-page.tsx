import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getInventoryProducts } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Plus, Edit2, Utensils, Settings } from 'lucide-react'
import { AddInventoryItemDialog } from '../components/add-inventory-item-dialog'
import { AdjustInventoryItemDialog } from '../components/adjust-inventory-item-dialog'
import { ConsumeInventoryItemDialog } from '../components/consume-inventory-item-dialog'
import { InventorySettingsDialog } from '../components/inventory-settings-dialog'
import type { InventoryProductDetail } from '../types'

export function InventoryPage() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const [addDialogOpen, setAddDialogOpen] = useState(false)
  const [settingsDialogOpen, setSettingsDialogOpen] = useState(false)
  const [adjustItem, setAdjustItem] = useState<InventoryProductDetail | null>(
    null
  )
  const [consumeItem, setConsumeItem] = useState<InventoryProductDetail | null>(
    null
  )

  const {
    data: products,
    isLoading,
    error,
    refetch,
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
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Inventory</h1>
          <p className="text-gray-500">Manage your current stock levels.</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Button
            variant="outline"
            size="icon"
            onClick={() => setSettingsDialogOpen(true)}
            title="Settings"
          >
            <Settings className="h-4 w-4" />
          </Button>
          <Button
            onClick={() => setAddDialogOpen(true)}
            className="flex-1 sm:flex-initial"
          >
            <Plus className="mr-2 h-4 w-4" />
            Add Item
          </Button>
        </div>
      </div>

      {/* Desktop Table View */}
      <div className="hidden overflow-hidden rounded-lg border bg-white shadow-sm md:block">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-b bg-gray-50 font-medium text-gray-700">
              <tr>
                <th className="px-4 py-3">Product</th>
                <th className="px-4 py-3">Brand</th>
                <th className="px-4 py-3">Variant</th>
                <th className="px-4 py-3 text-right">Quantity</th>
                <th className="px-4 py-3">Unit</th>
                <th className="px-4 py-3 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {!products || products.length === 0 ? (
                <tr>
                  <td
                    colSpan={6}
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
                    <td className="px-4 py-3 text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => setConsumeItem(item)}
                          title="Consume"
                        >
                          <Utensils className="h-4 w-4 text-orange-500" />
                          <span className="sr-only">Consume</span>
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => setAdjustItem(item)}
                          title="Adjust"
                        >
                          <Edit2 className="h-4 w-4 text-gray-500" />
                          <span className="sr-only">Adjust</span>
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Mobile Card View */}
      <div className="grid gap-4 md:hidden">
        {!products || products.length === 0 ? (
          <div className="rounded-lg border border-dashed p-8 text-center text-gray-500">
            No items in inventory.
          </div>
        ) : (
          products?.map((item) => (
            <Card key={item.id}>
              <CardContent className="flex items-center justify-between p-4">
                <div className="space-y-1">
                  <div className="font-semibold text-midnight-slate">
                    {item.product_name}
                  </div>
                  <div className="text-sm text-slate-500">
                    {item.brand && (
                      <span className="mr-2 font-medium">{item.brand}</span>
                    )}
                    {item.variant_name && <span>{item.variant_name}</span>}
                  </div>
                  <div className="text-sm text-slate-400">
                    {item.size ? `${item.size} ${item.product_unit}` : ''}
                  </div>
                </div>

                <div className="flex flex-col items-end gap-3">
                  <div className="text-lg font-bold text-midnight-slate">
                    {item.quantity}{' '}
                    <span className="text-xs font-normal text-slate-500">
                      {item.unit || item.product_unit}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8 rounded-full border-soft-pebble text-orange-500 hover:text-orange-600 hover:bg-orange-50"
                      onClick={() => setConsumeItem(item)}
                    >
                      <Utensils className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8 rounded-full border-soft-pebble text-slate-500 hover:text-midnight-slate"
                      onClick={() => setAdjustItem(item)}
                    >
                      <Edit2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      <AddInventoryItemDialog
        open={addDialogOpen}
        onOpenChange={setAddDialogOpen}
        onSuccess={() => refetch()}
      />

      <AdjustInventoryItemDialog
        open={!!adjustItem}
        onOpenChange={(open) => !open && setAdjustItem(null)}
        item={adjustItem}
        onSuccess={() => refetch()}
      />

      <ConsumeInventoryItemDialog
        open={!!consumeItem}
        onOpenChange={(open) => !open && setConsumeItem(null)}
        item={consumeItem}
        onSuccess={() => refetch()}
      />

      <InventorySettingsDialog
        open={settingsDialogOpen}
        onOpenChange={setSettingsDialogOpen}
      />
    </div>
  )
}
