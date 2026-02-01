import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getVariants, deleteVariant } from '../api'
import { Button } from '@/components/ui/button'
import { EditVariantForm } from './edit-variant-form'

interface Props {
  productId: string
}

export function VariantList({ productId }: Props) {
  const [editingVariantId, setEditingVariantId] = useState<string | null>(null)
  const queryClient = useQueryClient()

  const {
    data: variants,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['variants', productId],
    queryFn: () => getVariants(productId),
  })

  const deleteMutation = useMutation({
    mutationFn: deleteVariant,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['variants', productId] })
    },
  })

  if (isLoading)
    return <div className="text-xs text-gray-500">Loading variants...</div>
  if (error)
    return <div className="text-xs text-red-500">Failed to load variants</div>

  if (!variants || variants.length === 0) {
    return <div className="text-sm text-gray-400 italic">No variants yet.</div>
  }

  return (
    <div className="space-y-2">
      {variants.map((variant) => {
        if (editingVariantId === variant.id) {
          return (
            <EditVariantForm
              key={variant.id}
              variant={variant}
              onSuccess={() => setEditingVariantId(null)}
              onCancel={() => setEditingVariantId(null)}
            />
          )
        }

        return (
          <div
            key={variant.id}
            className="flex items-center justify-between text-sm p-2 bg-gray-50 rounded group"
          >
            <div>
              <div className="font-medium">{variant.variant_name}</div>
              <div className="text-gray-500 text-xs">
                {variant.size && variant.unit
                  ? `${variant.size} ${variant.unit}`
                  : ''}
                {variant.sku ? ` • SKU: ${variant.sku}` : ''}
              </div>
            </div>
            <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <Button
                variant="ghost"
                size="sm"
                className="h-6 text-xs px-2"
                onClick={() => setEditingVariantId(variant.id)}
              >
                Edit
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="h-6 text-xs px-2 text-red-500 hover:text-red-700"
                onClick={() => {
                  if (confirm('Delete this variant?')) {
                    deleteMutation.mutate(variant.id)
                  }
                }}
              >
                Delete
              </Button>
            </div>
          </div>
        )
      })}
    </div>
  )
}
