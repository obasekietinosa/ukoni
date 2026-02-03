import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Plus, ShoppingCart, Activity } from 'lucide-react'

export function QuickActions() {
  const navigate = useNavigate()

  return (
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-2">
        <Button
          variant="outline"
          className="w-full justify-start"
          onClick={() => navigate('/inventory')}
        >
          <Plus className="mr-2 h-4 w-4" />
          Manage Inventory
        </Button>
        <Button
          variant="outline"
          className="w-full justify-start"
          onClick={() => navigate('/shopping-lists')}
        >
          <ShoppingCart className="mr-2 h-4 w-4" />
          Shopping Lists
        </Button>
        <Button
          variant="outline"
          className="w-full justify-start"
          onClick={() => navigate('/consumption')}
        >
          <Activity className="mr-2 h-4 w-4" />
          Record Consumption
        </Button>
      </CardContent>
    </Card>
  )
}
