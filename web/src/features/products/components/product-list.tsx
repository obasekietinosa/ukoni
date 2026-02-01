import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getProducts, deleteProduct } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { VariantList } from './variant-list'
import { CreateVariantForm } from './create-variant-form'
import { Button } from '@/components/ui/button'
import { EditProductDialog } from './edit-product-dialog'
import type { Product } from '../types'

interface Props {
  canonicalProductId: string
}

export function ProductList({ canonicalProductId }: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)

  const {
    data: products,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['products', activeInventoryId, canonicalProductId],
    queryFn: () => getProducts(activeInventoryId!, canonicalProductId),
    enabled: !!activeInventoryId && !!canonicalProductId,
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['products', activeInventoryId],
      })
    },
  })

  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this brand/product?')) {
      deleteMutation.mutate(id)
    }
  }

  if (isLoading) return <div>Loading brands...</div>
  if (error) return <div className="text-red-500">Failed to load brands</div>

  if (!products || products.length === 0) {
    return (
      <div className="text-gray-500">No brands found for this product.</div>
    )
  }

  return (
    <div className="space-y-4">
      {products.map((product) => (
        <div key={product.id} className="rounded-lg border p-4 bg-white">
          <div className="flex items-center justify-between mb-4">
            <div>
              <div className="font-medium text-lg">{product.name}</div>
              {product.brand && (
                <div className="text-sm text-gray-500">
                  Brand: {product.brand}
                </div>
              )}
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setEditingProduct(product)}
              >
                Edit
              </Button>
              <Button
                variant="destructive"
                size="sm"
                onClick={() => handleDelete(product.id)}
              >
                Delete
              </Button>
            </div>
          </div>

          <div className="pl-4 border-l-2 border-gray-100">
            <div className="text-sm text-gray-500 mb-2">Variants:</div>
            <VariantList productId={product.id} />
            <CreateVariantForm productId={product.id} />
          </div>
        </div>
      ))}

      {editingProduct && (
        <EditProductDialog
          product={editingProduct}
          open={!!editingProduct}
          onOpenChange={(open) => !open && setEditingProduct(null)}
        />
      )}
    </div>
  )
}
