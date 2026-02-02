import { useQuery, useQueries } from '@tanstack/react-query'
import { getSellers, getOutlets } from '../api'
import type { Outlet } from '../types'

export type OutletWithSeller = Outlet & {
  sellerName: string
}

export function useAllOutlets() {
  const { data: sellers, isLoading: isSellersLoading } = useQuery({
    queryKey: ['sellers'],
    queryFn: getSellers,
  })

  const outletQueries = useQueries({
    queries: (sellers || []).map((seller) => ({
      queryKey: ['outlets', seller.id],
      queryFn: () => getOutlets(seller.id),
      enabled: !!seller.id,
    })),
  })

  const isLoadingOutlets = outletQueries.some((q) => q.isLoading)
  const isError = outletQueries.some((q) => q.isError)

  const outlets: OutletWithSeller[] = []

  if (sellers && !isLoadingOutlets) {
    outletQueries.forEach((query, index) => {
      if (query.data) {
        const seller = sellers[index]
        query.data.forEach((outlet) => {
          outlets.push({
            ...outlet,
            sellerName: seller.name,
          })
        })
      }
    })
  }

  return {
    outlets,
    isLoading: isSellersLoading || isLoadingOutlets,
    isError,
  }
}
