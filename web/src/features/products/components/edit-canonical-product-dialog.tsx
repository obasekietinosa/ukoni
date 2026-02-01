import { useState, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { updateCanonicalProduct } from '../api'
import type { CanonicalProduct } from '../types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

interface Props {
  product: CanonicalProduct
  open: boolean
  onClose: () => void
}

export function EditCanonicalProductDialog({ product, open, onClose }: Props) {
  const [name, setName] = useState(product.name)
  const [description, setDescription] = useState(product.description || '')
  const queryClient = useQueryClient()

  useEffect(() => {
    if (open) {
      setName(product.name)
      setDescription(product.description || '')
    }
  }, [open, product])

  const mutation = useMutation({
    mutationFn: (data: { name: string; description?: string }) =>
      updateCanonicalProduct(product.id, {
        ...data,
        category_id: product.category_id,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['canonical-product', product.id],
      })
      queryClient.invalidateQueries({
        queryKey: ['canonical-products', product.inventory_id],
      })
      onClose()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name) return
    mutation.mutate({ name, description })
  }

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
        <h3 className="text-lg font-medium mb-4">Edit Product</h3>
        {mutation.error && (
          <div className="text-red-500 text-sm mb-4">
            {mutation.error instanceof Error
              ? mutation.error.message
              : 'Failed to update product'}
          </div>
        )}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="edit-name" className="text-sm font-medium">
              Product Name
            </label>
            <Input
              id="edit-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="edit-description" className="text-sm font-medium">
              Description (Optional)
            </label>
            <Input
              id="edit-description"
              type="text"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
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
