import { useEffect, useRef } from 'react'
import { createPortal } from 'react-dom'
import { Button } from './button'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  children: React.ReactNode
  footer?: React.ReactNode
}

export function Modal({ isOpen, onClose, title, children, footer }: ModalProps) {
  const dialogRef = useRef<HTMLDialogElement>(null)

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    if (isOpen) {
        if (!dialog.open) {
            dialog.showModal()
        }
    } else {
        if (dialog.open) {
            dialog.close()
        }
    }
  }, [isOpen])

  // Handle ESC key and backdrop click via built-in dialog behavior?
  // <dialog> handles ESC by default.
  // Backdrop click needs manual handling or CSS.

  const handleClose = () => {
      onClose()
  }

  if (!isOpen) return null

  return createPortal(
    <dialog
      ref={dialogRef}
      className="backdrop:bg-black/50 backdrop:backdrop-blur-sm bg-white rounded-lg shadow-xl p-0 w-full max-w-md open:animate-in open:fade-in open:zoom-in-95"
      onClose={handleClose}
      onClick={(e) => {
          if (e.target === dialogRef.current) {
              handleClose()
          }
      }}
    >
      <div className="flex items-center justify-between p-4 border-b">
        <h2 className="text-lg font-semibold">{title}</h2>
        <Button variant="ghost" size="sm" onClick={handleClose}>
          X
        </Button>
      </div>
      <div className="p-4">{children}</div>
      {footer && <div className="p-4 border-t bg-gray-50 flex justify-end gap-2">{footer}</div>}
    </dialog>,
    document.body
  )
}
