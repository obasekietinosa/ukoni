import * as React from 'react'
import { cn } from '@/lib/utils'
import { Button } from './button'

const DialogContext = React.createContext<{
  open: boolean
  onOpenChange: (open: boolean) => void
} | null>(null)

export function DialogTrigger({
  asChild,
  children,
}: {
  asChild?: boolean
  children: React.ReactNode
}) {
  const context = React.useContext(DialogContext)
  if (!context) throw new Error('DialogTrigger must be used within Dialog')
  const { onOpenChange } = context

  if (asChild && React.isValidElement(children)) {
    const child = children as React.ReactElement<{
      onClick?: React.MouseEventHandler
    }>
    return React.cloneElement(child, {
      onClick: (e: React.MouseEvent) => {
        child.props.onClick?.(e)
        onOpenChange(true)
      },
    })
  }

  return <button onClick={() => onOpenChange(true)}>{children}</button>
}

export function Dialog({
  open,
  onOpenChange,
  children,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  children: React.ReactNode
}) {
  return (
    <DialogContext.Provider value={{ open, onOpenChange }}>
      {children}
    </DialogContext.Provider>
  )
}

export function DialogContent({
  children,
  className,
}: {
  children: React.ReactNode
  className?: string
}) {
  const context = React.useContext(DialogContext)
  if (!context) throw new Error('DialogContent must be used within Dialog')
  const { open, onOpenChange } = context

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
        'backdrop:bg-slate-900/20 backdrop:backdrop-blur-sm p-0 rounded-2xl shadow-soft-xl w-[calc(100%-2rem)] max-w-lg m-auto border border-soft-pebble bg-white/95',
        'open:animate-in open:fade-in open:zoom-in-95 backdrop:animate-in backdrop:fade-in duration-300 ease-gentle',
        className
      )}
      onClose={() => onOpenChange(false)}
      onClick={(e) => {
        if (e.target === dialogRef.current) {
          onOpenChange(false)
        }
      }}
    >
      <div className="p-6 space-y-4 relative">
        {children}
        <Button
          variant="ghost"
          size="sm"
          className="absolute right-4 top-4 h-8 w-8 p-0 opacity-70 hover:opacity-100 rounded-full"
          onClick={() => onOpenChange(false)}
        >
          ✕<span className="sr-only">Close</span>
        </Button>
      </div>
    </dialog>
  )
}

export function DialogHeader({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'flex flex-col space-y-1.5 text-center sm:text-left',
        className
      )}
      {...props}
    />
  )
}

export function DialogFooter({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2',
        className
      )}
      {...props}
    />
  )
}

export function DialogTitle({
  className,
  ...props
}: React.HTMLAttributes<HTMLHeadingElement>) {
  return (
    <h2
      className={cn(
        'text-xl font-bold leading-none tracking-tight text-midnight-slate',
        className
      )}
      {...props}
    />
  )
}
