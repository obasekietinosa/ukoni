import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  getPlanGroup,
  deletePlanGroup,
  createShoppingListFromGroup,
  unlinkShoppingListFromGroup,
} from '../api'
import { getShoppingLists } from '@/features/shopping-lists/api'
import { PlanList } from '../components/plan-list'
import { Button } from '@/components/ui/button'
import { Loader2, Plus, Trash2, ShoppingBag, Unlink } from 'lucide-react'
import { AddPlanToGroupDialog } from '../components/add-plan-to-group-dialog'
import { LinkShoppingListToGroupDialog } from '../components/link-shopping-list-to-group-dialog'

export function PlanGroupDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const [addPlanOpen, setAddPlanOpen] = useState(false)
  const [linkListOpen, setLinkListOpen] = useState(false)

  const { data: group, isLoading } = useQuery({
    queryKey: ['plan-group', id],
    queryFn: () => getPlanGroup(id!),
    enabled: !!id,
  })

  // Fetch all shopping lists to match IDs
  const { data: allShoppingLists } = useQuery({
    queryKey: ['shopping-lists', group?.inventory_id],
    queryFn: () => getShoppingLists(group!.inventory_id),
    enabled: !!group?.inventory_id,
  })

  const linkedLists = allShoppingLists?.filter((l) =>
    group?.shopping_lists?.includes(l.id)
  )

  const deleteMutation = useMutation({
    mutationFn: deletePlanGroup,
    onSuccess: () => {
      navigate('/plans')
      queryClient.invalidateQueries({ queryKey: ['plan-groups'] })
    },
  })

  const unlinkMutation = useMutation({
    mutationFn: (listId: string) => unlinkShoppingListFromGroup(id!, listId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan-group', id] })
    },
  })

  const createListMutation = useMutation({
    mutationFn: (groupId: string) => createShoppingListFromGroup(groupId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plan-group', id] })
      queryClient.invalidateQueries({ queryKey: ['shopping-lists'] })
    },
  })

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-electric-mint" />
      </div>
    )
  }

  if (!group) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <h2 className="text-xl font-semibold text-slate-900">
          Group not found
        </h2>
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
            <Link
              to="/plans"
              className="hover:text-slate-900 transition-colors"
            >
              Plans
            </Link>
            <span>/</span>
            <span className="font-medium text-slate-900 truncate max-w-[200px]">
              {group.title}
            </span>
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">
            {group.title}
          </h1>
          {group.description && (
            <p className="text-slate-500 max-w-2xl">{group.description}</p>
          )}
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              if (confirm('Delete this group?')) {
                deleteMutation.mutate(group.id)
              }
            }}
            className="text-red-600 hover:bg-red-50 hover:text-red-700 border-red-200"
          >
            <Trash2 className="mr-2 h-4 w-4" />
            Delete Group
          </Button>
        </div>
      </div>

      {/* Plans Section */}
      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-slate-900">Plans</h2>
          <div className="flex items-center gap-2">
            {group.plans && group.plans.length > 0 && (
              <Button
                size="sm"
                variant="outline"
                onClick={() => createListMutation.mutate(group.id)}
                disabled={createListMutation.isPending}
              >
                {createListMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <ShoppingBag className="mr-2 h-4 w-4" />
                )}
                Create Shopping List
              </Button>
            )}
            <Button
              size="sm"
              variant="outline"
              onClick={() => setAddPlanOpen(true)}
            >
              <Plus className="mr-2 h-4 w-4" />
              Add Plan
            </Button>
          </div>
        </div>
        {group.plans && group.plans.length > 0 ? (
          <PlanList plans={group.plans} />
        ) : (
          <div className="text-sm text-slate-500 italic">
            No plans in this group.
          </div>
        )}
      </section>

      {/* Linked Shopping Lists Section */}
      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-slate-900">
            Linked Shopping Lists
          </h2>
          <Button
            size="sm"
            variant="outline"
            onClick={() => setLinkListOpen(true)}
          >
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
                <Link
                  to={`/shopping-lists/${list.id}`}
                  className="flex items-center gap-3 flex-1 min-w-0"
                >
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
          <div className="text-sm text-slate-500 italic">
            No linked shopping lists.
          </div>
        )}
      </section>

      <AddPlanToGroupDialog
        open={addPlanOpen}
        onOpenChange={setAddPlanOpen}
        groupId={group.id}
        existingPlanIds={group.plans?.map((p) => p.id) || []}
      />
      <LinkShoppingListToGroupDialog
        open={linkListOpen}
        onOpenChange={setLinkListOpen}
        groupId={group.id}
        linkedListIds={group.shopping_lists || []}
      />
    </div>
  )
}
