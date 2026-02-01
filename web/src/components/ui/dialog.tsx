import * as React from 'react'
import { cn } from '@/lib/utils'
import { Button } from './button'

interface DialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  children: React.ReactNode
  title?: string
}

export function Dialog({ open, onOpenChange, children, title }: DialogProps) {
  const dialogRef = React.useRef<HTMLDialogElement>(null)

  React.useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    if (open && !dialog.open) {
      dialog.showModal()
    } else if (!open && dialog.open) {
      dialog.close()
    }
  }, [open])

  return (
    <dialog
      ref={dialogRef}
      className={cn(
        'backdrop:bg-black/50 p-0 rounded-lg shadow-lg min-w-[400px]',
        'open:animate-in open:fade-in open:zoom-in-95 backdrop:animate-in backdrop:fade-in'
      )}
      onClose={() => onOpenChange(false)}
      onClick={(e) => {
        if (e.target === dialogRef.current) {
          onOpenChange(false)
        }
      }}
    >
      <div className="p-6 space-y-4">
        <div className="flex items-center justify-between">
          {title && <h2 className="text-lg font-semibold">{title}</h2>}
          <Button
            variant="ghost"
            size="sm"
            className="h-6 w-6 p-0"
            onClick={() => onOpenChange(false)}
          >
            ✕
          </Button>
        </div>
        {children}
      </div>
    </dialog>
  )
}
