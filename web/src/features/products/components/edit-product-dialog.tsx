import { useMutation, useQueryClient } from '@tanstack/react-query'
import { updateProduct } from '../api'
import type { Product } from '../types'
import { ProductForm } from './product-form'
import { Dialog } from '@/components/ui/dialog'
import { useInventoryStore } from '@/store/inventory'

interface Props {
  product: Product
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function EditProductDialog({ product, open, onOpenChange }: Props) {
  const queryClient = useQueryClient()
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const mutation = useMutation({
    mutationFn: (data: { name: string; brand?: string }) =>
      updateProduct(product.id, {
        ...data,
        canonical_product_id: product.canonical_product_id,
        category_id: product.category_id,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['products', activeInventoryId],
      })
      onOpenChange(false)
    },
  })

  return (
    <Dialog open={open} onOpenChange={onOpenChange} title="Edit Brand/Product">
      {mutation.error && (
        <div className="text-red-500 text-sm mb-4">
          {mutation.error instanceof Error
            ? mutation.error.message
            : 'Failed to update product'}
        </div>
      )}
      <ProductForm
        initialData={product}
        onSubmit={(data) => mutation.mutate(data)}
        isLoading={mutation.isPending}
        submitLabel="Save Changes"
      />
    </Dialog>
  )
}
