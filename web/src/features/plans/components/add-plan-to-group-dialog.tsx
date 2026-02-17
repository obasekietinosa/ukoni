import { useState, useMemo } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { getPlans, addPlanToGroup } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { Loader2, Search, Plus } from 'lucide-react'
import type { Plan } from '../types'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  groupId: string
  existingPlanIds: string[]
}

export function AddPlanToGroupDialog({
  open,
  onOpenChange,
  groupId,
  existingPlanIds,
}: Props) {
  const [searchQuery, setSearchQuery] = useState('')
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const { data: plans, isLoading } = useQuery({
    queryKey: ['plans', activeInventoryId, 'all'],
    queryFn: () => getPlans(activeInventoryId!, { limit: 100 }),
    enabled: !!activeInventoryId && open,
  })

  const availablePlans = useMemo(() => {
    if (!plans) return []

    return plans.filter((plan) => {
      // Exclude already added plans
      if (existingPlanIds.includes(plan.id)) return false

      // Filter by search query
      if (
        searchQuery &&
        !plan.title.toLowerCase().includes(searchQuery.toLowerCase())
      ) {
        return false
      }

      return true
    })
  }, [plans, existingPlanIds, searchQuery])

  const mutation = useMutation({
    mutationFn: async (plan: Plan) => {
      return addPlanToGroup(groupId, plan.id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['plan-group', groupId],
      })
      onOpenChange(false)
      setSearchQuery('')
    },
  })

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Plan to Group</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="relative">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-500" />
            <Input
              placeholder="Search plans..."
              className="pl-9"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              autoFocus
            />
          </div>

          <div className="max-h-[300px] overflow-y-auto space-y-2">
            {isLoading && (
              <div className="flex justify-center p-4">
                <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
              </div>
            )}

            {!isLoading && availablePlans.length === 0 && (
              <div className="text-center p-4 text-slate-500">
                {searchQuery
                  ? 'No matching plans found.'
                  : 'No available plans to add.'}
              </div>
            )}

            {availablePlans.map((plan) => (
              <div
                key={plan.id}
                className="flex items-center justify-between p-3 rounded-md hover:bg-slate-50 border border-transparent hover:border-slate-200 transition-colors"
              >
                <div className="min-w-0 flex-1 mr-4">
                  <div className="font-medium text-slate-900 truncate">
                    {plan.title}
                  </div>
                  {plan.description && (
                    <div className="text-sm text-slate-500 truncate">
                      {plan.description}
                    </div>
                  )}
                </div>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => mutation.mutate(plan)}
                  disabled={mutation.isPending}
                >
                  {mutation.isPending && mutation.variables?.id === plan.id ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Plus className="h-4 w-4" />
                  )}
                  Add
                </Button>
              </div>
            ))}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
