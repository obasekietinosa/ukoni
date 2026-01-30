import { useQuery } from '@tanstack/react-query'
import { getInventories } from '../api'
import { CreateInventoryForm } from '../components/create-inventory-form'
import { useInventoryStore } from '@/store/inventory'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'

export function InventorySelectionRoute() {
  const {
    data: inventories,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['inventories'],
    queryFn: getInventories,
  })
  const setActiveInventoryId = useInventoryStore(
    (state) => state.setActiveInventoryId
  )
  const navigate = useNavigate()

  const handleSelect = (id: string) => {
    setActiveInventoryId(id)
    navigate('/')
  }

  return (
    <div className="flex h-screen w-full items-center justify-center bg-gray-50 p-4">
      <div className="w-full max-w-2xl space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Select Household</h1>
          <p className="text-gray-500">
            Choose a household to manage or create a new one.
          </p>
        </div>

        <div className="grid gap-8 md:grid-cols-2">
          <div className="space-y-4 rounded-lg border bg-white p-6 shadow-sm">
            <h2 className="text-xl font-semibold">Your Households</h2>
            {isLoading && <div>Loading...</div>}
            {error && (
              <div className="text-red-500">Failed to load households</div>
            )}

            <div className="space-y-2">
              {inventories?.length === 0 && (
                <p className="text-sm text-gray-500">No households found.</p>
              )}
              {inventories?.map((inventory) => (
                <Button
                  key={inventory.id}
                  variant="outline"
                  className="w-full justify-start"
                  onClick={() => handleSelect(inventory.id)}
                >
                  {inventory.name}
                </Button>
              ))}
            </div>
          </div>

          <div className="space-y-4 rounded-lg border bg-white p-6 shadow-sm">
            <h2 className="text-xl font-semibold">Create New</h2>
            <CreateInventoryForm onSuccess={() => navigate('/')} />
          </div>
        </div>
      </div>
    </div>
  )
}
