import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createCanonicalProduct } from '../api'
import { CanonicalProductForm } from './canonical-product-form'
import { useInventoryStore } from '@/store/inventory'

interface Props {
  onSuccess?: () => void
}

export function CreateCanonicalProductForm({ onSuccess }: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: { name: string; description?: string }) =>
      createCanonicalProduct(activeInventoryId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['canonical-products', activeInventoryId],
      })
      onSuccess?.()
    },
  })

  return (
    <div className="rounded-lg border p-4 shadow-sm bg-white">
      <h3 className="text-lg font-medium mb-4">Add New Product</h3>
      {mutation.error && (
        <div className="text-red-500 text-sm mb-4">
          {mutation.error instanceof Error
            ? mutation.error.message
            : 'Failed to create product'}
        </div>
      )}
      <CanonicalProductForm
        onSubmit={(data) => mutation.mutate(data)}
        isLoading={mutation.isPending}
        submitLabel="Add Product"
      />
    </div>
  )
}
