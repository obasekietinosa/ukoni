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
  const [size, setSize] = useState('')
  const [unit, setUnit] = useState('')
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: {
      variant_name: string
      size?: number
      unit?: string
    }) => createVariant(productId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['variants', productId] })
      setVariantName('')
      setSize('')
      setUnit('')
      onSuccess?.()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!variantName) return
    mutation.mutate({
      variant_name: variantName,
      size: size ? parseFloat(size) : undefined,
      unit: unit || undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 mt-4 border-t pt-4">
      <h4 className="text-sm font-medium">Add Variant</h4>
      {mutation.error && (
        <div className="text-red-500 text-sm">
          {mutation.error instanceof Error
            ? mutation.error.message
            : 'Failed to create variant'}
        </div>
      )}
      <div className="grid grid-cols-2 gap-2">
        <div className="space-y-1">
          <label htmlFor="variantName" className="text-xs font-medium">
            Name (e.g. 1L Bottle)
          </label>
          <Input
            id="variantName"
            type="text"
            value={variantName}
            onChange={(e) => setVariantName(e.target.value)}
            placeholder="e.g. 1L"
            required
            className="h-8"
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="size" className="text-xs font-medium">
            Size
          </label>
          <Input
            id="size"
            type="number"
            step="0.01"
            value={size}
            onChange={(e) => setSize(e.target.value)}
            placeholder="1.0"
            className="h-8"
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="unit" className="text-xs font-medium">
            Unit
          </label>
          <Input
            id="unit"
            type="text"
            value={unit}
            onChange={(e) => setUnit(e.target.value)}
            placeholder="L"
            className="h-8"
          />
        </div>
      </div>
      <Button type="submit" size="sm" disabled={mutation.isPending}>
        {mutation.isPending ? 'Adding...' : 'Add Variant'}
      </Button>
    </form>
  )
}
