import { useState, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { updateProduct } from '../api'
import type { Product } from '../types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

interface Props {
  product: Product
  open: boolean
  onClose: () => void
}

export function EditProductDialog({ product, open, onClose }: Props) {
  const [name, setName] = useState(product.name)
  const [brand, setBrand] = useState(product.brand || '')
  const [description, setDescription] = useState(product.description || '')
  const queryClient = useQueryClient()

  useEffect(() => {
    if (open) {
      setName(product.name)
      setBrand(product.brand || '')
      setDescription(product.description || '')
    }
  }, [open, product])

  const mutation = useMutation({
    mutationFn: (data: { name: string; brand?: string; description?: string }) =>
      updateProduct(product.id, {
        ...data,
        category_id: product.category_id,
        canonical_product_id: product.canonical_product_id,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['products', product.inventory_id],
      })
      onClose()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name) return
    mutation.mutate({ name, brand, description })
  }

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
        <h3 className="text-lg font-medium mb-4">Edit Brand/Product</h3>
        {mutation.error && (
          <div className="text-red-500 text-sm mb-4">
            {mutation.error instanceof Error
              ? mutation.error.message
              : 'Failed to update product'}
          </div>
        )}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="edit-product-brand" className="text-sm font-medium">
              Brand
            </label>
            <Input
              id="edit-product-brand"
              type="text"
              value={brand}
              onChange={(e) => setBrand(e.target.value)}
              placeholder="e.g. Tesco"
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="edit-product-name" className="text-sm font-medium">
              Product Name
            </label>
            <Input
              id="edit-product-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Whole Milk"
              required
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="edit-product-description" className="text-sm font-medium">
              Description (Optional)
            </label>
            <Input
              id="edit-product-description"
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
