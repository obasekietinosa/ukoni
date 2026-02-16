import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useInventoryStore } from '@/store/inventory'
import { getPlans } from '../api'
import { PlanList } from '../components/plan-list'
import { CreatePlanDialog } from '../components/create-plan-dialog'
import { Button } from '@/components/ui/button'
import { Plus, Loader2 } from 'lucide-react'

export function PlansPage() {
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const { data: plans, isLoading } = useQuery({
    queryKey: ['plans', activeInventoryId],
    queryFn: () => getPlans(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  // Filter top-level plans (those without a parent)
  // The API supports filtering by parent_plan_id, but current implementation of getPlans
  // takes optional params. If we don't pass parent_plan_id, it might return all plans?
  // Let's check api.ts.
  // getPlans(invId, { parent_plan_id: ... })
  // If I call getPlans(invId), it calls `/inventories/{id}/plans`.
  // The backend `ListPlans` handler checks query param `parent_plan_id`.
  // If it's empty, `parentPlanID` is nil.
  // `PlanModel.List` checks `if parentPlanID != nil`.
  // So if I don't provide it, it returns ALL plans (including sub-plans).
  // Wait, I want only top-level plans here.
  // So I should pass `parent_plan_id` as undefined? No, that would mean "don't filter".
  // Actually, how do I say "where parent_plan_id IS NULL"?
  // The backend `PlanModel.List`:
  // if parentPlanID != nil { query += " AND parent_plan_id = $..." }
  // This filters by EQUALITY.
  // It does NOT filter by IS NULL if parentPlanID is nil.
  // So `List` returns ALL plans if `parentPlanID` is nil.
  // This is a bit problematic if I only want top-level plans.
  // I can filter them client-side.

  const topLevelPlans = plans?.filter((p) => !p.parent_plan_id) || []

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">
            Plans
          </h1>
          <p className="text-slate-500">
            Manage your household activities and meals.
          </p>
        </div>
        <Button onClick={() => setCreateDialogOpen(true)} className="w-full sm:w-auto">
          <Plus className="mr-2 h-4 w-4" />
          Create Plan
        </Button>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-electric-mint" />
        </div>
      ) : (
        <PlanList plans={topLevelPlans} />
      )}

      <CreatePlanDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
      />
    </div>
  )
}
