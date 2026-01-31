import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getCanonicalProduct } from '../api'
import { Button } from '@/components/ui/button'
import { ProductList } from '../components/product-list'
import { CreateProductForm } from '../components/create-product-form'

export function ProductDetailsPage() {
  const { id } = useParams<{ id: string }>()

  const {
    data: product,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['canonical-product', id],
    queryFn: () => getCanonicalProduct(id!),
    enabled: !!id,
  })

  if (isLoading) return <div>Loading...</div>
  if (error) return <div className="text-red-500">Failed to load product</div>
  if (!product) return <div className="text-gray-500">Product not found</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/products">
          <Button variant="outline" size="sm">
            Back
          </Button>
        </Link>
        <h1 className="text-2xl font-bold">{product.name}</h1>
      </div>

      {product.description && (
        <p className="text-gray-600">{product.description}</p>
      )}

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
