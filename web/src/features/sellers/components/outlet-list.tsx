import { useState } from 'react'
import type { Outlet } from '../types'
import { Button } from '@/components/ui/button'
import { OutletDialog } from './outlet-dialog'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { deleteOutlet } from '../api'
import { Edit2, Trash2, Plus, MapPin, Globe } from 'lucide-react'

interface Props {
  sellerId: string
  outlets: Outlet[]
}

export function OutletList({ sellerId, outlets }: Props) {
  const [editingOutlet, setEditingOutlet] = useState<Outlet | undefined>()
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: deleteOutlet,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['outlets', sellerId] })
    },
  })

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this outlet?')) {
      deleteMutation.mutate(id)
    }
  }

  const handleEdit = (outlet: Outlet) => {
    setEditingOutlet(outlet)
    setIsDialogOpen(true)
  }

  const handleCreate = () => {
    setEditingOutlet(undefined)
    setIsDialogOpen(true)
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-xl font-semibold">Outlets</h3>
        <Button size="sm" onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          Add Outlet
        </Button>
      </div>

      <div className="border rounded-md">
        <table className="w-full text-sm text-left">
          <thead className="bg-gray-50 text-gray-700 font-medium">
            <tr>
              <th className="p-4">Name</th>
              <th className="p-4">Channel</th>
              <th className="p-4">Details</th>
              <th className="p-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {outlets.map((outlet) => (
              <tr key={outlet.id} className="hover:bg-gray-50">
                <td className="p-4 font-medium">{outlet.name}</td>
                <td className="p-4 capitalize">{outlet.channel}</td>
                <td className="p-4">
                  <div className="flex flex-col gap-1 text-xs text-gray-500">
                    {outlet.address && (
                      <div className="flex items-center gap-1">
                        <MapPin className="w-3 h-3" />
                        {outlet.address}
                      </div>
                    )}
                    {outlet.website_url && (
                      <div className="flex items-center gap-1">
                        <Globe className="w-3 h-3" />
                        <a
                          href={outlet.website_url}
                          target="_blank"
                          rel="noreferrer"
                          className="hover:underline text-blue-500"
                        >
                          Website
                        </a>
                      </div>
                    )}
                  </div>
                </td>
                <td className="p-4 text-right space-x-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(outlet)}
                  >
                    <Edit2 className="w-4 h-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-red-500 hover:text-red-600"
                    onClick={() => handleDelete(outlet.id)}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </td>
              </tr>
            ))}
            {outlets.length === 0 && (
              <tr>
                <td colSpan={4} className="p-8 text-center text-gray-500">
                  No outlets found for this seller.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <OutletDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        sellerId={sellerId}
        outlet={editingOutlet}
      />
    </div>
  )
}
