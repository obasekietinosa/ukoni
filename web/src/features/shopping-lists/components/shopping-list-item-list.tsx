import type { ShoppingListItem } from '../types'
import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { deleteShoppingListItem } from '../api'

interface Props {
  items: ShoppingListItem[]
  listId: string
}

export function ShoppingListItemList({ items, listId }: Props) {
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: deleteShoppingListItem,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['shopping-list-items', listId],
      })
    },
  })

  if (items.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        List is empty. Add some items!
      </div>
    )
  }

  const getItemName = (item: ShoppingListItem) => {
    if (item.target_type === 'canonical_product') {
      return item.canonical_product?.name || 'Unknown Product'
    }
    if (item.target_type === 'product_variant') {
      const brand = item.product?.brand
      const name = item.product_variant?.variant_name || 'Unknown Variant'
      return brand ? `${brand} ${name}` : name
    }
    return 'Unknown Item'
  }

  return (
    <div className="space-y-2">
      {items.map((item) => (
        <div
          key={item.id}
          className="flex items-center justify-between rounded-lg border p-3 bg-white"
        >
          <div className="flex-1">
            <div className="flex items-baseline gap-2">
              <span className="font-medium">{getItemName(item)}</span>
              {(item.quantity || item.unit) && (
                <span className="text-sm text-gray-500">
                  {item.quantity} {item.unit}
                </span>
              )}
            </div>
            {item.notes && (
              <div className="text-sm text-gray-500 italic mt-1">
                "{item.notes}"
              </div>
            )}
            {item.preferred_outlet && (
              <div className="text-xs text-blue-600 mt-1">
                @ {item.preferred_outlet.name}
              </div>
            )}
          </div>
          <Button
            variant="ghost"
            size="sm"
            className="text-red-500 hover:text-red-700 hover:bg-red-50"
            onClick={() => {
              if (confirm('Remove this item?')) {
                deleteMutation.mutate(item.id)
              }
            }}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      ))}
    </div>
  )
}
