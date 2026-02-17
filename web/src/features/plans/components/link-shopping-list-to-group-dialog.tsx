import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useInventoryStore } from '@/store/inventory'
import { getShoppingLists } from '@/features/shopping-lists/api'
import { linkShoppingListToGroup } from '../api'
import { Loader2 } from 'lucide-react'

interface Props {
  groupId: string
  linkedListIds: string[]
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function LinkShoppingListToGroupDialog({
  groupId,
  linkedListIds,
  open,
  onOpenChange,
}: Props) {
  const [selectedListId, setSelectedListId] = useState<string>('')
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const { data: shoppingLists, isLoading } = useQuery({
    queryKey: ['shopping-lists', activeInventoryId],
    queryFn: () => getShoppingLists(activeInventoryId!),
    enabled: !!activeInventoryId && open,
  })

  const mutation = useMutation({
    mutationFn: async () => {
      if (!selectedListId) return
      return linkShoppingListToGroup(groupId, selectedListId)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['plan-group', groupId],
      })
      onOpenChange(false)
      setSelectedListId('')
    },
  })

  const availableLists =
    shoppingLists?.filter((l) => !linkedListIds.includes(l.id)) || []

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Link Shopping List</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {isLoading ? (
            <div className="flex justify-center p-4">
              <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
            </div>
          ) : availableLists.length === 0 ? (
            <div className="text-center p-4 text-slate-500">
              No available shopping lists to link.
            </div>
          ) : (
            <div className="space-y-2">
              <label className="text-sm font-medium">Select List</label>
              <div className="grid gap-2 max-h-[300px] overflow-y-auto">
                {availableLists?.map((list) => (
                  <button
                    key={list.id}
                    onClick={() => setSelectedListId(list.id)}
                    className={`p-3 text-left border rounded-md transition-colors ${selectedListId === list.id ? 'border-electric-mint bg-electric-mint/10' : 'hover:bg-slate-50'}`}
                  >
                    <div className="font-medium text-slate-900">
                      {list.name}
                    </div>
                    <div className="text-xs text-slate-500">
                      {new Date(list.last_updated_at).toLocaleDateString()}
                    </div>
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={() => mutation.mutate()}
            disabled={!selectedListId || mutation.isPending}
          >
            {mutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            Link List
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
