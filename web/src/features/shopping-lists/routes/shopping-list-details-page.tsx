import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useQuery, useMutation } from '@tanstack/react-query'
import {
  Plus,
  ArrowLeft,
  Loader2,
  Trash2,
  ShoppingBag,
  Search,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  getShoppingList,
  getShoppingListItems,
  deleteShoppingList,
} from '../api'
import { ShoppingListItemList } from '../components/shopping-list-item-list'
import { AddShoppingListItemDialog } from '../components/add-shopping-list-item-dialog'
import { CreateTransactionFromListDialog } from '@/features/transactions/components/create-transaction-from-list-dialog'

export function ShoppingListDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [isTransactionOpen, setIsTransactionOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [shoppedItems, setShoppedItems] = useState<Set<string>>(new Set())

  const { data: list, isLoading: isListLoading } = useQuery({
    queryKey: ['shopping-list', id],
    queryFn: () => getShoppingList(id!),
    enabled: !!id,
  })

  const { data: items, isLoading: isItemsLoading } = useQuery({
    queryKey: ['shopping-list-items', id],
    queryFn: () => getShoppingListItems(id!),
    enabled: !!id,
  })

  const deleteListMutation = useMutation({
    mutationFn: deleteShoppingList,
    onSuccess: () => {
      navigate('/shopping-lists')
    },
  })

  const toggleShopped = (itemId: string) => {
    setShoppedItems((prev) => {
      const next = new Set(prev)
      if (next.has(itemId)) {
        next.delete(itemId)
      } else {
        next.add(itemId)
      }
      return next
    })
  }

  const filteredItems = items?.filter((item) => {
    const term = searchTerm.toLowerCase()
    if (!term) return true

    // Check canonical product name
    if (
      item.target_type === 'canonical_product' &&
      item.canonical_product?.name.toLowerCase().includes(term)
    ) {
      return true
    }
    // Check variant name and brand
    if (item.target_type === 'product_variant') {
      const variantName = item.product_variant?.variant_name.toLowerCase() || ''
      const brand = item.product?.brand.toLowerCase() || ''
      if (variantName.includes(term) || brand.includes(term)) {
        return true
      }
    }
    // Check notes
    if (item.notes?.toLowerCase().includes(term)) {
      return true
    }
    return false
  })

  if (isListLoading) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    )
  }

  if (!list) return <div>List not found</div>

  return (
    <div className="container mx-auto max-w-2xl py-6">
      <div className="mb-6">
        <Link
          to="/shopping-lists"
          className="inline-flex items-center text-sm text-gray-500 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="mr-1 h-4 w-4" />
          Back to Lists
        </Link>
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <h1 className="text-2xl font-bold">{list.name}</h1>
          <div className="flex flex-wrap items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              className="text-red-500 hover:text-red-700 hover:bg-red-50 border-red-200"
              onClick={() => {
                if (confirm('Delete this shopping list?')) {
                  deleteListMutation.mutate(list.id)
                }
              }}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete List
            </Button>
            <Button
              variant="secondary"
              onClick={() => setIsTransactionOpen(true)}
              disabled={!items || items.length === 0}
            >
              <ShoppingBag className="mr-2 h-4 w-4" />
              Complete Shopping
              {shoppedItems.size > 0 && ` (${shoppedItems.size})`}
            </Button>
            <Button onClick={() => setIsAddOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Add Item
            </Button>
          </div>
        </div>

        {/* Search Bar */}
        <div className="mt-4 relative">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-500" />
          <Input
            placeholder="Search items..."
            className="pl-9 bg-white"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>

      {isItemsLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
        </div>
      ) : (
        <ShoppingListItemList
          items={filteredItems || []}
          listId={list.id}
          shoppedItems={shoppedItems}
          onToggleShopped={toggleShopped}
        />
      )}

      <AddShoppingListItemDialog
        listId={list.id}
        open={isAddOpen}
        onOpenChange={setIsAddOpen}
      />

      {isTransactionOpen && (
        <CreateTransactionFromListDialog
          open={isTransactionOpen}
          onOpenChange={setIsTransactionOpen}
          items={items || []}
          initialSelection={shoppedItems}
          onSuccess={() => {
            // Maybe navigate to transactions list or show success message?
            // For now, staying on the list is fine as items are not automatically deleted
            // (unless backend does it, which we discussed it doesn't yet)
            setShoppedItems(new Set()) // Reset shopped items
          }}
        />
      )}
    </div>
  )
}
