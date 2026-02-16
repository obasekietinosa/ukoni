import type { PlanItem } from '../types'
import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { removePlanItem } from '../api'

interface Props {
  items: PlanItem[]
  planId: string
}

export function PlanItemList({ items, planId }: Props) {
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: removePlanItem,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['plan', planId],
      })
    },
  })

  if (items.length === 0) {
    return (
      <div className="text-center py-12 rounded-lg border border-dashed border-slate-200 bg-slate-50/50">
        <p className="text-sm text-slate-500">
          No items in this plan yet.
        </p>
      </div>
    )
  }

  const getItemName = (item: PlanItem) => {
    if (item.target_type === 'canonical_product') {
      return item.canonical_product?.name || 'Unknown Product'
    }
    if (item.target_type === 'product_variant') {
      const brand = item.product?.brand
      const name = item.product_variant?.variant_name || 'Unknown Variant'
      return brand ? `${brand} ${name}` : name
    }
    // Fallback for direct product reference if implemented
    if (item.target_type === 'product') {
        return item.product?.name || 'Unknown Product'
    }
    return 'Unknown Item'
  }

  return (
    <div className="divide-y divide-slate-100 border rounded-lg overflow-hidden bg-white shadow-sm">
      {items.map((item) => (
        <div
          key={item.id}
          className="flex items-center justify-between p-4 hover:bg-slate-50 transition-colors"
        >
          <div className="flex-1 min-w-0 pr-4">
            <div className="flex items-baseline gap-2 flex-wrap">
              <span className="font-medium text-slate-900 truncate">
                {getItemName(item)}
              </span>
              {(item.quantity || item.unit) && (
                <span className="text-sm font-medium text-electric-mint px-2 py-0.5 rounded-full bg-electric-mint/10">
                  {item.quantity} {item.unit}
                </span>
              )}
            </div>
            {item.note && (
              <div className="text-sm text-slate-500 mt-1 italic">
                "{item.note}"
              </div>
            )}
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-full"
              onClick={() => {
                if (confirm('Remove this item?')) {
                  deleteMutation.mutate(item.id)
                }
              }}
              title="Remove item"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}
