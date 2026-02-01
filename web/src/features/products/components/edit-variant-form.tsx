import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { updateVariant } from '../api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import type { ProductVariant } from '../types'

interface Props {
  variant: ProductVariant
  onSuccess?: () => void
  onCancel?: () => void
}

export function EditVariantForm({ variant, onSuccess, onCancel }: Props) {
  const [variantName, setVariantName] = useState(variant.variant_name)
  const [sku, setSku] = useState(variant.sku || '')
  const [unit, setUnit] = useState(variant.unit || '')
  const [size, setSize] = useState(variant.size?.toString() || '')

  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: {
      variant_name: string
      sku?: string
      unit?: string
      size?: number
    }) => updateVariant(variant.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['variants', variant.product_id] })
      onSuccess?.()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!variantName) return
    mutation.mutate({
      variant_name: variantName,
      sku: sku || undefined,
      unit: unit || undefined,
      size: size ? parseFloat(size) : undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3 p-2 border rounded bg-gray-50">
      <h4 className="text-sm font-medium">Edit Variant</h4>
      {mutation.error && (
        <div className="text-red-500 text-xs">
          {mutation.error instanceof Error
            ? mutation.error.message
            : 'Failed to update variant'}
        </div>
      )}
      <div className="grid grid-cols-2 gap-2">
        <div className="col-span-2">
          <Input
            value={variantName}
            onChange={(e) => setVariantName(e.target.value)}
            placeholder="Variant Name (e.g. 1L Bottle)"
            required
            className="h-8 text-sm"
          />
        </div>
        <div>
          <Input
            value={sku}
            onChange={(e) => setSku(e.target.value)}
            placeholder="SKU (Optional)"
            className="h-8 text-sm"
          />
        </div>
        <div className="flex gap-2">
          <Input
            value={size}
            onChange={(e) => setSize(e.target.value)}
            placeholder="Size"
            type="number"
            step="0.01"
            className="h-8 text-sm w-1/2"
          />
          <Input
            value={unit}
            onChange={(e) => setUnit(e.target.value)}
            placeholder="Unit"
            className="h-8 text-sm w-1/2"
          />
        </div>
      </div>
      <div className="flex justify-end gap-2">
        <Button
          type="button"
          size="sm"
          variant="outline"
          onClick={onCancel}
          disabled={mutation.isPending}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          size="sm"
          disabled={mutation.isPending}
        >
          {mutation.isPending ? 'Saving...' : 'Save Changes'}
        </Button>
      </div>
    </form>
  )
}
