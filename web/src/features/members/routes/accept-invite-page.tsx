import { useState } from 'react'
import { useParams, useSearchParams, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { acceptInvite } from '../api'
import { Button } from '@/components/ui/button'
import { CheckCircle2, XCircle } from 'lucide-react'

export function AcceptInvitePage() {
  const { id } = useParams<{ id: string }>()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token')
  const navigate = useNavigate()
  const [mutationError, setMutationError] = useState<string | null>(null)

  const isInvalidLink = !id || !token

  const acceptMutation = useMutation({
    mutationFn: async () => {
      if (!id || !token) throw new Error('Invalid invitation link')
      return acceptInvite(id, token)
    },
    onSuccess: () => {
      // Redirect to home/inventory selection after short delay
      setTimeout(() => navigate('/'), 2000)
    },
    onError: (err: Error) => {
      setMutationError(err.message || 'Failed to accept invitation')
    },
  })

  const error = isInvalidLink ? 'Invalid invitation link' : mutationError

  return (
    <div className="flex h-screen w-full items-center justify-center bg-gray-50">
      <div className="w-full max-w-md rounded-lg border bg-white p-8 shadow-sm text-center space-y-6">
        {acceptMutation.isSuccess ? (
          <>
            <div className="flex justify-center text-green-500">
              <CheckCircle2 className="h-16 w-16" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900">
              Invitation Accepted!
            </h1>
            <p className="text-gray-500">
              You have successfully joined the inventory. Redirecting...
            </p>
          </>
        ) : error ? (
          <>
            <div className="flex justify-center text-red-500">
              <XCircle className="h-16 w-16" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900">Error</h1>
            <p className="text-red-500">{error}</p>
            <Button onClick={() => navigate('/')} variant="outline">
              Go Home
            </Button>
          </>
        ) : (
          <>
            <h1 className="text-2xl font-bold text-gray-900">
              Accept Invitation
            </h1>
            <p className="text-gray-500">
              You have been invited to join an inventory.
            </p>
            <div className="pt-4">
              <Button
                size="lg"
                className="w-full"
                onClick={() => acceptMutation.mutate()}
                disabled={acceptMutation.isPending}
              >
                {acceptMutation.isPending ? 'Accepting...' : 'Join Inventory'}
              </Button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
