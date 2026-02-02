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
import type { InventoryProductDetail } from '@/features/inventory/types'
import { Loader2 } from 'lucide-react'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  item: InventoryProductDetail | null
  onSuccess?: () => void
}

function AdjustInventoryItemForm({
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
  const [newQuantity, setNewQuantity] = useState(item.quantity.toString())

  const diff = parseFloat(newQuantity || '0') - item.quantity

  const adjustMutation = useMutation({
    mutationFn: () =>
      createTransaction(activeInventoryId!, {
        transaction_date: new Date().toISOString(),
        items: [
          {
            product_variant_id: item.product_variant_id,
            quantity: diff,
          },
        ],
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['inventory-products', activeInventoryId],
      })
      onSuccess?.()
      onClose()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (diff === 0) {
      onClose()
      return
    }
    adjustMutation.mutate()
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="text-sm text-gray-500">
        Adjusting:{' '}
        <span className="font-medium text-gray-900">
          {item.product_name} - {item.variant_name}
        </span>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="text-sm font-medium text-gray-500">Current</label>
          <div className="p-2 bg-gray-50 rounded-md text-gray-900 font-medium">
            {item.quantity} {item.unit}
          </div>
        </div>
        <div className="space-y-2">
          <label htmlFor="new-quantity" className="text-sm font-medium">
            New Quantity
          </label>
          <Input
            id="new-quantity"
            type="number"
            step="0.01"
            min="0"
            value={newQuantity}
            onChange={(e) => setNewQuantity(e.target.value)}
            required
            autoFocus
          />
        </div>
      </div>

      {diff !== 0 && (
        <div
          className={`text-sm font-medium ${diff > 0 ? 'text-green-600' : 'text-red-600'}`}
        >
          Will {diff > 0 ? 'add' : 'remove'} {Math.abs(diff).toFixed(2)}{' '}
          {item.unit}
        </div>
      )}

      <DialogFooter>
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={adjustMutation.isPending || diff === 0}>
          {adjustMutation.isPending && (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          )}
          Confirm Adjustment
        </Button>
      </DialogFooter>
    </form>
  )
}

export function AdjustInventoryItemDialog({
  open,
  onOpenChange,
  item,
  onSuccess,
}: Props) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Adjust Inventory</DialogTitle>
        </DialogHeader>
        {item && (
          <AdjustInventoryItemForm
            key={item.id}
            item={item}
            onClose={() => onOpenChange(false)}
            onSuccess={onSuccess}
          />
        )}
      </DialogContent>
    </Dialog>
  )
}
