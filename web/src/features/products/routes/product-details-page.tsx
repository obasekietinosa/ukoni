import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getCanonicalProduct, deleteCanonicalProduct } from '../api'
import { Button } from '@/components/ui/button'
import { ProductList } from '../components/product-list'
import { CreateProductForm } from '../components/create-product-form'
import { EditCanonicalProductDialog } from '../components/edit-canonical-product-dialog'

export function ProductDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [isEditOpen, setIsEditOpen] = useState(false)

  const {
    data: product,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['canonical-product', id],
    queryFn: () => getCanonicalProduct(id!),
    enabled: !!id,
  })

  const deleteMutation = useMutation({
    mutationFn: () => deleteCanonicalProduct(id!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['canonical-products'] })
      navigate('/products')
    },
  })

  const handleDelete = () => {
    if (window.confirm('Are you sure you want to delete this product?')) {
      deleteMutation.mutate()
    }
  }

  if (isLoading) return <div>Loading...</div>
  if (error) return <div className="text-red-500">Failed to load product</div>
  if (!product) return <div className="text-gray-500">Product not found</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link to="/products">
            <Button variant="outline" size="sm">
              Back
            </Button>
          </Link>
          <h1 className="text-2xl font-bold">{product.name}</h1>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={() => setIsEditOpen(true)}>
            Edit
          </Button>
          <Button
            variant="destructive"
            size="sm"
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
          </Button>
        </div>
      </div>

      {product.description && (
        <p className="text-gray-600">{product.description}</p>
      )}

      <EditCanonicalProductDialog
        product={product}
        open={isEditOpen}
        onOpenChange={setIsEditOpen}
      />

      <div className="border-t pt-4">
        <h2 className="text-xl font-semibold mb-4">Brands & Variants</h2>
        <div className="grid gap-6 md:grid-cols-3">
          <div className="md:col-span-2">
            <ProductList canonicalProductId={id!} />
          </div>
          <div>
            <CreateProductForm canonicalProductId={id!} />
          </div>
        </div>
      </div>
    </div>
  )
}
