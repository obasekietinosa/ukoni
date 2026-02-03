import { useState } from 'react'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Logo } from '@/components/ui/logo'
import { useAuthStore } from '@/store/auth'
import { useInventoryStore } from '@/store/inventory'
import { useQuery } from '@tanstack/react-query'
import { getInventory } from '@/features/inventory/api'
import { useCurrentUserRole } from '@/features/inventory/hooks/use-current-user-role'
import { Menu, X } from 'lucide-react'

const NavLink = ({
  to,
  children,
  onClick,
}: {
  to: string
  children: React.ReactNode
  onClick?: () => void
}) => (
  <Link
    to={to}
    onClick={onClick}
    className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-midnight-slate"
  >
    {children}
  </Link>
)

const MobileNavLink = ({
  to,
  children,
  onClose,
}: {
  to: string
  children: React.ReactNode
  onClose: () => void
}) => (
  <Link
    to={to}
    onClick={onClose}
    className="block rounded-lg px-4 py-3 text-base font-medium text-slate-600 hover:bg-slate-100 hover:text-midnight-slate"
  >
    {children}
  </Link>
)

export function MainLayout() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const location = useLocation()

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

  // Close mobile menu when route changes
  if (mobileMenuOpen && location.state) {
    setMobileMenuOpen(false)
  }

  return (
    <div className="flex min-h-screen w-full flex-col bg-warm-chalk">
      <header className="sticky top-0 z-50 border-b border-soft-pebble bg-white/80 px-4 py-3 backdrop-blur-md md:px-6 md:py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3 md:gap-6">
            <Link to="/">
              <Logo className="h-8" iconClassName="h-8" />
            </Link>
            {inventory && (
              <div className="flex items-center gap-2 md:gap-3">
                <span className="max-w-[100px] truncate rounded-full bg-slate-100 px-2 py-1 text-xs font-medium text-slate-700 md:max-w-none md:px-3 md:text-sm">
                  {inventory.name}
                </span>
                {role && (
                  <span className="hidden rounded-full border border-soft-pebble px-3 py-1 text-xs uppercase tracking-wide text-slate-500 md:inline-block">
                    {role}
                  </span>
                )}
              </div>
            )}

            {/* Desktop Nav */}
            <nav className="hidden ml-2 items-center gap-1 md:flex">
              <NavLink to="/inventory">Inventory</NavLink>
              <NavLink to="/products">Products</NavLink>
              <NavLink to="/sellers">Sellers</NavLink>
              <NavLink to="/shopping-lists">Lists</NavLink>
              <NavLink to="/transactions">History</NavLink>
              <NavLink to="/consumption">Consumption</NavLink>
              <NavLink to="/members">Members</NavLink>
            </nav>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              className="hidden md:inline-flex"
              onClick={handleLogout}
            >
              Logout
            </Button>

            {/* Mobile Menu Toggle */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? (
                <X className="h-5 w-5" />
              ) : (
                <Menu className="h-5 w-5" />
              )}
            </Button>
          </div>
        </div>

        {/* Mobile Nav Overlay */}
        {mobileMenuOpen && (
          <div className="absolute left-0 top-full h-[calc(100vh-60px)] w-full overflow-y-auto border-b border-soft-pebble bg-white/95 px-4 py-6 backdrop-blur-md animate-in slide-in-from-top-2 md:hidden">
            <nav className="flex flex-col space-y-2">
              <MobileNavLink
                to="/inventory"
                onClose={() => setMobileMenuOpen(false)}
              >
                Inventory
              </MobileNavLink>
              <MobileNavLink
                to="/products"
                onClose={() => setMobileMenuOpen(false)}
              >
                Products
              </MobileNavLink>
              <MobileNavLink
                to="/sellers"
                onClose={() => setMobileMenuOpen(false)}
              >
                Sellers
              </MobileNavLink>
              <MobileNavLink
                to="/shopping-lists"
                onClose={() => setMobileMenuOpen(false)}
              >
                Lists
              </MobileNavLink>
              <MobileNavLink
                to="/transactions"
                onClose={() => setMobileMenuOpen(false)}
              >
                History
              </MobileNavLink>
              <MobileNavLink
                to="/consumption"
                onClose={() => setMobileMenuOpen(false)}
              >
                Consumption
              </MobileNavLink>
              <MobileNavLink
                to="/members"
                onClose={() => setMobileMenuOpen(false)}
              >
                Members
              </MobileNavLink>
              <div className="my-4 border-t border-soft-pebble pt-4">
                <Button
                  variant="ghost"
                  className="w-full justify-start"
                  onClick={handleLogout}
                >
                  Logout
                </Button>
              </div>
            </nav>
          </div>
        )}
      </header>

      <main className="flex-1 overflow-auto p-4 md:p-8">
        <div className="mx-auto max-w-7xl animate-in fade-in slide-in-from-bottom-4 duration-500 ease-gentle">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
