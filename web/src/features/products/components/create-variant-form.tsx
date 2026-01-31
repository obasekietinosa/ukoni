import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createVariant } from '../api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

interface Props {
  productId: string
  onSuccess?: () => void
}

export function CreateVariantForm({ productId, onSuccess }: Props) {
  const [variantName, setVariantName] = useState('')
  const [sku, setSku] = useState('')
  const [unit, setUnit] = useState('')
  const [size, setSize] = useState('')

  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: {
      variant_name: string
      sku?: string
      unit?: string
      size?: number
    }) => createVariant(productId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['variants', productId] })
      setVariantName('')
      setSku('')
      setUnit('')
      setSize('')
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
    <form onSubmit={handleSubmit} className="space-y-3 mt-4 border-t pt-4">
      <h4 className="text-sm font-medium">Add Variant</h4>
      {mutation.error && (
        <div className="text-red-500 text-xs">
          {mutation.error instanceof Error
            ? mutation.error.message
            : 'Failed to create variant'}
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
      <Button
        type="submit"
        size="sm"
        variant="secondary"
        disabled={mutation.isPending}
        className="w-full"
      >
        {mutation.isPending ? 'Adding...' : 'Add Variant'}
      </Button>
    </form>
  )
}
