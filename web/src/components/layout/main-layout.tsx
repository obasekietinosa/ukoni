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
    <div className="flex min-h-screen w-full flex-col bg-warm-chalk">
      <header className="sticky top-0 z-50 flex items-center justify-between border-b border-soft-pebble bg-white/80 px-6 py-4 backdrop-blur-md">
        <div className="flex items-center gap-6">
          <Link
            to="/"
            className="text-xl font-bold tracking-tight text-midnight-slate"
          >
            Ukoni
          </Link>
          {inventory && (
            <div className="flex items-center gap-3">
              <span className="rounded-full bg-slate-100 px-3 py-1 text-sm font-medium text-slate-700">
                {inventory.name}
              </span>
              {role && (
                <span className="rounded-full border border-soft-pebble px-3 py-1 text-xs uppercase tracking-wide text-slate-500">
                  {role}
                </span>
              )}
            </div>
          )}
          <nav className="ml-2 flex items-center gap-1">
            <Link
              to="/inventory"
              className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
            >
              Inventory
            </Link>
            <Link
              to="/products"
              className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
            >
              Products
            </Link>
            <Link
              to="/sellers"
              className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
            >
              Sellers
            </Link>
            <Link
              to="/shopping-lists"
              className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
            >
              Lists
            </Link>
            <Link
              to="/transactions"
              className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
            >
              History
            </Link>
            <Link
              to="/consumption"
              className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
            >
              Consumption
            </Link>
            <Link
              to="/members"
              className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
            >
              Members
            </Link>
          </nav>
        </div>
        <Button variant="ghost" onClick={handleLogout}>
          Logout
        </Button>
      </header>
      <main className="flex-1 overflow-auto p-8">
        <div className="mx-auto max-w-7xl animate-in fade-in slide-in-from-bottom-4 duration-500 ease-gentle">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
