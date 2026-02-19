import { useAuthStore } from '@/store/auth'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { LoginForm } from '@/features/auth/components/login-form'

export function SessionExpiredDialog() {
  const sessionExpired = useAuthStore((state) => state.sessionExpired)
  const setSessionExpired = useAuthStore((state) => state.setSessionExpired)
  const clearAuth = useAuthStore((state) => state.clearAuth)

  return (
    <Dialog open={sessionExpired} onOpenChange={setSessionExpired}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Session Expired</DialogTitle>
        </DialogHeader>
        <div className="py-4">
          <p className="text-sm text-slate-600 mb-4">
            Your session has expired. Please log in again to continue.
          </p>
          <LoginForm onSuccess={() => {}} />
          <Button
            variant="ghost"
            onClick={() => clearAuth()}
            className="w-full mt-2"
          >
            Log Out
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
