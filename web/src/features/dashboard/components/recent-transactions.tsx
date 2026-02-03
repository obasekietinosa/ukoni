import { useQuery } from '@tanstack/react-query'
import { getTransactions } from '@/features/transactions/api'
import { useAllOutlets } from '@/features/sellers/hooks/use-all-outlets'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2 } from 'lucide-react'

interface RecentTransactionsProps {
  inventoryId: string
}

export function RecentTransactions({ inventoryId }: RecentTransactionsProps) {
  const { data: transactions, isLoading: isTransactionsLoading } = useQuery({
    queryKey: ['transactions', inventoryId, 5],
    queryFn: () => getTransactions(inventoryId, 5),
  })

  const { outlets, isLoading: isOutletsLoading } = useAllOutlets()

  const isLoading = isTransactionsLoading || isOutletsLoading

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Recent Transactions</CardTitle>
        </CardHeader>
        <CardContent className="flex justify-center p-6">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    )
  }

  if (!transactions || transactions.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Recent Transactions</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            No recent transactions.
          </p>
        </CardContent>
      </Card>
    )
  }

  const getOutletName = (outletId?: string) => {
    if (!outletId) return 'Unknown Outlet'
    const outlet = outlets.find((o) => o.id === outletId)
    return outlet ? outlet.name : 'Unknown Outlet'
  }

  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle>Recent Transactions</CardTitle>
      </CardHeader>
      <CardContent>
        <ul className="space-y-4">
          {transactions.map((tx) => (
            <li key={tx.id} className="flex flex-col space-y-1">
              <div className="flex justify-between font-medium">
                <span>{getOutletName(tx.outlet_id)}</span>
                <span className="text-muted-foreground">
                  {tx.items ? tx.items.length : 0} items
                </span>
              </div>
              <div className="text-xs text-muted-foreground">
                {new Date(tx.transaction_date).toLocaleDateString()}
              </div>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  )
}
