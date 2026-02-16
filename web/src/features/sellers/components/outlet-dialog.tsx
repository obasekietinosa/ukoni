import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createOutlet, updateOutlet } from '../api'
import type { Outlet } from '../types'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

interface Props {
  sellerId: string
  outlet?: Outlet
  open: boolean
  onOpenChange: (open: boolean) => void
}

function OutletForm({
  sellerId,
  outlet,
  onSuccess,
  onCancel,
}: {
  sellerId: string
  outlet?: Outlet
  onSuccess: () => void
  onCancel: () => void
}) {
  const [name, setName] = useState(outlet?.name || '')
  const [channel, setChannel] = useState<Outlet['channel']>(
    outlet?.channel || 'physical'
  )
  const [address, setAddress] = useState(outlet?.address || '')
  const [websiteUrl, setWebsiteUrl] = useState(outlet?.website_url || '')
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: (data: {
      name: string
      channel: Outlet['channel']
      address?: string
      website_url?: string
    }) => {
      if (outlet) {
        return updateOutlet(outlet.id, data)
      }
      return createOutlet(sellerId, data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['outlets', sellerId] })
      onSuccess()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    mutation.mutate({
      name,
      channel,
      address: address || undefined,
      website_url: websiteUrl || undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <label htmlFor="name" className="text-sm font-medium">
          Name
        </label>
        <Input
          id="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Outlet Name"
          required
        />
      </div>
      <div className="space-y-2">
        <label htmlFor="channel" className="text-sm font-medium">
          Channel
        </label>
        <select
          id="channel"
          value={channel}
          onChange={(e) => setChannel(e.target.value as Outlet['channel'])}
          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-base md:text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        >
          <option value="physical">Physical</option>
          <option value="online">Online</option>
        </select>
      </div>
      <div className="space-y-2">
        <label htmlFor="address" className="text-sm font-medium">
          Address (Optional)
        </label>
        <Input
          id="address"
          value={address}
          onChange={(e) => setAddress(e.target.value)}
          placeholder="123 Main St"
        />
      </div>
      <div className="space-y-2">
        <label htmlFor="websiteUrl" className="text-sm font-medium">
          Website URL (Optional)
        </label>
        <Input
          id="websiteUrl"
          type="url"
          value={websiteUrl}
          onChange={(e) => setWebsiteUrl(e.target.value)}
          placeholder="https://example.com"
        />
      </div>
      <div className="flex justify-end gap-2">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={mutation.isPending}>
          {mutation.isPending ? 'Saving...' : 'Save'}
        </Button>
      </div>
    </form>
  )
}

export function OutletDialog({ sellerId, outlet, open, onOpenChange }: Props) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{outlet ? 'Edit Outlet' : 'Add Outlet'}</DialogTitle>
        </DialogHeader>
        {open && (
          <OutletForm
            key={outlet?.id || 'new'}
            sellerId={sellerId}
            outlet={outlet}
            onSuccess={() => onOpenChange(false)}
            onCancel={() => onOpenChange(false)}
          />
        )}
      </DialogContent>
    </Dialog>
  )
}
