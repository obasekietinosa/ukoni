import { Outlet } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { useQuery } from '@tanstack/react-query'
import { getInventory } from '@/features/inventory/api'

export function MainLayout() {
  const clearAuth = useAuthStore((state) => state.clearAuth)
  const clearInventory = useInventoryStore(
    (state) => state.clearActiveInventory
  )
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const { data: inventory } = useQuery({
    queryKey: ['inventory', activeInventoryId],
    queryFn: () => getInventory(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  const handleLogout = () => {
    clearInventory()
    clearAuth()
    window.location.reload()
  }

  return (
    <div className="flex h-screen w-full flex-col">
      <header className="flex items-center justify-between border-b px-6 py-4">
        <div className="flex items-center gap-4">
          <h1 className="text-xl font-bold">Ukoni</h1>
          {inventory && (
            <span className="rounded bg-gray-100 px-2 py-1 text-sm text-gray-700">
              {inventory.name}
            </span>
          )}
        </div>
        <Button variant="ghost" onClick={handleLogout}>
          Logout
        </Button>
      </header>
      <main className="flex-1 overflow-auto p-6">
        <Outlet />
      </main>
    </div>
  )
}
