import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useInventoryStore } from '@/store/inventory'
import { createTransaction } from '@/features/transactions/api'
import { createConsumptionEvent } from '@/features/consumption/api'
import type { InventoryProductDetail } from '@/features/inventory/types'
import { Loader2 } from 'lucide-react'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  item: InventoryProductDetail | null
  onSuccess?: () => void
}

function ConsumeInventoryItemForm({
  item,
  onClose,
  onSuccess,
}: {
  item: InventoryProductDetail
  onClose: () => void
  onSuccess?: () => void
}) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()
  const [consumeQuantity, setConsumeQuantity] = useState('1')

  const consumeQty = parseFloat(consumeQuantity || '0')

  const consumeMutation = useMutation({
    mutationFn: async () => {
      // 1. Record Consumption
      await createConsumptionEvent(activeInventoryId!, {
        canonical_product_id: item.canonical_product_id,
        product_variant_id: item.product_variant_id,
        quantity: consumeQty,
        unit: item.unit || item.product_unit,
        source: 'quick-action',
        consumed_at: new Date().toISOString(),
      })

      // 2. Reduce Inventory via Transaction
      return createTransaction(activeInventoryId!, {
        transaction_date: new Date().toISOString(),
        items: [
          {
            product_variant_id: item.product_variant_id,
            quantity: -consumeQty, // Negative for consumption
          },
        ],
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['inventory-products', activeInventoryId],
      })
      queryClient.invalidateQueries({
        queryKey: ['consumption-events', activeInventoryId],
      })
      onSuccess?.()
      onClose()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (consumeQty <= 0) {
      return
    }
    consumeMutation.mutate()
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="text-sm text-gray-500">
        Consuming:{' '}
        <span className="font-medium text-gray-900">
          {item.product_name} - {item.variant_name}
        </span>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="text-sm font-medium text-gray-500">Current</label>
          <div className="p-2 bg-gray-50 rounded-md text-gray-900 font-medium">
            {item.quantity} {item.unit || item.product_unit}
          </div>
        </div>
        <div className="space-y-2">
          <label htmlFor="consume-quantity" className="text-sm font-medium">
            Quantity to Consume
          </label>
          <Input
            id="consume-quantity"
            type="number"
            step="0.01"
            min="0.01"
            max={item.quantity > 0 ? item.quantity : undefined}
            value={consumeQuantity}
            onChange={(e) => setConsumeQuantity(e.target.value)}
            required
            autoFocus
          />
        </div>
      </div>

      <div className="text-sm font-medium text-red-600">
        Will remove {consumeQty} {item.unit || item.product_unit} from inventory.
      </div>

      <DialogFooter>
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={consumeMutation.isPending || consumeQty <= 0}
        >
          {consumeMutation.isPending && (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          )}
          Confirm Consumption
        </Button>
      </DialogFooter>
    </form>
  )
}

export function ConsumeInventoryItemDialog({
  open,
  onOpenChange,
  item,
  onSuccess,
}: Props) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Consume Item</DialogTitle>
        </DialogHeader>
        {item && open && (
          <ConsumeInventoryItemForm
            item={item}
            onClose={() => onOpenChange(false)}
            onSuccess={onSuccess}
          />
        )}
      </DialogContent>
    </Dialog>
  )
}
