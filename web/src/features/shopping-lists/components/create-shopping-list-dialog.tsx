import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createShoppingList } from '../api'
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

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateShoppingListDialog({ open, onOpenChange }: Props) {
  const [name, setName] = useState('')
  const queryClient = useQueryClient()
  const inventoryId = useInventoryStore((state) => state.activeInventoryId)

  const mutation = useMutation({
    mutationFn: (data: { name: string }) => {
      if (!inventoryId) throw new Error('No active inventory')
      return createShoppingList(inventoryId, data.name)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['shopping-lists', inventoryId],
      })
      onOpenChange(false)
      setName('')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) return
    mutation.mutate({ name })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Shopping List</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="name" className="text-sm font-medium">
              List Name
            </label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Weekly Groceries"
              required
              autoFocus
            />
          </div>
          {mutation.error && (
            <div className="text-sm text-red-500">
              {mutation.error instanceof Error
                ? mutation.error.message
                : 'Failed to create list'}
            </div>
          )}
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? 'Creating...' : 'Create List'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
