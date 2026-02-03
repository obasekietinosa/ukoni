import { MemberList } from '../components/member-list'
import { InviteMemberDialog } from '../components/invite-member-dialog'

export function MembersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Members</h1>
          <p className="text-gray-500">
            Manage who has access to this inventory.
          </p>
        </div>
        <InviteMemberDialog />
      </div>
      <MemberList />
    </div>
  )
}
