import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getPlan, deletePlan, unlinkShoppingList } from '../api'
import { getShoppingLists } from '@/features/shopping-lists/api'
import { PlanItemList } from '../components/plan-item-list'
import { PlanList } from '../components/plan-list'
import { CreatePlanDialog } from '../components/create-plan-dialog'
import { AddPlanItemDialog } from '../components/add-plan-item-dialog'
import { LinkShoppingListDialog } from '../components/link-shopping-list-dialog'
import { Button } from '@/components/ui/button'
import {
  Loader2,
  Plus,
  Trash2,
  ShoppingBag,
  Unlink,
} from 'lucide-react'

export function PlanDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const [createSubPlanOpen, setCreateSubPlanOpen] = useState(false)
  const [addItemOpen, setAddItemOpen] = useState(false)
  const [linkListOpen, setLinkListOpen] = useState(false)

  const { data: plan, isLoading } = useQuery({
    queryKey: ['plan', id],
    queryFn: () => getPlan(id!),
    enabled: !!id,
  })

  // Fetch all shopping lists to match IDs
  const { data: allShoppingLists } = useQuery({
    queryKey: ['shopping-lists', plan?.inventory_id],
    queryFn: () => getShoppingLists(plan!.inventory_id),
    enabled: !!plan?.inventory_id,
  })

  const linkedLists = allShoppingLists?.filter((l) =>
    plan?.shopping_lists?.includes(l.id)
  )

  const deleteMutation = useMutation({
    mutationFn: deletePlan,
    onSuccess: () => {
      navigate('/plans')
      queryClient.invalidateQueries({ queryKey: ['plans'] })
    },
  })

  const unlinkMutation = useMutation({
    mutationFn: (listId: string) => unlinkShoppingList(id!, listId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan', id] })
    },
  })

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-electric-mint" />
      </div>
    )
  }

  if (!plan) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <h2 className="text-xl font-semibold text-slate-900">Plan not found</h2>
        <Button variant="link" onClick={() => navigate('/plans')}>
          Back to Plans
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-8 pb-12">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div className="space-y-1">
          <div className="flex items-center gap-2 text-sm text-slate-500 mb-2">
            <Link to="/plans" className="hover:text-slate-900 transition-colors">
              Plans
            </Link>
            <span>/</span>
            <span className="font-medium text-slate-900 truncate max-w-[200px]">
              {plan.title}
            </span>
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">
            {plan.title}
          </h1>
          {plan.description && (
            <p className="text-slate-500 max-w-2xl">{plan.description}</p>
          )}
        </div>
        <div className="flex items-center gap-2">
           <Button
            variant="outline"
            size="sm"
            onClick={() => {
              if (confirm('Delete this plan?')) {
                deleteMutation.mutate(plan.id)
              }
            }}
            className="text-red-600 hover:bg-red-50 hover:text-red-700 border-red-200"
          >
            <Trash2 className="mr-2 h-4 w-4" />
            Delete Plan
          </Button>
        </div>
      </div>

      {/* Items Section */}
      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-slate-900">Items</h2>
          <Button size="sm" onClick={() => setAddItemOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Item
          </Button>
        </div>
        <PlanItemList items={plan.items || []} planId={plan.id} />
      </section>

      {/* Sub-plans Section */}
      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-slate-900">Sub-plans</h2>
          <Button size="sm" variant="outline" onClick={() => setCreateSubPlanOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Create Sub-plan
          </Button>
        </div>
        {plan.children && plan.children.length > 0 ? (
           <PlanList plans={plan.children} />
        ) : (
           <div className="text-sm text-slate-500 italic">No sub-plans.</div>
        )}
      </section>

       {/* Linked Shopping Lists Section */}
       <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-slate-900">Linked Shopping Lists</h2>
          <Button size="sm" variant="outline" onClick={() => setLinkListOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Link List
          </Button>
        </div>

        {linkedLists && linkedLists.length > 0 ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {linkedLists.map((list) => (
              <div
                key={list.id}
                className="group relative flex items-center justify-between overflow-hidden rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition-all hover:border-slate-300 hover:shadow-md"
              >
                <Link to={`/shopping-lists/${list.id}`} className="flex items-center gap-3 flex-1 min-w-0">
                  <div className="flex h-10 w-10 items-center justify-center rounded-full bg-slate-50 text-slate-400 group-hover:bg-electric-mint/10 group-hover:text-electric-mint transition-colors">
                    <ShoppingBag className="h-5 w-5" />
                  </div>
                  <div className="truncate">
                    <h3 className="font-medium text-slate-900 group-hover:text-electric-mint transition-colors">
                      {list.name}
                    </h3>
                    <p className="text-xs text-slate-500">
                      {new Date(list.last_updated_at).toLocaleDateString()}
                    </p>
                  </div>
                </Link>
                 <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-full"
                    onClick={() => {
                        if (confirm(`Unlink ${list.name}?`)) {
                            unlinkMutation.mutate(list.id)
                        }
                    }}
                    title="Unlink list"
                >
                    <Unlink className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-sm text-slate-500 italic">No linked shopping lists.</div>
        )}
      </section>

      <CreatePlanDialog
        open={createSubPlanOpen}
        onOpenChange={setCreateSubPlanOpen}
        parentPlanId={plan.id}
      />
      <AddPlanItemDialog
        open={addItemOpen}
        onOpenChange={setAddItemOpen}
        planId={plan.id}
      />
      <LinkShoppingListDialog
        open={linkListOpen}
        onOpenChange={setLinkListOpen}
        planId={plan.id}
        linkedListIds={plan.shopping_lists || []}
      />
    </div>
  )
}
