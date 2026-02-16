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
import { Textarea } from '@/components/ui/textarea'
import { useInventoryStore } from '@/store/inventory'
import {
  getCanonicalProducts,
  getProducts,
  getVariants,
  createCanonicalProduct,
  createProduct,
  createVariant,
} from '@/features/products/api'
import { addPlanItem } from '../api'
import type {
  CanonicalProduct,
  Product,
  ProductVariant,
} from '@/features/products/types'
import { Loader2, Search, Plus, ArrowLeft } from 'lucide-react'

interface Props {
  planId: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function AddPlanItemDialog({ planId, open, onOpenChange }: Props) {
  const [step, setStep] = useState<
    | 'search'
    | 'create-canonical'
    | 'details'
    | 'select-brand'
    | 'create-brand'
    | 'select-variant'
    | 'create-variant'
  >('search')

  const [searchQuery, setSearchQuery] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')

  const [selectedCanonical, setSelectedCanonical] =
    useState<CanonicalProduct | null>(null)
  const [selectedBrand, setSelectedBrand] = useState<Product | null>(null)
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(
    null
  )

  // Form State
  const [quantity, setQuantity] = useState<string>('1')
  const [unit, setUnit] = useState('')
  const [notes, setNotes] = useState('')

  // Creation Form State
  const [newCanonicalName, setNewCanonicalName] = useState('')
  const [newBrandName, setNewBrandName] = useState('')
  const [newVariantName, setNewVariantName] = useState('')
  const [newVariantSize, setNewVariantSize] = useState('')
  const [newVariantUnit, setNewVariantUnit] = useState('')

  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  // Debounce Search
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchQuery)
    }, 500)
    return () => clearTimeout(timer)
  }, [searchQuery])

  // Search Query
  const { data: searchResults, isLoading: isSearching } = useQuery({
    queryKey: ['canonical-products-search', activeInventoryId, debouncedSearch],
    queryFn: () =>
      getCanonicalProducts(activeInventoryId!, {
        search: debouncedSearch,
        limit: 10,
      }),
    enabled:
      !!activeInventoryId && debouncedSearch.length > 0 && step === 'search',
  })

  const { data: products, isLoading: isLoadingProducts } = useQuery({
    queryKey: ['products', activeInventoryId, selectedCanonical?.id],
    queryFn: () => getProducts(activeInventoryId!, selectedCanonical!.id),
    enabled:
      !!activeInventoryId && !!selectedCanonical && step === 'select-brand',
  })

  const { data: variants, isLoading: isLoadingVariants } = useQuery({
    queryKey: ['variants', selectedBrand?.id],
    queryFn: () => getVariants(selectedBrand!.id),
    enabled: !!selectedBrand && step === 'select-variant',
  })

  const createCanonicalMutation = useMutation({
    mutationFn: (data: { name: string }) =>
      createCanonicalProduct(activeInventoryId!, data),
    onSuccess: (newCanonical) => {
      queryClient.invalidateQueries({
        queryKey: ['canonical-products', activeInventoryId],
      })
      setSelectedCanonical(newCanonical)
      setStep('details')
      setNewCanonicalName('')
    },
  })

  const handleCreateCanonical = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newCanonicalName) return
    createCanonicalMutation.mutate({ name: newCanonicalName })
  }

  const createBrandMutation = useMutation({
    mutationFn: (data: { name: string }) =>
      createProduct(activeInventoryId!, {
        name: data.name,
        brand: data.name,
        canonical_product_id: selectedCanonical!.id,
      }),
    onSuccess: (newBrand) => {
      queryClient.invalidateQueries({
        queryKey: ['products', activeInventoryId, selectedCanonical!.id],
      })
      setSelectedBrand(newBrand)
      setStep('select-variant')
      setNewBrandName('')
    },
  })

  const createVariantMutation = useMutation({
    mutationFn: (data: {
      variant_name: string
      unit?: string
      size?: number
    }) => createVariant(selectedBrand!.id, data),
    onSuccess: (newVariant) => {
      queryClient.invalidateQueries({
        queryKey: ['variants', selectedBrand!.id],
      })
      setSelectedVariant(newVariant)
      setStep('details')
      setNewVariantName('')
      setNewVariantSize('')
      setNewVariantUnit('')
    },
  })

  const handleCreateBrand = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newBrandName) return
    createBrandMutation.mutate({ name: newBrandName })
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

  const mutation = useMutation({
    mutationFn: async () => {
      if (selectedVariant) {
        return addPlanItem(planId, {
          target_type: 'product_variant',
          target_id: selectedVariant.id,
          quantity: parseFloat(quantity) || 1,
          unit: unit || undefined,
          note: notes || undefined,
        })
      }

      if (selectedCanonical) {
        return addPlanItem(planId, {
          target_type: 'canonical_product',
          target_id: selectedCanonical.id,
          quantity: parseFloat(quantity) || 1,
          unit: unit || undefined,
          note: notes || undefined,
        })
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['plan', planId],
      })
      resetAndClose()
    },
  })

  const resetAndClose = () => {
    onOpenChange(false)
    setTimeout(() => {
      setStep('search')
      setSearchQuery('')
      setSelectedCanonical(null)
      setSelectedBrand(null)
      setSelectedVariant(null)
      setQuantity('1')
      setUnit('')
      setNotes('')
      setNewCanonicalName('')
      setNewBrandName('')
      setNewVariantName('')
      setNewVariantSize('')
      setNewVariantUnit('')
    }, 300)
  }

  const handleSelectCanonical = (product: CanonicalProduct) => {
    setSelectedCanonical(product)
    setStep('details')
  }

  return (
    <Dialog open={open} onOpenChange={(val) => !val && resetAndClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Item to Plan</DialogTitle>
        </DialogHeader>

        {step === 'search' ? (
          <div className="space-y-4">
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-500" />
              <Input
                placeholder="Search products..."
                className="pl-9"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                autoFocus
              />
            </div>

            <div className="max-h-[300px] overflow-y-auto space-y-2">
              {isSearching && (
                <div className="flex justify-center p-4">
                  <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
                </div>
              )}

              {!isSearching &&
                searchResults?.length === 0 &&
                debouncedSearch && (
                  <div className="text-center p-4 text-slate-500">
                    No products found.
                  </div>
                )}

              {searchResults?.map((product) => (
                <button
                  key={product.id}
                  onClick={() => handleSelectCanonical(product)}
                  className="w-full text-left p-3 rounded-md hover:bg-slate-100 transition-colors border border-transparent hover:border-slate-200"
                >
                  <div className="font-medium text-slate-900">
                    {product.name}
                  </div>
                  {product.description && (
                    <div className="text-sm text-slate-500 truncate">
                      {product.description}
                    </div>
                  )}
                </button>
              ))}

              <Button
                variant="outline"
                className="w-full justify-start mt-2"
                onClick={() => setStep('create-canonical')}
              >
                <Plus className="mr-2 h-4 w-4" />
                Create New Product
              </Button>
            </div>
          </div>
        ) : step === 'create-canonical' ? (
          <form onSubmit={handleCreateCanonical} className="space-y-4">
            {createCanonicalMutation.error && (
              <div className="text-red-500 text-sm">
                {createCanonicalMutation.error instanceof Error
                  ? createCanonicalMutation.error.message
                  : 'Failed to create product'}
              </div>
            )}
            <div className="space-y-2">
              <label className="text-sm font-medium">Product Name</label>
              <Input
                value={newCanonicalName}
                onChange={(e) => setNewCanonicalName(e.target.value)}
                placeholder="e.g. Milk, Bread"
                required
                autoFocus
              />
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setStep('search')}
              >
                Back
              </Button>
              <Button
                type="submit"
                disabled={createCanonicalMutation.isPending}
              >
                {createCanonicalMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Create & Select
              </Button>
            </DialogFooter>
          </form>
        ) : step === 'select-brand' ? (
          <div className="space-y-2">
            <div className="flex items-center gap-2 mb-2">
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => setStep('details')}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <span className="font-medium">
                Select Brand for {selectedCanonical?.name}
              </span>
            </div>
            {isLoadingProducts ? (
              <div className="flex justify-center p-4">
                <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
              </div>
            ) : (
              <>
                {products?.map((p) => (
                  <div
                    key={p.id}
                    className="p-3 border rounded-md hover:bg-slate-50 cursor-pointer transition-colors"
                    onClick={() => {
                      setSelectedBrand(p)
                      setStep('select-variant')
                    }}
                  >
                    <div className="font-medium">{p.name}</div>
                    {p.brand && (
                      <div className="text-sm text-slate-500">{p.brand}</div>
                    )}
                  </div>
                ))}
                <Button
                  variant="outline"
                  className="w-full justify-start mt-2"
                  onClick={() => setStep('create-brand')}
                >
                  <Plus className="mr-2 h-4 w-4" />
                  Create New Brand
                </Button>
              </>
            )}
          </div>
        ) : step === 'create-brand' ? (
          <form onSubmit={handleCreateBrand} className="space-y-4">
            <div className="flex items-center gap-2 mb-2">
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => setStep('select-brand')}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <span className="font-medium">Create Brand</span>
            </div>
            {createBrandMutation.error && (
              <div className="text-red-500 text-sm">
                {createBrandMutation.error instanceof Error
                  ? createBrandMutation.error.message
                  : 'Failed to create brand'}
              </div>
            )}
            <div className="space-y-2">
              <label className="text-sm font-medium">Brand Name</label>
              <Input
                value={newBrandName}
                onChange={(e) => setNewBrandName(e.target.value)}
                placeholder={`e.g. Tesco ${selectedCanonical?.name}`}
                required
                autoFocus
              />
            </div>
            <DialogFooter>
              <Button type="submit" disabled={createBrandMutation.isPending}>
                {createBrandMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Create & Select
              </Button>
            </DialogFooter>
          </form>
        ) : step === 'select-variant' ? (
          <div className="space-y-2">
            <div className="flex items-center gap-2 mb-2">
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => setStep('select-brand')}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <span className="font-medium">
                Select Variant for {selectedBrand?.name}
              </span>
            </div>
            {isLoadingVariants ? (
              <div className="flex justify-center p-4">
                <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
              </div>
            ) : (
              <>
                {variants?.map((v) => (
                  <div
                    key={v.id}
                    className="p-3 border rounded-md hover:bg-slate-50 cursor-pointer transition-colors"
                    onClick={() => {
                      setSelectedVariant(v)
                      setStep('details')
                    }}
                  >
                    <div className="font-medium">{v.variant_name}</div>
                    <div className="text-sm text-slate-500">
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
        ) : step === 'create-variant' ? (
          <form onSubmit={handleCreateVariant} className="space-y-4">
            <div className="flex items-center gap-2 mb-2">
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6"
                onClick={() => setStep('select-variant')}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <span className="font-medium">Create Variant</span>
            </div>
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
              <Button type="submit" disabled={createVariantMutation.isPending}>
                {createVariantMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Create & Select
              </Button>
            </DialogFooter>
          </form>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center justify-between p-3 bg-blue-50 text-blue-900 rounded-md">
              <div>
                <span className="font-medium block">
                  {selectedCanonical?.name}
                </span>
                {selectedVariant && (
                  <span className="text-sm text-blue-700">
                    {selectedBrand?.name} - {selectedVariant.variant_name} (
                    {selectedVariant.size} {selectedVariant.unit})
                  </span>
                )}
              </div>
              <Button
                variant="ghost"
                size="sm"
                className="h-auto p-1 text-blue-700 hover:text-blue-900"
                onClick={() => setStep('search')}
              >
                Change
              </Button>
            </div>

            {!selectedVariant ? (
              <Button
                variant="outline"
                size="sm"
                className="w-full"
                onClick={() => setStep('select-brand')}
              >
                Select Specific Brand/Variant
              </Button>
            ) : (
              <Button
                variant="ghost"
                size="sm"
                className="w-full text-red-500 hover:text-red-700 hover:bg-red-50"
                onClick={() => {
                  setSelectedVariant(null)
                  setSelectedBrand(null)
                }}
              >
                Clear Variant (Use Generic)
              </Button>
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Quantity</label>
                <Input
                  type="number"
                  min="0"
                  step="0.1"
                  value={quantity}
                  onChange={(e) => setQuantity(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Unit (opt)</label>
                <Input
                  placeholder="e.g. kg, pcs"
                  value={unit}
                  onChange={(e) => setUnit(e.target.value)}
                />
              </div>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Notes (opt)</label>
              <Textarea
                placeholder="Brand preference, etc."
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
              />
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={resetAndClose}>
                Cancel
              </Button>
              <Button
                onClick={() => mutation.mutate()}
                disabled={mutation.isPending}
              >
                {mutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Add Item
              </Button>
            </DialogFooter>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
