import { useState, useCallback } from 'react'
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
import { Loader2 } from 'lucide-react'
import { useAllOutlets } from '@/features/sellers/hooks/use-all-outlets'
import { createTransaction } from '../api'
import {
  TransactionWizardItem,
  type WizardItemState,
} from './transaction-wizard-item'
import type { ShoppingListItem } from '@/features/shopping-lists/types'
import { useInventoryStore } from '@/store/inventory'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  items: ShoppingListItem[]
  onSuccess?: () => void
}

export function CreateTransactionFromListDialog({
  open,
  onOpenChange,
  items,
  onSuccess,
}: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()
  const { outlets, isLoading: isLoadingOutlets } = useAllOutlets()

  // Transaction State
  const [transactionDate, setTransactionDate] = useState(
    new Date().toISOString().split('T')[0]
  )

  const [selectedOutletId, setSelectedOutletId] = useState<string>(() => {
    if (items.length > 0) {
      const preferred = items.find(
        (i) => i.preferred_outlet_id
      )?.preferred_outlet_id
      if (preferred) return preferred
    }
    return ''
  })

  const [itemStates, setItemStates] = useState<Record<string, WizardItemState>>(
    {}
  )

  const handleItemUpdate = useCallback((state: WizardItemState) => {
    setItemStates((prev) => ({
      ...prev,
      [state.shoppingListItemId]: state,
    }))
  }, [])

  const createMutation = useMutation({
    mutationFn: async () => {
      const transactionItems = Object.values(itemStates)
        .filter((s) => s.included && s.selectedVariantId) // Only included and valid variants
        .map((s) => ({
          product_variant_id: s.selectedVariantId!,
          quantity: s.quantity,
          price_per_unit: s.price || undefined,
          shopping_list_item_id: s.shoppingListItemId,
        }))

      if (transactionItems.length === 0) return // Nothing to submit

      return createTransaction(activeInventoryId!, {
        transaction_date: new Date(transactionDate).toISOString(),
        outlet_id: selectedOutletId || undefined,
        items: transactionItems,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['inventory-products', activeInventoryId],
      })
      queryClient.invalidateQueries({ queryKey: ['shopping-list-items'] })
      onSuccess?.()
      onOpenChange(false)
    },
  })

  // Validation
  const includedItems = Object.values(itemStates).filter((s) => s.included)
  const isValid =
    includedItems.length > 0 &&
    includedItems.every((s) => !!s.selectedVariantId)

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl max-h-[90vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>Complete Shopping Trip</DialogTitle>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto py-4 space-y-6">
          {/* Meta */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Date</label>
              <Input
                type="date"
                value={transactionDate}
                onChange={(e) => setTransactionDate(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Outlet</label>
              <select
                value={selectedOutletId}
                onChange={(e) => setSelectedOutletId(e.target.value)}
                disabled={isLoadingOutlets}
                className="flex h-10 w-full rounded-md border border-gray-200 bg-white px-3 py-2 text-base md:text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-950 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              >
                <option value="">Select Store (Optional)</option>
                {outlets.map((outlet) => (
                  <option key={outlet.id} value={outlet.id}>
                    {outlet.sellerName} - {outlet.name || 'Main'}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Items */}
          <div className="space-y-4">
            <h3 className="font-medium">Items</h3>
            <div className="space-y-3">
              {items.map((item) => (
                <TransactionWizardItem
                  key={item.id}
                  item={item}
                  onUpdate={handleItemUpdate}
                />
              ))}
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={() => createMutation.mutate()}
            disabled={!isValid || createMutation.isPending}
          >
            {createMutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            Complete Transaction
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
