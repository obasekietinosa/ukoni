import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createCanonicalProduct } from '../api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useInventoryStore } from '@/store/inventory'

interface Props {
  onSuccess?: () => void
}

export function CreateCanonicalProductForm({ onSuccess }: Props) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const activeInventoryId = useInventoryStore((state) => state.activeInventoryId)
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: { name: string; description?: string }) =>
      createCanonicalProduct(activeInventoryId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['canonical-products', activeInventoryId] })
      setName('')
      setDescription('')
      onSuccess?.()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name) return
    mutation.mutate({ name, description })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 rounded-lg border p-4 shadow-sm bg-white">
      <h3 className="text-lg font-medium">Add New Product</h3>
      {mutation.error && (
        <div className="text-red-500 text-sm">
          {mutation.error instanceof Error ? mutation.error.message : 'Failed to create product'}
        </div>
      )}
      <div className="space-y-2">
        <label htmlFor="name" className="text-sm font-medium">
          Product Name
        </label>
        <Input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g. Milk, Bread"
          required
        />
      </div>
       <div className="space-y-2">
        <label htmlFor="description" className="text-sm font-medium">
          Description (Optional)
        </label>
        <Input
          id="description"
          type="text"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="e.g. Dairy"
        />
      </div>
      <Button type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? 'Adding...' : 'Add Product'}
      </Button>
    </form>
  )
}
