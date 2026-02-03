import { useMutation, useQueryClient } from '@tanstack/react-query'
import { InventoryMembership } from '@/features/inventory/types'
import { removeMember, updateMemberRole } from '../api'
import { Button } from '@/components/ui/button'
import { Trash2, Loader2 } from 'lucide-react'
import { useAuthStore } from '@/store/auth'

interface MemberListProps {
  inventoryId: string
  ownerId: string
  members: InventoryMembership[]
}

export function MemberList({ inventoryId, ownerId, members }: MemberListProps) {
  const user = useAuthStore((state) => state.user)
  const isOwner = user?.id === ownerId
  const queryClient = useQueryClient()

  const removeMutation = useMutation({
    mutationFn: (userId: string) => removeMember(inventoryId, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members', inventoryId] })
    },
  })

  const updateRoleMutation = useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: string }) =>
      updateMemberRole(inventoryId, userId, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members', inventoryId] })
    },
  })

  return (
    <div className="rounded-md border">
      <div className="grid grid-cols-12 border-b bg-muted/50 p-4 text-sm font-medium">
        <div className="col-span-4">Name</div>
        <div className="col-span-4">Email</div>
        <div className="col-span-3">Role</div>
        <div className="col-span-1"></div>
      </div>
      <div className="divide-y">
        {members.map((member) => {
          const isMemberOwner = member.user_id === ownerId
          const isMe = member.user_id === user?.id

          return (
            <div
              key={member.id}
              className="grid grid-cols-12 items-center p-4 text-sm"
            >
              <div className="col-span-4 font-medium">
                {member.user_name || 'Unknown'}
                {isMemberOwner && (
                  <span className="ml-2 rounded bg-yellow-100 px-2 py-0.5 text-xs font-normal text-yellow-800">
                    Owner
                  </span>
                )}
                {isMe && (
                  <span className="ml-2 text-xs text-muted-foreground">
                    (You)
                  </span>
                )}
              </div>
              <div className="col-span-4 text-muted-foreground">
                {member.user_email}
              </div>
              <div className="col-span-3">
                <select
                  className="w-full max-w-[120px] rounded-md border border-input bg-transparent px-2 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                  value={member.role}
                  disabled={
                    !isOwner || isMemberOwner || removeMutation.isPending
                  }
                  onChange={(e) =>
                    updateRoleMutation.mutate({
                      userId: member.user_id,
                      role: e.target.value,
                    })
                  }
                >
                  <option value="viewer">Viewer</option>
                  <option value="editor">Editor</option>
                  <option value="admin">Admin</option>
                </select>
              </div>
              <div className="col-span-1 flex justify-end">
                {isOwner && !isMemberOwner && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-red-500 hover:bg-red-50 hover:text-red-600"
                    onClick={() => removeMutation.mutate(member.user_id)}
                    disabled={removeMutation.isPending}
                  >
                    {removeMutation.isPending ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Trash2 className="h-4 w-4" />
                    )}
                  </Button>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
