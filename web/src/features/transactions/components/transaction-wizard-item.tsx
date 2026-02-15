import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Input } from '@/components/ui/input'
import type { ShoppingListItem } from '@/features/shopping-lists/types'
import {
  getProducts,
  getVariants,
  createProduct,
  createVariant,
} from '@/features/products/api'
import { useInventoryStore } from '@/store/inventory'
import { Loader2 } from 'lucide-react'

export interface WizardItemState {
  shoppingListItemId: string
  selectedVariantId: string | null
  quantity: number
  price: number | null
  included: boolean
}

interface Props {
  item: ShoppingListItem
  onUpdate: (state: WizardItemState) => void
  initialIncluded?: boolean
}

const CREATE_GENERIC_VALUE = '__CREATE_GENERIC__'
const CREATE_DEFAULT_VARIANT_VALUE = '__CREATE_DEFAULT_VARIANT__'

export function TransactionWizardItem({
  item,
  onUpdate,
  initialIncluded,
}: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  // State
  const [included, setIncluded] = useState(initialIncluded ?? true)
  const [quantity, setQuantity] = useState(item.quantity || 1)
  const [price, setPrice] = useState<string>('') // User input string, parsed to float
  const [selectedBrandId, setSelectedBrandId] = useState<string | null>(
    item.product?.id || null
  )
  const [selectedVariantId, setSelectedVariantId] = useState<string | null>(
    item.target_type === 'product_variant' ? item.target_id : null
  )

  // Notify parent of changes
  useEffect(() => {
    onUpdate({
      shoppingListItemId: item.id,
      selectedVariantId: included ? selectedVariantId : null,
      quantity,
      price: price ? parseFloat(price) : null,
      included,
    })
  }, [included, quantity, price, selectedVariantId, item.id, onUpdate])

  // Queries for resolving canonical -> variant
  const isCanonical = item.target_type === 'canonical_product'

  const { data: products, isLoading: isLoadingProducts } = useQuery({
    queryKey: ['products', activeInventoryId, item.target_id],
    queryFn: () => getProducts(activeInventoryId!, item.target_id),
    enabled: !!activeInventoryId && isCanonical && included,
  })

  const { data: variants, isLoading: isLoadingVariants } = useQuery({
    queryKey: ['variants', selectedBrandId],
    queryFn: () => getVariants(selectedBrandId!),
    enabled: !!selectedBrandId && included,
  })

  // Mutations for creating generic product/variant
  const createProductMutation = useMutation({
    mutationFn: async () => {
      if (!item.canonical_product) throw new Error('No canonical product')
      const newProduct = await createProduct(activeInventoryId!, {
        name: item.canonical_product.name,
        canonical_product_id: item.target_id,
        // No brand implies generic
      })

      // Immediately create a variant for this product
      const newVariant = await createVariant(newProduct.id, {
        variant_name: 'Regular',
      })

      return { product: newProduct, variant: newVariant }
    },
    onSuccess: ({ product, variant }) => {
      queryClient.invalidateQueries({
        queryKey: ['products', activeInventoryId, item.target_id],
      })
      queryClient.setQueryData(['variants', product.id], [variant])
      setSelectedBrandId(product.id)
      setSelectedVariantId(variant.id)
    },
  })

  const createVariantMutation = useMutation({
    mutationFn: async () => {
      if (!selectedBrandId) throw new Error('No brand selected')
      return createVariant(selectedBrandId, {
        variant_name: 'Regular',
        // defaults
      })
    },
    onSuccess: (newVariant) => {
      queryClient.invalidateQueries({
        queryKey: ['variants', selectedBrandId],
      })
      setSelectedVariantId(newVariant.id)
    },
  })

  const isCreating =
    createProductMutation.isPending || createVariantMutation.isPending

  // Derived state for dropdowns
  const hasGenericProduct = products?.some((p) => !p.brand)
  const hasDefaultVariant = variants?.some((v) => v.variant_name === 'Regular')

  // Validation UI
  const showBrandError = included && isCanonical && !selectedBrandId
  const showVariantError =
    included && isCanonical && selectedBrandId && !selectedVariantId

  // Handlers
  const handleBrandChange = (val: string) => {
    if (val === CREATE_GENERIC_VALUE) {
      createProductMutation.mutate()
    } else {
      setSelectedBrandId(val)
      setSelectedVariantId(null)
    }
  }

  const handleVariantChange = (val: string) => {
    if (val === CREATE_DEFAULT_VARIANT_VALUE) {
      createVariantMutation.mutate()
    } else {
      setSelectedVariantId(val)
    }
  }

  const selectBaseClassName =
    'flex h-8 w-full rounded-md border px-3 py-1 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-950 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50'

  return (
    <div
      className={`rounded-lg border p-4 transition-colors ${
        included ? 'bg-white border-gray-200' : 'bg-gray-50 border-gray-100'
      }`}
    >
      <div className="flex items-start gap-4">
        <input
          type="checkbox"
          checked={included}
          onChange={(e) => setIncluded(e.target.checked)}
          className="mt-1 h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
        />

        <div className="flex-1 space-y-3">
          {/* Header: Item Name */}
          <div className="flex items-center justify-between">
            <div>
              <div
                className={`font-medium ${!included ? 'text-gray-400' : ''}`}
              >
                {isCanonical
                  ? item.canonical_product?.name
                  : `${item.product?.brand} ${item.product_variant?.variant_name}`}
              </div>
              {item.notes && (
                <div className="text-sm text-gray-500 italic">
                  "{item.notes}"
                </div>
              )}
            </div>
            {isCreating && <Loader2 className="h-4 w-4 animate-spin" />}
          </div>

          {included && (
            <div className="space-y-3">
              {/* Resolution for Canonical Items */}
              {isCanonical && (
                <div className="grid grid-cols-2 gap-2">
                  {/* Brand Select */}
                  <div>
                    <label className="text-xs text-gray-500 font-medium flex justify-between">
                      Brand
                      {showBrandError && (
                        <span className="text-red-500">* Required</span>
                      )}
                    </label>
                    <select
                      value={selectedBrandId || ''}
                      onChange={(e) => handleBrandChange(e.target.value)}
                      disabled={isLoadingProducts || isCreating}
                      className={`${selectBaseClassName} ${
                        showBrandError
                          ? 'border-red-300 bg-red-50'
                          : 'border-gray-200 bg-white'
                      }`}
                    >
                      <option value="" disabled>
                        {isLoadingProducts ? 'Loading...' : 'Select Brand'}
                      </option>
                      {!hasGenericProduct && (
                        <option value={CREATE_GENERIC_VALUE}>
                          Generic (No Brand)
                        </option>
                      )}
                      {products?.map((p) => (
                        <option key={p.id} value={p.id}>
                          {p.brand
                            ? `${p.brand} - ${p.name}`
                            : `${p.name} (Generic)`}
                        </option>
                      ))}
                    </select>
                  </div>

                  {/* Variant Select */}
                  <div>
                    <label className="text-xs text-gray-500 font-medium flex justify-between">
                      Variant
                      {showVariantError && (
                        <span className="text-red-500">* Required</span>
                      )}
                    </label>
                    <select
                      value={selectedVariantId || ''}
                      onChange={(e) => handleVariantChange(e.target.value)}
                      disabled={
                        !selectedBrandId || isLoadingVariants || isCreating
                      }
                      className={`${selectBaseClassName} ${
                        showVariantError
                          ? 'border-red-300 bg-red-50'
                          : 'border-gray-200 bg-white'
                      }`}
                    >
                      <option value="" disabled>
                        {isLoadingVariants
                          ? 'Loading...'
                          : !selectedBrandId
                            ? 'Select Brand First'
                            : 'Select Variant'}
                      </option>
                      {!hasDefaultVariant && selectedBrandId && (
                        <option value={CREATE_DEFAULT_VARIANT_VALUE}>
                          Regular / Default
                        </option>
                      )}
                      {variants?.map((v) => (
                        <option key={v.id} value={v.id}>
                          {v.variant_name} ({v.size} {v.unit})
                        </option>
                      ))}
                    </select>
                  </div>
                </div>
              )}

              {/* Quantity and Price */}
              <div className="flex gap-2">
                <div className="flex-1">
                  <label className="text-xs text-gray-500 font-medium">
                    Qty
                  </label>
                  <Input
                    type="number"
                    step="0.01"
                    min="0"
                    value={quantity}
                    onChange={(e) => setQuantity(parseFloat(e.target.value))}
                    className="h-8"
                  />
                </div>
                <div className="flex-1">
                  <label className="text-xs text-gray-500 font-medium">
                    Price (Optional)
                  </label>
                  <Input
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="0.00"
                    value={price}
                    onChange={(e) => setPrice(e.target.value)}
                    className="h-8"
                  />
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
