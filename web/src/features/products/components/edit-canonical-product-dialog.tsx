import { useMutation, useQueryClient } from '@tanstack/react-query'
import { updateCanonicalProduct } from '../api'
import type { CanonicalProduct } from '../types'
import { CanonicalProductForm } from './canonical-product-form'
import { Dialog } from '@/components/ui/dialog'

interface Props {
  product: CanonicalProduct
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function EditCanonicalProductDialog({
  product,
  open,
  onOpenChange,
}: Props) {
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: { name: string; description?: string }) =>
      updateCanonicalProduct(product.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['canonical-products'],
      })
      queryClient.invalidateQueries({
        queryKey: ['canonical-product', product.id],
      })
      onOpenChange(false)
    },
  })

  return (
    <Dialog
      open={open}
      onOpenChange={onOpenChange}
      title="Edit Canonical Product"
    >
      {mutation.error && (
        <div className="text-red-500 text-sm mb-4">
          {mutation.error instanceof Error
            ? mutation.error.message
            : 'Failed to update product'}
        </div>
      )}
      <CanonicalProductForm
        initialData={product}
        onSubmit={(data) => mutation.mutate(data)}
        isLoading={mutation.isPending}
        submitLabel="Save Changes"
      />
    </Dialog>
  )
}
