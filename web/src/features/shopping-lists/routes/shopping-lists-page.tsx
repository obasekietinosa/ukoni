import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Plus, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useInventoryStore } from '@/store/inventory'
import { getShoppingLists } from '../api'
import { ShoppingListCard } from '../components/shopping-list-card'
import { CreateShoppingListDialog } from '../components/create-shopping-list-dialog'

export function ShoppingListsPage() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  const {
    data: lists,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['shopping-lists', activeInventoryId],
    queryFn: () => getShoppingLists(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  if (!activeInventoryId) return null

  if (isLoading) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-8 text-center text-red-500">
        Failed to load shopping lists.
      </div>
    )
  }

  return (
    <div className="container mx-auto max-w-2xl py-6">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold">Shopping Lists</h1>
        <Button onClick={() => setIsCreateOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          New List
        </Button>
      </div>

      <div className="space-y-4">
        {lists?.length === 0 ? (
          <div className="text-center py-12 border rounded-lg bg-gray-50 text-gray-500">
            You don't have any shopping lists yet.
          </div>
        ) : (
          lists?.map((list) => <ShoppingListCard key={list.id} list={list} />)
        )}
      </div>

      <CreateShoppingListDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />
    </div>
  )
}
