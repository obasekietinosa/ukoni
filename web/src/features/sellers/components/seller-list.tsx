import { useState } from 'react'
import { Link } from 'react-router-dom'
import type { Seller } from '../types'
import { Button } from '@/components/ui/button'
import { SellerDialog } from './seller-dialog'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { deleteSeller } from '../api'
import { Edit2, Trash2, Plus } from 'lucide-react'

interface Props {
  sellers: Seller[]
}

export function SellerList({ sellers }: Props) {
  const [editingSeller, setEditingSeller] = useState<Seller | undefined>()
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: deleteSeller,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sellers'] })
    },
  })

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this seller?')) {
      deleteMutation.mutate(id)
    }
  }

  const handleEdit = (seller: Seller) => {
    setEditingSeller(seller)
    setIsDialogOpen(true)
  }

  const handleCreate = () => {
    setEditingSeller(undefined)
    setIsDialogOpen(true)
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Sellers</h2>
        <Button onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          Add Seller
        </Button>
      </div>

      <div className="border rounded-md">
        <table className="w-full text-sm text-left">
          <thead className="bg-gray-50 text-gray-700 font-medium">
            <tr>
              <th className="p-4">Name</th>
              <th className="p-4">Type</th>
              <th className="p-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {sellers.map((seller) => (
              <tr key={seller.id} className="hover:bg-gray-50">
                <td className="p-4">
                  <Link
                    to={`/sellers/${seller.id}`}
                    className="font-medium text-blue-600 hover:underline"
                  >
                    {seller.name}
                  </Link>
                </td>
                <td className="p-4 capitalize">{seller.type}</td>
                <td className="p-4 text-right space-x-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(seller)}
                  >
                    <Edit2 className="w-4 h-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-red-500 hover:text-red-600"
                    onClick={() => handleDelete(seller.id)}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </td>
              </tr>
            ))}
            {sellers.length === 0 && (
              <tr>
                <td colSpan={3} className="p-8 text-center text-gray-500">
                  No sellers found. Add one to get started.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <SellerDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        seller={editingSeller}
      />
    </div>
  )
}
