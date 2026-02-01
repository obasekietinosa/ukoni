import { useState, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createSeller, updateSeller } from '../api'
import type { Seller } from '../types'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

interface Props {
  seller?: Seller // If provided, edit mode. If null, create mode.
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SellerDialog({ seller, open, onOpenChange }: Props) {
  const [name, setName] = useState('')
  const [type, setType] = useState<Seller['type']>('chain')
  const queryClient = useQueryClient()

  useEffect(() => {
    if (open) {
      if (seller) {
        setName(seller.name)
        setType(seller.type)
      } else {
        setName('')
        setType('chain')
      }
    }
  }, [open, seller])

  const mutation = useMutation({
    mutationFn: (data: { name: string; type: Seller['type'] }) => {
      if (seller) {
        return updateSeller(seller.id, data)
      }
      return createSeller(data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sellers'] })
      if (seller) {
        queryClient.invalidateQueries({ queryKey: ['seller', seller.id] })
      }
      onOpenChange(false)
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    mutation.mutate({ name, type })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{seller ? 'Edit Seller' : 'Add Seller'}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="name" className="text-sm font-medium">
              Name
            </label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Seller Name"
              required
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="type" className="text-sm font-medium">
              Type
            </label>
            <select
              id="type"
              value={type}
              onChange={(e) => setType(e.target.value as Seller['type'])}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="chain">Chain</option>
              <option value="independent">Independent</option>
              <option value="online">Online</option>
            </select>
          </div>
          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? 'Saving...' : 'Save'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
