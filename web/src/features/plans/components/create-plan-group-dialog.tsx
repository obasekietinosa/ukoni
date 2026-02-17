import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
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
import { createPlanGroup } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { Loader2 } from 'lucide-react'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreatePlanGroupDialog({ open, onOpenChange }: Props) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: async () => {
      if (!activeInventoryId) throw new Error('No inventory selected')
      return createPlanGroup(activeInventoryId, {
        title,
        description: description || undefined,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['plan-groups', activeInventoryId],
      })
      onOpenChange(false)
      setTitle('')
      setDescription('')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim()) return
    mutation.mutate()
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create Plan Group</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="group-title" className="text-sm font-medium">
              Title
            </label>
            <Input
              id="group-title"
              placeholder="e.g. Weekly Meal Plans"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={mutation.isPending}
              autoFocus
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="group-description" className="text-sm font-medium">
              Description (optional)
            </label>
            <Textarea
              id="group-description"
              placeholder="e.g. Group for weekly meals..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={mutation.isPending}
            />
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={mutation.isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={!title.trim() || mutation.isPending}
            >
              {mutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Create
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
