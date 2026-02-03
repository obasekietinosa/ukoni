import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useInventoryStore } from '@/store/inventory'
import { getInventory, getMembers } from '@/features/inventory/api'
import { Button } from '@/components/ui/button'
import { Plus, Users, Loader2 } from 'lucide-react'
import { InviteMemberDialog } from '../components/invite-member-dialog'
import { MemberList } from '../components/member-list'

export function MembersPage() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const [inviteOpen, setInviteOpen] = useState(false)

  const { data: inventory, isLoading: isInventoryLoading } = useQuery({
    queryKey: ['inventory', activeInventoryId],
    queryFn: () => getInventory(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  const { data: members, isLoading: isMembersLoading } = useQuery({
    queryKey: ['members', activeInventoryId],
    queryFn: () => getMembers(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  if (!activeInventoryId) return null

  if (isInventoryLoading || isMembersLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (!inventory || !members) return <div>Error loading data</div>

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Users className="h-6 w-6" />
          <h1 className="text-2xl font-bold">Members</h1>
        </div>
        <Button onClick={() => setInviteOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Invite Member
        </Button>
      </div>

      <div className="rounded-lg bg-card p-6 shadow-sm">
        <div className="mb-4">
          <h2 className="text-lg font-semibold">Household Members</h2>
          <p className="text-sm text-muted-foreground">
            Manage who has access to {inventory.name}.
          </p>
        </div>

        <MemberList
          inventoryId={activeInventoryId}
          ownerId={inventory.owner_user_id}
          members={members}
        />
      </div>

      <InviteMemberDialog
        inventoryId={activeInventoryId}
        open={inviteOpen}
        onOpenChange={setInviteOpen}
      />
    </div>
  )
}
