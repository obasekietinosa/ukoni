import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Plus,
  ShoppingCart,
  Activity,
  Calendar,
  History,
  Package,
  Users,
} from 'lucide-react'

export function QuickActions() {
  const navigate = useNavigate()

  const actions = [
    {
      label: 'Manage Inventory',
      icon: Plus,
      onClick: () => navigate('/inventory'),
    },
    {
      label: 'Shopping Lists',
      icon: ShoppingCart,
      onClick: () => navigate('/shopping-lists'),
    },
    {
      label: 'Plans',
      icon: Calendar,
      onClick: () => navigate('/plans'),
    },
    {
      label: 'Record Consumption',
      icon: Activity,
      onClick: () => navigate('/consumption'),
    },
    {
      label: 'Transactions',
      icon: History,
      onClick: () => navigate('/transactions'),
    },
    {
      label: 'Products',
      icon: Package,
      onClick: () => navigate('/products'),
    },
    {
      label: 'Members',
      icon: Users,
      onClick: () => navigate('/members'),
    },
  ]

  return (
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
      </CardHeader>
      <CardContent className="grid grid-cols-1 gap-2 sm:grid-cols-2">
        {actions.map((action) => (
          <Button
            key={action.label}
            variant="outline"
            className="w-full justify-start"
            onClick={action.onClick}
          >
            <action.icon className="mr-2 h-4 w-4" />
            {action.label}
          </Button>
        ))}
      </CardContent>
    </Card>
  )
}
