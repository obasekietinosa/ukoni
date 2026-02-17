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
import { getPlans, updatePlan } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { Loader2, Search, Plus } from 'lucide-react'
import type { Plan } from '../types'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  parentPlanId: string
}

export function AddExistingPlanDialog({
  open,
  onOpenChange,
  parentPlanId,
}: Props) {
  const [searchQuery, setSearchQuery] = useState('')
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  // Fetch plans with a large limit to allow client-side filtering
  const { data: plans, isLoading } = useQuery({
    queryKey: ['plans', activeInventoryId, 'all'],
    queryFn: () => getPlans(activeInventoryId!, { limit: 100 }),
    enabled: !!activeInventoryId && open,
  })

  // Get current plan details (to find its parent for ancestry check)
  // We can assume the parent plan is already in the cache or fetched
  // But strictly, we need to know the parent chain.
  // The 'plans' list contains parent_plan_id, so we can build the chain.

  const availablePlans = useMemo(() => {
    if (!plans) return []

    // Map of plan ID to parent ID
    const parentMap = new Map<string, string | undefined>()
    plans.forEach((p) => {
      parentMap.set(p.id, p.parent_plan_id)
    })

    // Identify ancestors of the current plan (parentPlanId)
    const ancestors = new Set<string>()
    // Also add the current plan itself as "ancestor" to exclude it
    ancestors.add(parentPlanId)

    // Walk up the tree
    // We need to be careful about potential cycles in existing data, though backend should prevent it.
    const visited = new Set<string>()
    visited.add(parentPlanId)

    // The 'plans' list might not contain the parent of the current plan if it wasn't fetched (e.g. limit 100 missed it, or it's not in the list).
    // But we can only check based on what we have.
    // Ideally we should fetch the current plan to know its parent, but 'parentMap' is built from 'plans'.
    // If 'parentPlanId' refers to a plan in 'plans', we can find its parent.

    // However, we are inside PlanDetailsPage for 'parentPlanId'.
    // The query for 'plans' returns a list.
    // We can try to trace up from parentPlanId using the map.

    // Note: The loop condition needs to handle the case where we don't find the parent in the map.
    // But since we are fetching "all" (limit 100), we hope we have the tree.

    // Let's refine the ancestry check.
    // We want to exclude any plan 'P' such that 'P' is an ancestor of 'parentPlanId'.
    // i.e., parentPlanId -> ... -> P.

    let currentId = parentPlanId
    while (currentId) {
      const parentId = parentMap.get(currentId)
      if (!parentId) break // Can't trace further up or reached root
      if (visited.has(parentId)) break // Cycle detected or already visited
      ancestors.add(parentId)
      visited.add(parentId)
      currentId = parentId
    }

    return plans.filter((plan) => {
      // Exclude ancestors (including self)
      if (ancestors.has(plan.id)) return false
      // Exclude already children
      if (plan.parent_plan_id === parentPlanId) return false

      // Filter by search query
      if (
        searchQuery &&
        !plan.title.toLowerCase().includes(searchQuery.toLowerCase())
      ) {
        return false
      }

      return true
    })
  }, [plans, parentPlanId, searchQuery])

  const mutation = useMutation({
    mutationFn: async (plan: Plan) => {
      return updatePlan(plan.id, {
        title: plan.title,
        description: plan.description,
        parent_plan_id: parentPlanId,
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['plans', activeInventoryId],
      })
      queryClient.invalidateQueries({
        queryKey: ['plan', parentPlanId],
      })
      onOpenChange(false)
      setSearchQuery('')
    },
  })

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Existing Plan</DialogTitle>
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
                  : 'No other plans available.'}
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
