import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useInventoryStore } from '@/store/inventory'

export function InventoryGuard() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const location = useLocation()

  if (!activeInventoryId) {
    return (
      <Navigate to="/select-inventory" state={{ from: location }} replace />
    )
  }

  return <Outlet />
}
