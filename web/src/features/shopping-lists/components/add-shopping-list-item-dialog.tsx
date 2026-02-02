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
import { getCanonicalProducts } from '@/features/products/api'
import { addShoppingListItem } from '../api'
import type { CanonicalProduct } from '@/features/products/types'
import { Loader2, Search } from 'lucide-react'

interface Props {
  listId: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function AddShoppingListItemDialog({
  listId,
  open,
  onOpenChange,
}: Props) {
  const [step, setStep] = useState<'search' | 'details'>('search')
  const [searchQuery, setSearchQuery] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')
  const [selectedProduct, setSelectedProduct] =
    useState<CanonicalProduct | null>(null)

  // Form State
  const [quantity, setQuantity] = useState<string>('1')
  const [unit, setUnit] = useState('')
  const [notes, setNotes] = useState('')

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

  const mutation = useMutation({
    mutationFn: async () => {
      if (!selectedProduct) return
      // Default to adding as canonical product for now
      // TODO: Add support for specific variant selection
      return addShoppingListItem(listId, {
        target_type: 'canonical_product',
        target_id: selectedProduct.id,
        quantity: parseFloat(quantity) || 1,
        unit: unit || undefined,
        notes: notes || undefined,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['shopping-list-items', listId],
      })
      resetAndClose()
    },
  })

  const resetAndClose = () => {
    onOpenChange(false)
    setTimeout(() => {
      setStep('search')
      setSearchQuery('')
      setSelectedProduct(null)
      setQuantity('1')
      setUnit('')
      setNotes('')
    }, 300)
  }

  const handleSelectProduct = (product: CanonicalProduct) => {
    setSelectedProduct(product)
    setStep('details')
  }

  return (
    <Dialog open={open} onOpenChange={(val) => !val && resetAndClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Item</DialogTitle>
        </DialogHeader>

        {step === 'search' ? (
          <div className="space-y-4">
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
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
                  <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                </div>
              )}

              {!isSearching &&
                searchResults?.length === 0 &&
                debouncedSearch && (
                  <div className="text-center p-4 text-gray-500">
                    No products found.
                  </div>
                )}

              {searchResults?.map((product) => (
                <button
                  key={product.id}
                  onClick={() => handleSelectProduct(product)}
                  className="w-full text-left p-3 rounded-md hover:bg-gray-100 transition-colors border"
                >
                  <div className="font-medium">{product.name}</div>
                  {product.description && (
                    <div className="text-sm text-gray-500 truncate">
                      {product.description}
                    </div>
                  )}
                </button>
              ))}
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center justify-between p-3 bg-blue-50 text-blue-900 rounded-md">
              <span className="font-medium">{selectedProduct?.name}</span>
              <Button
                variant="ghost"
                size="sm"
                className="h-auto p-1 text-blue-700 hover:text-blue-900"
                onClick={() => setStep('search')}
              >
                Change
              </Button>
            </div>

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
              <textarea
                className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
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
