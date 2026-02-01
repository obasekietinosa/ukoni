import { useQuery } from '@tanstack/react-query'
import { getSellers } from '../api'
import { SellerList } from '../components/seller-list'
import { Loader2 } from 'lucide-react'

export function SellersPage() {
  const { data: sellers, isLoading, error } = useQuery({
    queryKey: ['sellers'],
    queryFn: getSellers,
  })

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-500" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-red-500">
        Error loading sellers: {error instanceof Error ? error.message : 'Unknown error'}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Sellers</h1>
        <p className="text-muted-foreground">
          Manage sellers and their outlets.
        </p>
      </div>

      <SellerList sellers={sellers || []} />
    </div>
  )
}
