import { useEffect, useRef } from 'react'
import { useParams, useSearchParams, useNavigate } from 'react-router-dom'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { acceptInvite } from '../api'
import { Button } from '@/components/ui/button'
import { Loader2, CheckCircle, XCircle } from 'lucide-react'

export function AcceptInvitePage() {
  const { id } = useParams<{ id: string }>()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token')
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const attempted = useRef(false)

  const acceptMutation = useMutation({
    mutationFn: () => acceptInvite(id!, token!),
    onSuccess: () => {
      // Invalidate inventories so the new one shows up
      queryClient.invalidateQueries({ queryKey: ['inventories'] })
      setTimeout(() => navigate('/'), 2000)
    },
  })

  useEffect(() => {
    if (id && token && !attempted.current) {
      attempted.current = true
      acceptMutation.mutate()
    }
  }, [id, token, acceptMutation])

  if (!id || !token) {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-4">
        <XCircle className="h-12 w-12 text-red-500" />
        <h1 className="text-xl font-bold">Invalid Invitation Link</h1>
        <p>Missing ID or token.</p>
        <Button onClick={() => navigate('/')}>Go Home</Button>
      </div>
    )
  }

  return (
    <div className="flex h-screen flex-col items-center justify-center gap-4">
      {acceptMutation.isPending && (
        <>
          <Loader2 className="h-12 w-12 animate-spin text-primary" />
          <p>Accepting invitation...</p>
        </>
      )}
      {acceptMutation.isSuccess && (
        <>
          <CheckCircle className="h-12 w-12 text-green-500" />
          <h1 className="text-xl font-bold">Invitation Accepted!</h1>
          <p>Redirecting you to the dashboard...</p>
        </>
      )}
      {acceptMutation.isError && (
        <>
          <XCircle className="h-12 w-12 text-red-500" />
          <h1 className="text-xl font-bold">Failed to Accept Invitation</h1>
          <p className="text-red-500">
            {(acceptMutation.error as Error).message}
          </p>
          <Button onClick={() => navigate('/')}>Go Home</Button>
        </>
      )}
    </div>
  )
}
