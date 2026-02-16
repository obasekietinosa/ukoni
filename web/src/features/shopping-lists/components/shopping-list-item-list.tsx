import { useState } from 'react'
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

  // State for ticked items
  const [checkedItems, setCheckedItems] = useState<Set<string>>(() => {
    try {
      const saved = localStorage.getItem(`shopping-list-checked-${listId}`)
      return saved ? new Set(JSON.parse(saved)) : new Set()
    } catch (e) {
      console.error('Failed to parse checked items from local storage', e)
      return new Set()
    }
  })

  const updateLocalStorage = (newCheckedItems: Set<string>) => {
    try {
      localStorage.setItem(
        `shopping-list-checked-${listId}`,
        JSON.stringify(Array.from(newCheckedItems))
      )
    } catch (e) {
      console.error('Failed to save checked items to local storage', e)
    }
  }

  const toggleItem = (itemId: string) => {
    const next = new Set(checkedItems)
    if (next.has(itemId)) {
      next.delete(itemId)
    } else {
      next.add(itemId)
    }
    setCheckedItems(next)
    updateLocalStorage(next)
  }

  const deleteMutation = useMutation({
    mutationFn: deleteShoppingListItem,
    onSuccess: (_, itemId) => {
      // Cleanup local storage for deleted item
      const next = new Set(checkedItems)
      if (next.has(itemId)) {
        next.delete(itemId)
        setCheckedItems(next)
        updateLocalStorage(next)
      }
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
      {items.map((item) => {
        const isChecked = checkedItems.has(item.id)
        return (
          <div
            key={item.id}
            className={`flex items-center justify-between rounded-lg border p-3 bg-white transition-colors duration-200 ${
              isChecked ? 'bg-gray-50' : ''
            }`}
          >
            <div className="mr-3">
              <input
                type="checkbox"
                checked={isChecked}
                onChange={() => toggleItem(item.id)}
                className="h-5 w-5 rounded border-gray-300 text-electric-mint focus:ring-electric-mint cursor-pointer accent-electric-mint"
                aria-label={`Mark ${getItemName(item)} as purchased`}
              />
            </div>
            <div className={`flex-1 ${isChecked ? 'opacity-50' : ''}`}>
              <div className="flex items-baseline gap-2">
                <span
                  className={`font-medium ${
                    isChecked ? 'line-through text-gray-500' : ''
                  }`}
                >
                  {getItemName(item)}
                </span>
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
        )
      })}
    </div>
  )
}
