import { useQuery } from '@tanstack/react-query'
import { getVariants } from '../api'
import { Button } from '@/components/ui/button'
import { useState } from 'react'
import { EditVariantDialog } from './edit-variant-dialog'
import type { ProductVariant } from '../types'

interface Props {
  productId: string
}

export function VariantList({ productId }: Props) {
  const [editingVariant, setEditingVariant] = useState<ProductVariant | null>(
    null
  )
  const {
    data: variants,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['variants', productId],
    queryFn: () => getVariants(productId),
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
      {editingVariant && (
        <EditVariantDialog
          variant={editingVariant}
          open={!!editingVariant}
          onClose={() => setEditingVariant(null)}
        />
      )}
      {variants.map((variant) => (
        <div
          key={variant.id}
          className="flex items-center justify-between text-sm p-2 bg-gray-50 rounded"
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
          <Button
            variant="ghost"
            size="sm"
            className="h-6 px-2 text-xs"
            onClick={() => setEditingVariant(variant)}
          >
            Edit
          </Button>
        </div>
      ))}
    </div>
  )
}
