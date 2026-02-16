import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { listMembers, removeMember, updateMemberRole } from '../api'
import { useInventoryStore } from '@/store/inventory'
import { useCurrentUserRole } from '@/features/inventory/hooks/use-current-user-role'

export function MemberList() {
  const inventoryId = useInventoryStore((state) => state.activeInventoryId)
  const queryClient = useQueryClient()
  const { role: currentUserRole } = useCurrentUserRole(inventoryId)

  const { data: members, isLoading } = useQuery({
    queryKey: ['members', inventoryId],
    queryFn: () => listMembers(inventoryId!),
    enabled: !!inventoryId,
  })

  const removeMutation = useMutation({
    mutationFn: (userId: string) => removeMember(inventoryId!, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members', inventoryId] })
    },
  })

  const updateRoleMutation = useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: string }) =>
      updateMemberRole(inventoryId!, userId, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members', inventoryId] })
    },
  })

  if (isLoading) {
    return (
      <div className="p-4 text-center text-gray-500">Loading members...</div>
    )
  }

  // Assuming 'admin' and 'owner' (if it exists as a role string, or just implied by ownership) can manage.
  // Backend returns role string.
  const canManage = currentUserRole === 'admin' || currentUserRole === 'owner'

  return (
    <div className="space-y-4">
      <div className="rounded-md border bg-white shadow-sm">
        <table className="w-full text-sm text-left">
          <thead className="bg-gray-50 text-gray-500">
            <tr>
              <th className="px-4 py-3 font-medium">Name</th>
              <th className="px-4 py-3 font-medium">Email</th>
              <th className="px-4 py-3 font-medium">Role</th>
              <th className="px-4 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {members?.map((member) => (
              <tr key={member.id}>
                <td className="px-4 py-3 font-medium text-gray-900">
                  {member.user_name || 'Pending...'}
                </td>
                <td className="px-4 py-3 text-gray-500">
                  {member.user_email || 'Hidden'}
                </td>
                <td className="px-4 py-3 text-gray-500">
                  {canManage && member.role !== 'owner' ? (
                    <select
                      value={member.role}
                      onChange={(e) =>
                        updateRoleMutation.mutate({
                          userId: member.user_id,
                          role: e.target.value,
                        })
                      }
                      className="rounded border border-gray-200 bg-white px-2 py-1 text-base md:text-xs focus:border-gray-400 focus:outline-none"
                    >
                      <option value="viewer">Viewer</option>
                      <option value="editor">Editor</option>
                      <option value="admin">Admin</option>
                    </select>
                  ) : (
                    <span className="capitalize">{member.role}</span>
                  )}
                </td>
                <td className="px-4 py-3 text-right">
                  {canManage && member.role !== 'owner' && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => {
                        if (confirm('Remove this member?')) {
                          removeMutation.mutate(member.user_id)
                        }
                      }}
                      className="text-red-600 hover:text-red-700 hover:bg-red-50"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </td>
              </tr>
            ))}
            {members?.length === 0 && (
              <tr>
                <td colSpan={4} className="px-4 py-8 text-center text-gray-500">
                  No members found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
