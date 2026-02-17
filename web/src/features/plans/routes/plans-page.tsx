import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useInventoryStore } from '@/store/inventory'
import { getPlans, getPlanGroups } from '../api'
import { PlanList } from '../components/plan-list'
import { PlanGroupList } from '../components/plan-group-list'
import { CreatePlanDialog } from '../components/create-plan-dialog'
import { CreatePlanGroupDialog } from '../components/create-plan-group-dialog'
import { Button } from '@/components/ui/button'
import { Plus, Loader2 } from 'lucide-react'

export function PlansPage() {
  const [createPlanOpen, setCreatePlanOpen] = useState(false)
  const [createGroupOpen, setCreateGroupOpen] = useState(false)
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const { data: plans, isLoading: plansLoading } = useQuery({
    queryKey: ['plans', activeInventoryId],
    queryFn: () => getPlans(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  const { data: groups, isLoading: groupsLoading } = useQuery({
    queryKey: ['plan-groups', activeInventoryId],
    queryFn: () => getPlanGroups(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  const isLoading = plansLoading || groupsLoading

  return (
    <div className="space-y-8">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">
            Plans & Groups
          </h1>
          <p className="text-slate-500">
            Manage your household activities, meals, and plan groups.
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Button
            variant="outline"
            onClick={() => setCreateGroupOpen(true)}
            className="flex-1 sm:flex-none"
          >
            <Plus className="mr-2 h-4 w-4" />
            New Group
          </Button>
          <Button
            onClick={() => setCreatePlanOpen(true)}
            className="flex-1 sm:flex-none"
          >
            <Plus className="mr-2 h-4 w-4" />
            New Plan
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-electric-mint" />
        </div>
      ) : (
        <>
          <section className="space-y-4">
            <h2 className="text-lg font-semibold text-slate-900">Plan Groups</h2>
            <PlanGroupList groups={groups || []} />
          </section>

          <section className="space-y-4">
            <h2 className="text-lg font-semibold text-slate-900">All Plans</h2>
            <PlanList plans={plans || []} />
          </section>
        </>
      )}

      <CreatePlanDialog
        open={createPlanOpen}
        onOpenChange={setCreatePlanOpen}
      />
      <CreatePlanGroupDialog
        open={createGroupOpen}
        onOpenChange={setCreateGroupOpen}
      />
    </div>
  )
}
