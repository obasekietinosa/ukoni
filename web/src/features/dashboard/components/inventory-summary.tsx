import { useQuery } from '@tanstack/react-query'
import { getInventoryProducts } from '@/features/inventory/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2 } from 'lucide-react'

interface InventorySummaryProps {
  inventoryId: string
}

export function InventorySummary({ inventoryId }: InventorySummaryProps) {
  const { data: products, isLoading } = useQuery({
    queryKey: ['inventory-products', inventoryId],
    queryFn: () => getInventoryProducts(inventoryId),
  })

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Inventory Summary</CardTitle>
        </CardHeader>
        <CardContent className="flex justify-center p-6">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    )
  }

  if (!products) return null

  const totalItems = products.length
  const totalQuantity = products.reduce((acc, p) => acc + p.quantity, 0)
  const lowStockItems = products
    .filter((p) => p.quantity <= 2)
    .sort((a, b) => a.quantity - b.quantity)
    .slice(0, 5)

  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle>Inventory Summary</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1">
            <p className="text-sm font-medium text-muted-foreground">
              Total Variants
            </p>
            <p className="text-2xl font-bold">{totalItems}</p>
          </div>
          <div className="space-y-1">
            <p className="text-sm font-medium text-muted-foreground">
              Total Quantity
            </p>
            <p className="text-2xl font-bold">{Math.round(totalQuantity)}</p>
          </div>
        </div>

        {lowStockItems.length > 0 && (
          <div className="space-y-2">
            <h4 className="text-sm font-semibold">Low Stock</h4>
            <ul className="space-y-2">
              {lowStockItems.map((item) => (
                <li
                  key={item.id}
                  className="flex items-center justify-between text-sm"
                >
                  <span className="truncate pr-2">
                    {item.product_name}
                    {item.variant_name && ` - ${item.variant_name}`}
                  </span>
                  <span className="font-medium text-destructive">
                    {item.quantity} {item.unit || item.product_unit}
                  </span>
                </li>
              ))}
            </ul>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
