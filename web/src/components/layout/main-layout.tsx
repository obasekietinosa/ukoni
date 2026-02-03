import { Link, Outlet } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { useQuery } from '@tanstack/react-query'
import { getInventory } from '@/features/inventory/api'
import { useCurrentUserRole } from '@/features/inventory/hooks/use-current-user-role'

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

  const { role } = useCurrentUserRole(activeInventoryId)

  const handleLogout = () => {
    clearInventory()
    clearAuth()
    window.location.reload()
  }

  return (
    <div className="flex h-screen w-full flex-col">
      <header className="flex items-center justify-between border-b px-6 py-4">
        <div className="flex items-center gap-4">
          <Link to="/" className="text-xl font-bold">
            Ukoni
          </Link>
          {inventory && (
            <div className="flex items-center gap-2">
              <span className="rounded bg-gray-100 px-2 py-1 text-sm text-gray-700">
                {inventory.name}
              </span>
              {role && (
                <span className="rounded border border-gray-200 px-2 py-1 text-xs uppercase text-gray-500">
                  {role}
                </span>
              )}
            </div>
          )}
          <nav className="ml-4 flex items-center gap-4 border-l pl-4">
            <Link
              to="/inventory"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              Inventory
            </Link>
            <Link
              to="/products"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              Products
            </Link>
            <Link
              to="/sellers"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              Sellers
            </Link>
            <Link
              to="/shopping-lists"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              Lists
            </Link>
            <Link
              to="/transactions"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              History
            </Link>
            <Link
              to="/consumption"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              Consumption
            </Link>
            <Link
              to="/members"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              Members
            </Link>
          </nav>
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
