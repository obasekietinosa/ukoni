import { CanonicalProductList } from '../components/canonical-product-list'
import { CreateCanonicalProductForm } from '../components/create-canonical-product-form'

export function ProductsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Product Catalog</h1>
        <p className="text-gray-500">Manage your canonical products here.</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <div>
          <CreateCanonicalProductForm />
        </div>
        <div>
          <h2 className="mb-4 text-xl font-semibold">All Products</h2>
          <CanonicalProductList />
        </div>
      </div>
    </div>
  )
}
