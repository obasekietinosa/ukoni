import { useAuthStore } from '@/store/auth'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

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
          <p className="text-sm text-slate-600">
            Your session has expired. You can cancel to stay on this page to
            save your work, or log out immediately.
          </p>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setSessionExpired(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={() => clearAuth()}>
            Log Out
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
