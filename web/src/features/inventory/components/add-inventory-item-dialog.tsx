import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createTransaction } from '../api'
import { Modal } from '@/components/ui/modal'
import { Button } from '@/components/ui/button'
import { useInventoryStore } from '@/store/inventory'

interface Props {
  isOpen: boolean
  onClose: () => void
  variantId: string
  variantName: string
  unit?: string
}

export function AddInventoryItemDialog({
  isOpen,
  onClose,
  variantId,
  variantName,
  unit,
}: Props) {
  const [quantity, setQuantity] = useState<string>('1')
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: async () => {
      if (!activeInventoryId) return
      return createTransaction(activeInventoryId, {
        transaction_date: new Date().toISOString(),
        items: [
          {
            product_variant_id: variantId,
            quantity: parseFloat(quantity),
          },
        ],
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['inventory-products', activeInventoryId],
      })
      onClose()
      setQuantity('1')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!quantity || isNaN(parseFloat(quantity))) return
    mutation.mutate()
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={`Add ${variantName} to Inventory`}
    >
      <form id="add-inventory-form" onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="quantity"
            className="block text-sm font-medium text-gray-700"
          >
            Quantity {unit ? `(${unit})` : ''}
          </label>
          <input
            type="number"
            id="quantity"
            min="0.01"
            step="any"
            required
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
          />
        </div>
        {mutation.error && (
            <div className="text-sm text-red-500">
                Failed to add item.
            </div>
        )}
      </form>
      <div className="mt-4 flex justify-end gap-2">
        <Button variant="outline" onClick={onClose} type="button">
            Cancel
        </Button>
        <Button type="submit" form="add-inventory-form" disabled={mutation.isPending}>
            {mutation.isPending ? 'Adding...' : 'Add to Inventory'}
        </Button>
      </div>
    </Modal>
  )
}
