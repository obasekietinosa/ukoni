import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createProduct } from '../api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useInventoryStore } from '@/store/inventory'

interface Props {
  canonicalProductId: string
  onSuccess?: () => void
}

export function CreateProductForm({ canonicalProductId, onSuccess }: Props) {
  const [name, setName] = useState('')
  const [brand, setBrand] = useState('')
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: {
      name: string
      brand?: string
      canonical_product_id: string
    }) => createProduct(activeInventoryId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['products', activeInventoryId],
      })
      setName('')
      setBrand('')
      onSuccess?.()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name) return
    mutation.mutate({
      name,
      brand: brand || undefined,
      canonical_product_id: canonicalProductId,
    })
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="space-y-4 rounded-lg border p-4 shadow-sm bg-white"
    >
      <h3 className="text-lg font-medium">Add Brand/Product</h3>
      {mutation.error && (
        <div className="text-red-500 text-sm">
          {mutation.error instanceof Error
            ? mutation.error.message
            : 'Failed to create product'}
        </div>
      )}
      <div className="space-y-2">
        <label htmlFor="product-name" className="text-sm font-medium">
          Product Name
        </label>
        <Input
          id="product-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g. Tesco Whole Milk"
          required
        />
      </div>
      <div className="space-y-2">
        <label htmlFor="brand" className="text-sm font-medium">
          Brand (Optional)
        </label>
        <Input
          id="brand"
          type="text"
          value={brand}
          onChange={(e) => setBrand(e.target.value)}
          placeholder="e.g. Tesco"
        />
      </div>
      <Button type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? 'Adding...' : 'Add Brand/Product'}
      </Button>
    </form>
  )
}
