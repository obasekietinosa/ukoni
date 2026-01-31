import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getCanonicalProduct } from '../api'
import { ProductList } from '../components/product-list'
import { CreateProductForm } from '../components/create-product-form'
import { VariantList } from '../components/variant-list'
import { CreateVariantForm } from '../components/create-variant-form'
import { Button } from '@/components/ui/button'

export function ProductDetailsPage() {
  const { canonicalId } = useParams<{ canonicalId: string }>()

  const {
    data: product,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['canonical-product', canonicalId],
    queryFn: () => getCanonicalProduct(canonicalId!),
    enabled: !!canonicalId,
  })

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>Error loading product</div>
  if (!product) return <div>Product not found</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/products">
          <Button variant="outline" size="sm">
            Back
          </Button>
        </Link>
        <div>
          <h1 className="text-2xl font-bold">{product.name}</h1>
          {product.description && (
            <p className="text-gray-500">{product.description}</p>
          )}
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <div>
          <CreateProductForm
            canonicalProductId={product.id}
            defaultName={product.name}
          />
        </div>
        <div>
          <h2 className="mb-4 text-xl font-semibold">
            Available Brands & Variants
          </h2>
          <ProductList canonicalProductId={product.id}>
            {(product) => (
              <>
                <VariantList productId={product.id} />
                <CreateVariantForm productId={product.id} />
              </>
            )}
          </ProductList>
        </div>
      </div>
    </div>
  )
}
