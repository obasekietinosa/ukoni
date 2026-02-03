import { useInventoryStore } from '@/store/inventory'
import { InventorySummary } from '../components/inventory-summary'
import { RecentConsumption } from '../components/recent-consumption'
import { RecentTransactions } from '../components/recent-transactions'
import { QuickActions } from '../components/quick-actions'

export function DashboardPage() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  if (!activeInventoryId) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-muted-foreground">Please select an inventory.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <div className="col-span-1 lg:col-span-2">
          <InventorySummary inventoryId={activeInventoryId} />
        </div>
        <div className="col-span-1 lg:col-span-2">
          <QuickActions />
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <div className="col-span-1">
          <RecentConsumption inventoryId={activeInventoryId} />
        </div>
        <div className="col-span-1">
          <RecentTransactions inventoryId={activeInventoryId} />
        </div>
      </div>
    </div>
  )
}
