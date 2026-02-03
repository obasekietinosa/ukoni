import { useQuery } from '@tanstack/react-query'
import { getConsumptionEvents } from '@/features/consumption/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2 } from 'lucide-react'

interface RecentConsumptionProps {
  inventoryId: string
}

export function RecentConsumption({ inventoryId }: RecentConsumptionProps) {
  const { data: events, isLoading } = useQuery({
    queryKey: ['consumption-events', inventoryId, 5],
    queryFn: () => getConsumptionEvents(inventoryId, { limit: 5 }),
  })

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Recent Consumption</CardTitle>
        </CardHeader>
        <CardContent className="flex justify-center p-6">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    )
  }

  if (!events || events.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Recent Consumption</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            No recent consumption events.
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle>Recent Consumption</CardTitle>
      </CardHeader>
      <CardContent>
        <ul className="space-y-4">
          {events.map((event) => (
            <li key={event.id} className="flex flex-col space-y-1">
              <div className="flex justify-between font-medium">
                <span className="truncate pr-2">
                  {event.product_name ||
                    event.canonical_product_name ||
                    'Unknown Product'}
                  {event.variant_name ? ` - ${event.variant_name}` : ''}
                </span>
                <span className="whitespace-nowrap text-muted-foreground">
                  {event.quantity} {event.unit || event.product_unit}
                </span>
              </div>
              <div className="text-xs text-muted-foreground">
                {new Date(event.consumed_at).toLocaleString()}
              </div>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  )
}
