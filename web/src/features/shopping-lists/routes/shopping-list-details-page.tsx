import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useQuery, useMutation } from '@tanstack/react-query'
import { Plus, ArrowLeft, Loader2, Trash2, ShoppingBag } from 'lucide-react'
import { Button } from '@/components/ui/button'
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
            </Button>
            <Button onClick={() => setIsAddOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Add Item
            </Button>
          </div>
        </div>
      </div>

      {isItemsLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
        </div>
      ) : (
        <ShoppingListItemList items={items || []} listId={list.id} />
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
          onSuccess={() => {
            // Maybe navigate to transactions list or show success message?
            // For now, staying on the list is fine as items are not automatically deleted
            // (unless backend does it, which we discussed it doesn't yet)
          }}
        />
      )}
    </div>
  )
}
