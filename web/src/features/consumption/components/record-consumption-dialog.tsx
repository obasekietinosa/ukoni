import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useInventoryStore } from '@/store/inventory'
import {
  getCanonicalProducts,
  getProducts,
  getVariants,
} from '@/features/products/api'
import { createConsumptionEvent } from '../api'
import type {
  CanonicalProduct,
  Product,
  ProductVariant,
} from '@/features/products/types'
import { Loader2, ArrowLeft, Search, Calendar } from 'lucide-react'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess?: () => void
}

export function RecordConsumptionDialog({
  open,
  onOpenChange,
  onSuccess,
}: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const [step, setStep] = useState<
    'search-canonical' | 'details' | 'select-product' | 'select-variant'
  >('search-canonical')

  const [searchQuery, setSearchQuery] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')

  const [selectedCanonical, setSelectedCanonical] =
    useState<CanonicalProduct | null>(null)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(
    null
  )

  const [quantity, setQuantity] = useState('1')
  const [unit, setUnit] = useState('')
  const [date, setDate] = useState(new Date().toISOString().slice(0, 16)) // datetime-local format

  // Debounce search
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchQuery)
    }, 500)
    return () => clearTimeout(timer)
  }, [searchQuery])

  // Reset state on close
  useEffect(() => {
    if (!open) {
      setTimeout(() => {
        setStep('search-canonical')
        setSearchQuery('')
        setSelectedCanonical(null)
        setSelectedProduct(null)
        setSelectedVariant(null)
        setQuantity('1')
        setUnit('')
        setDate(new Date().toISOString().slice(0, 16))
      }, 300)
    }
  }, [open])

  // Queries
  const { data: canonicalProducts, isLoading: isLoadingCanonical } = useQuery({
    queryKey: ['canonical-products', activeInventoryId, debouncedSearch],
    queryFn: () =>
      getCanonicalProducts(activeInventoryId!, {
        search: debouncedSearch,
        limit: 10,
      }),
    enabled: !!activeInventoryId && open && step === 'search-canonical',
  })

  const { data: products, isLoading: isLoadingProducts } = useQuery({
    queryKey: ['products', activeInventoryId, selectedCanonical?.id],
    queryFn: () => getProducts(activeInventoryId!, selectedCanonical!.id),
    enabled:
      !!activeInventoryId && !!selectedCanonical && step === 'select-product',
  })

  const { data: variants, isLoading: isLoadingVariants } = useQuery({
    queryKey: ['variants', selectedProduct?.id],
    queryFn: () => getVariants(selectedProduct!.id),
    enabled: !!selectedProduct && step === 'select-variant',
  })

  // Mutation
  const createConsumptionMutation = useMutation({
    mutationFn: () =>
      createConsumptionEvent(activeInventoryId!, {
        canonical_product_id: selectedCanonical?.id,
        product_variant_id: selectedVariant?.id,
        quantity: quantity ? parseFloat(quantity) : undefined,
        unit: unit || undefined,
        consumed_at: new Date(date).toISOString(),
        source: 'manual',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['consumption-events', activeInventoryId],
      })
      onSuccess?.()
      onOpenChange(false)
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createConsumptionMutation.mutate()
  }

  const handleBack = () => {
    if (step === 'details') {
      setStep('search-canonical')
      setSelectedCanonical(null)
      setSelectedVariant(null)
      setUnit('')
    } else if (step === 'select-product') {
      setStep('details')
    } else if (step === 'select-variant') {
      setStep('select-product')
      setSelectedProduct(null)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <div className="flex items-center gap-2">
            {step !== 'search-canonical' && (
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={handleBack}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
            )}
            <DialogTitle>
              {step === 'search-canonical' && 'Record Consumption'}
              {step === 'details' && 'Consumption Details'}
              {step === 'select-product' &&
                `Select Brand for ${selectedCanonical?.name}`}
              {step === 'select-variant' &&
                `Select Variant for ${selectedProduct?.name}`}
            </DialogTitle>
          </div>
        </DialogHeader>

        <div className="py-4">
          {/* Step 1: Search Canonical */}
          {step === 'search-canonical' && (
            <div className="space-y-4">
              <div className="relative">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
                <Input
                  placeholder="Search products (e.g. Milk)..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-9"
                  autoFocus
                />
              </div>
              <div className="max-h-[300px] overflow-y-auto space-y-2">
                {isLoadingCanonical ? (
                  <div className="flex justify-center p-4">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-500" />
                  </div>
                ) : (
                  canonicalProducts?.map((cp) => (
                    <div
                      key={cp.id}
                      className="p-3 border rounded-md hover:bg-gray-50 cursor-pointer transition-colors"
                      onClick={() => {
                        setSelectedCanonical(cp)
                        setStep('details')
                        setSearchQuery('')
                      }}
                    >
                      <div className="font-medium">{cp.name}</div>
                    </div>
                  ))
                )}
                {!isLoadingCanonical && canonicalProducts?.length === 0 && (
                  <div className="text-center text-sm text-gray-500 py-4">
                    No products found.
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Step 2: Details */}
          {step === 'details' && selectedCanonical && (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="p-3 bg-gray-50 rounded-md border">
                <div className="font-medium">{selectedCanonical.name}</div>
                {selectedVariant && (
                  <div className="text-sm text-gray-600 mt-1">
                    Variant: {selectedVariant.variant_name} (
                    {selectedVariant.size} {selectedVariant.unit})
                  </div>
                )}
                <Button
                  type="button"
                  variant="link"
                  className="p-0 h-auto text-sm mt-1"
                  onClick={() => setStep('select-product')}
                >
                  {selectedVariant
                    ? 'Change Variant'
                    : 'Select Specific Variant (Optional)'}
                </Button>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label htmlFor="quantity" className="text-sm font-medium">
                    Quantity
                  </label>
                  <Input
                    id="quantity"
                    type="number"
                    step="0.01"
                    value={quantity}
                    onChange={(e) => setQuantity(e.target.value)}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="unit" className="text-sm font-medium">
                    Unit
                  </label>
                  <Input
                    id="unit"
                    value={unit}
                    onChange={(e) => setUnit(e.target.value)}
                    placeholder={selectedVariant?.unit || 'e.g. kg, pcs'}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label htmlFor="date" className="text-sm font-medium">
                  Date & Time
                </label>
                <div className="relative">
                  <Calendar className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
                  <Input
                    id="date"
                    type="datetime-local"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                    className="pl-9"
                    required
                  />
                </div>
              </div>

              <DialogFooter>
                <Button
                  type="submit"
                  disabled={createConsumptionMutation.isPending}
                >
                  {createConsumptionMutation.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Record Consumption
                </Button>
              </DialogFooter>
            </form>
          )}

          {/* Step 3: Select Product */}
          {step === 'select-product' && (
            <div className="space-y-2">
              {isLoadingProducts ? (
                <div className="flex justify-center p-4">
                  <Loader2 className="h-6 w-6 animate-spin text-gray-500" />
                </div>
              ) : (
                products?.map((p) => (
                  <div
                    key={p.id}
                    className="p-3 border rounded-md hover:bg-gray-50 cursor-pointer transition-colors"
                    onClick={() => {
                      setSelectedProduct(p)
                      setStep('select-variant')
                    }}
                  >
                    <div className="font-medium">{p.name}</div>
                    {p.brand && (
                      <div className="text-sm text-gray-500">{p.brand}</div>
                    )}
                  </div>
                ))
              )}
              {!isLoadingProducts && products?.length === 0 && (
                <div className="text-center text-sm text-gray-500 py-4">
                  No brands found for this product.
                </div>
              )}
            </div>
          )}

          {/* Step 4: Select Variant */}
          {step === 'select-variant' && (
            <div className="space-y-2">
              {isLoadingVariants ? (
                <div className="flex justify-center p-4">
                  <Loader2 className="h-6 w-6 animate-spin text-gray-500" />
                </div>
              ) : (
                variants?.map((v) => (
                  <div
                    key={v.id}
                    className="p-3 border rounded-md hover:bg-gray-50 cursor-pointer transition-colors"
                    onClick={() => {
                      setSelectedVariant(v)
                      setUnit(v.unit || '')
                      setStep('details')
                    }}
                  >
                    <div className="font-medium">{v.variant_name}</div>
                    <div className="text-sm text-gray-500">
                      {v.size} {v.unit}
                    </div>
                  </div>
                ))
              )}
              {!isLoadingVariants && variants?.length === 0 && (
                <div className="text-center text-sm text-gray-500 py-4">
                  No variants found.
                </div>
              )}
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
