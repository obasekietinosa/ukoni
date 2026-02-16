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
import { createPlan } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { Loader2 } from 'lucide-react'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  parentPlanId?: string
}

export function CreatePlanDialog({ open, onOpenChange, parentPlanId }: Props) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: async () => {
      if (!activeInventoryId) throw new Error('No inventory selected')
      return createPlan(activeInventoryId, {
        title,
        description: description || undefined,
        parent_plan_id: parentPlanId,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['plans', activeInventoryId],
      })
      if (parentPlanId) {
        queryClient.invalidateQueries({
          queryKey: ['plan', parentPlanId],
        })
      }
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
          <DialogTitle>
            {parentPlanId ? 'Create Sub-plan' : 'Create New Plan'}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="title" className="text-sm font-medium">
              Title
            </label>
            <Input
              id="title"
              placeholder="e.g. Weekly Meal Plan"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={mutation.isPending}
              autoFocus
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="description" className="text-sm font-medium">
              Description (optional)
            </label>
            <Textarea
              id="description"
              placeholder="e.g. Plan for the week of..."
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
            <Button type="submit" disabled={!title.trim() || mutation.isPending}>
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
