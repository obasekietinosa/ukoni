import { useState, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { updateVariant } from '../api'
import type { ProductVariant } from '../types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

interface Props {
  variant: ProductVariant
  open: boolean
  onClose: () => void
}

export function EditVariantDialog({ variant, open, onClose }: Props) {
  const [variantName, setVariantName] = useState(variant.variant_name)
  const [sku, setSku] = useState(variant.sku || '')
  const [unit, setUnit] = useState(variant.unit || '')
  const [size, setSize] = useState(variant.size ? variant.size.toString() : '')
  const queryClient = useQueryClient()

  useEffect(() => {
    if (open) {
      setVariantName(variant.variant_name)
      setSku(variant.sku || '')
      setUnit(variant.unit || '')
      setSize(variant.size ? variant.size.toString() : '')
    }
  }, [open, variant])

  const mutation = useMutation({
    mutationFn: (data: {
      variant_name: string
      sku?: string
      unit?: string
      size?: number
    }) => updateVariant(variant.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['variants', variant.product_id],
      })
      onClose()
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

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
        <h3 className="text-lg font-medium mb-4">Edit Variant</h3>
        {mutation.error && (
          <div className="text-red-500 text-sm mb-4">
            {mutation.error instanceof Error
              ? mutation.error.message
              : 'Failed to update variant'}
          </div>
        )}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="edit-variant-name" className="text-sm font-medium">
              Variant Name
            </label>
            <Input
              id="edit-variant-name"
              type="text"
              value={variantName}
              onChange={(e) => setVariantName(e.target.value)}
              placeholder="e.g. 1L Bottle"
              required
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label htmlFor="edit-size" className="text-sm font-medium">
                Size
              </label>
              <Input
                id="edit-size"
                type="number"
                step="0.01"
                value={size}
                onChange={(e) => setSize(e.target.value)}
                placeholder="e.g. 1.0"
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="edit-unit" className="text-sm font-medium">
                Unit
              </label>
              <Input
                id="edit-unit"
                type="text"
                value={unit}
                onChange={(e) => setUnit(e.target.value)}
                placeholder="e.g. L, kg"
              />
            </div>
          </div>
          <div className="space-y-2">
            <label htmlFor="edit-sku" className="text-sm font-medium">
              SKU (Optional)
            </label>
            <Input
              id="edit-sku"
              type="text"
              value={sku}
              onChange={(e) => setSku(e.target.value)}
            />
          </div>
          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
