import { InventoryList } from '../components/inventory-list'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'

export function InventoryPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Inventory</h1>
          <p className="text-gray-500">Manage your current stock.</p>
        </div>
        <Link to="/products">
          <Button>Browse Catalog</Button>
        </Link>
      </div>

      <InventoryList />
    </div>
  )
}
