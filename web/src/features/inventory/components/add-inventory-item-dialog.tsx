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
  createVariant,
} from '@/features/products/api'
import { createTransaction } from '@/features/transactions/api'
import type {
  CanonicalProduct,
  Product,
  ProductVariant,
} from '@/features/products/types'
import { Loader2, ArrowLeft, Search, Plus } from 'lucide-react'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess?: () => void
}

export function AddInventoryItemDialog({ open, onOpenChange, onSuccess }: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  // Steps: 'search-canonical' | 'select-product' | 'select-variant' | 'create-variant' | 'input-quantity'
  const [step, setStep] = useState<
    'search-canonical' | 'select-product' | 'select-variant' | 'create-variant' | 'input-quantity'
  >('search-canonical')

  const [searchQuery, setSearchQuery] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')

  const [selectedCanonical, setSelectedCanonical] = useState<CanonicalProduct | null>(null)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(null)
  const [quantity, setQuantity] = useState('1')

  // New variant state
  const [newVariantName, setNewVariantName] = useState('')
  const [newVariantSize, setNewVariantSize] = useState('')
  const [newVariantUnit, setNewVariantUnit] = useState('')

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
        setNewVariantName('')
        setNewVariantSize('')
        setNewVariantUnit('')
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
    enabled: !!activeInventoryId && (debouncedSearch.length > 0 || step === 'search-canonical'),
  })

  const { data: products, isLoading: isLoadingProducts } = useQuery({
    queryKey: ['products', activeInventoryId, selectedCanonical?.id],
    queryFn: () => getProducts(activeInventoryId!, selectedCanonical!.id),
    enabled: !!activeInventoryId && !!selectedCanonical,
  })

  const { data: variants, isLoading: isLoadingVariants } = useQuery({
    queryKey: ['variants', selectedProduct?.id],
    queryFn: () => getVariants(selectedProduct!.id),
    enabled: !!selectedProduct,
  })

  // Mutations
  const createVariantMutation = useMutation({
    mutationFn: (data: {
      variant_name: string
      unit?: string
      size?: number
    }) => createVariant(selectedProduct!.id, data),
    onSuccess: (newVariant) => {
      queryClient.invalidateQueries({ queryKey: ['variants', selectedProduct!.id] })
      setSelectedVariant(newVariant)
      setStep('input-quantity')
      setNewVariantName('')
      setNewVariantSize('')
      setNewVariantUnit('')
    },
  })

  const addTransactionMutation = useMutation({
    mutationFn: () =>
      createTransaction(activeInventoryId!, {
        transaction_date: new Date().toISOString(),
        items: [
          {
            product_variant_id: selectedVariant!.id,
            quantity: parseFloat(quantity),
          },
        ],
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['inventory-products', activeInventoryId],
      })
      onSuccess?.()
      onOpenChange(false)
    },
  })

  const handleBack = () => {
    if (step === 'select-product') {
      setStep('search-canonical')
      setSelectedCanonical(null)
    } else if (step === 'select-variant') {
      setStep('select-product')
      setSelectedProduct(null)
    } else if (step === 'create-variant') {
      setStep('select-variant')
    } else if (step === 'input-quantity') {
      setStep('select-variant')
      setSelectedVariant(null)
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    addTransactionMutation.mutate()
  }

  const handleCreateVariant = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newVariantName) return
    createVariantMutation.mutate({
      variant_name: newVariantName,
      size: newVariantSize ? parseFloat(newVariantSize) : undefined,
      unit: newVariantUnit || undefined,
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <div className="flex items-center gap-2">
            {step !== 'search-canonical' && (
              <Button variant="ghost" size="icon" className="h-6 w-6" onClick={handleBack}>
                <ArrowLeft className="h-4 w-4" />
              </Button>
            )}
            <DialogTitle>
              {step === 'search-canonical' && 'Select Item to Add'}
              {step === 'select-product' && `Select ${selectedCanonical?.name} Brand`}
              {step === 'select-variant' && `Select ${selectedProduct?.name} Variant`}
              {step === 'create-variant' && `Create Variant for ${selectedProduct?.name}`}
              {step === 'input-quantity' && 'Enter Quantity'}
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
                  placeholder="Search products (e.g. Milk, Bread)..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-9"
                />
              </div>
              <div className="max-h-[300px] overflow-y-auto space-y-2">
                {isLoadingCanonical ? (
                  <div className="flex justify-center p-4">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-500" />
                  </div>
                ) : canonicalProducts?.length === 0 ? (
                  <div className="text-center text-gray-500 p-4">
                    No products found.
                  </div>
                ) : (
                  canonicalProducts?.map((cp) => (
                    <div
                      key={cp.id}
                      className="p-3 border rounded-md hover:bg-gray-50 cursor-pointer transition-colors"
                      onClick={() => {
                        setSelectedCanonical(cp)
                        setStep('select-product')
                        setSearchQuery('') // Clear search for next time or irrelevant now
                      }}
                    >
                      <div className="font-medium">{cp.name}</div>
                      {cp.description && (
                        <div className="text-sm text-gray-500">{cp.description}</div>
                      )}
                    </div>
                  ))
                )}
              </div>
            </div>
          )}

          {/* Step 2: Select Product */}
          {step === 'select-product' && (
            <div className="space-y-2">
               {isLoadingProducts ? (
                  <div className="flex justify-center p-4">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-500" />
                  </div>
                ) : products?.length === 0 ? (
                  <div className="text-center text-gray-500 p-4">
                    No brands found for {selectedCanonical?.name}.
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
            </div>
          )}

          {/* Step 3: Select Variant */}
          {step === 'select-variant' && (
            <div className="space-y-2">
               {isLoadingVariants ? (
                  <div className="flex justify-center p-4">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-500" />
                  </div>
                ) : (
                  <>
                    {variants?.map((v) => (
                      <div
                        key={v.id}
                        className="p-3 border rounded-md hover:bg-gray-50 cursor-pointer transition-colors"
                        onClick={() => {
                          setSelectedVariant(v)
                          setStep('input-quantity')
                        }}
                      >
                        <div className="font-medium">{v.variant_name}</div>
                        <div className="text-sm text-gray-500">
                          {v.size} {v.unit}
                        </div>
                      </div>
                    ))}
                    <Button
                      variant="outline"
                      className="w-full justify-start mt-2"
                      onClick={() => setStep('create-variant')}
                    >
                      <Plus className="mr-2 h-4 w-4" />
                      Create New Variant
                    </Button>
                  </>
                )}
            </div>
          )}

          {/* Step 3.5: Create Variant */}
          {step === 'create-variant' && (
            <form onSubmit={handleCreateVariant} className="space-y-4">
              {createVariantMutation.error && (
                <div className="text-red-500 text-sm">
                  {createVariantMutation.error instanceof Error
                    ? createVariantMutation.error.message
                    : 'Failed to create variant'}
                </div>
              )}
              <div className="space-y-2">
                <label className="text-sm font-medium">Variant Name</label>
                <Input
                  value={newVariantName}
                  onChange={(e) => setNewVariantName(e.target.value)}
                  placeholder="e.g. 1L Bottle"
                  required
                  autoFocus
                />
              </div>
              <div className="flex gap-4">
                <div className="space-y-2 flex-1">
                  <label className="text-sm font-medium">Size</label>
                  <Input
                    type="number"
                    step="0.01"
                    value={newVariantSize}
                    onChange={(e) => setNewVariantSize(e.target.value)}
                    placeholder="1"
                  />
                </div>
                <div className="space-y-2 flex-1">
                  <label className="text-sm font-medium">Unit</label>
                  <Input
                    value={newVariantUnit}
                    onChange={(e) => setNewVariantUnit(e.target.value)}
                    placeholder="L"
                  />
                </div>
              </div>
              <DialogFooter>
                 <Button type="button" variant="outline" onClick={handleBack}>
                  Back
                </Button>
                <Button type="submit" disabled={createVariantMutation.isPending}>
                  {createVariantMutation.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Create & Select
                </Button>
              </DialogFooter>
            </form>
          )}

           {/* Step 4: Input Quantity */}
           {step === 'input-quantity' && (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <div className="text-sm text-gray-500">
                  Adding: <span className="font-medium text-gray-900">
                    {selectedProduct?.name} - {selectedVariant?.variant_name} ({selectedVariant?.size} {selectedVariant?.unit})
                  </span>
                </div>
                <label htmlFor="quantity" className="text-sm font-medium">
                  Quantity
                </label>
                <Input
                  id="quantity"
                  type="number"
                  step="0.01"
                  min="0.01"
                  value={quantity}
                  onChange={(e) => setQuantity(e.target.value)}
                  placeholder="1"
                  required
                  autoFocus
                />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={handleBack}>
                  Back
                </Button>
                <Button type="submit" disabled={addTransactionMutation.isPending}>
                  {addTransactionMutation.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Add to Inventory
                </Button>
              </DialogFooter>
            </form>
          )}

        </div>
      </DialogContent>
    </Dialog>
  )
}
