import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { inviteUser } from '../api'
import { Invitation } from '../types'
import { Copy, Loader2, Check } from 'lucide-react'

interface InviteMemberDialogProps {
  inventoryId: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function InviteMemberDialog({
  inventoryId,
  open,
  onOpenChange,
}: InviteMemberDialogProps) {
  const [email, setEmail] = useState('')
  const [role, setRole] = useState('viewer')
  const [invited, setInvited] = useState<Invitation | null>(null)
  const [copied, setCopied] = useState(false)

  const queryClient = useQueryClient()

  const inviteMutation = useMutation({
    mutationFn: () => inviteUser(inventoryId, email, role),
    onSuccess: (data) => {
      setInvited(data)
      queryClient.invalidateQueries({ queryKey: ['members', inventoryId] })
    },
  })

  const reset = () => {
    setEmail('')
    setRole('viewer')
    setInvited(null)
    setCopied(false)
    inviteMutation.reset()
  }

  const handleClose = () => {
    onOpenChange(false)
    // slight delay to reset state after animation
    setTimeout(reset, 300)
  }

  const handleCopy = () => {
    if (!invited) return
    const link = `${window.location.origin}/accept-invite/${invited.id}?token=${invited.token}`
    navigator.clipboard.writeText(link)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Invite Member</DialogTitle>
        </DialogHeader>

        {!invited ? (
          <div className="flex flex-col gap-4 py-4">
            <div className="flex flex-col gap-2">
              <label htmlFor="email" className="text-sm font-medium">
                Email
              </label>
              <Input
                id="email"
                type="email"
                placeholder="colleague@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-2">
              <label htmlFor="role" className="text-sm font-medium">
                Role
              </label>
              <select
                id="role"
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                value={role}
                onChange={(e) => setRole(e.target.value)}
              >
                <option value="viewer">Viewer</option>
                <option value="editor">Editor</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            {inviteMutation.error && (
              <p className="text-sm text-red-500">
                {(inviteMutation.error as Error).message}
              </p>
            )}
          </div>
        ) : (
          <div className="flex flex-col gap-4 py-4">
            <div className="rounded-md bg-green-50 p-4 text-green-900">
              <p className="font-medium">Invitation created!</p>
              <p className="text-sm text-green-700">
                Share this link with {invited.email} to join.
              </p>
            </div>
            <div className="flex items-center gap-2">
              <Input
                readOnly
                value={`${window.location.origin}/accept-invite/${invited.id}?token=${invited.token}`}
              />
              <Button size="icon" variant="outline" onClick={handleCopy}>
                {copied ? (
                  <Check className="h-4 w-4" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </Button>
            </div>
          </div>
        )}

        <DialogFooter>
          {!invited ? (
            <>
              <Button variant="outline" onClick={handleClose}>
                Cancel
              </Button>
              <Button
                onClick={() => inviteMutation.mutate()}
                disabled={inviteMutation.isPending || !email}
              >
                {inviteMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Send Invite
              </Button>
            </>
          ) : (
            <Button onClick={handleClose}>Done</Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
