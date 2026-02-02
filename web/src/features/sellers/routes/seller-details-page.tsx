import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getSeller, getOutlets } from '../api'
import { OutletList } from '../components/outlet-list'
import { Loader2, ArrowLeft } from 'lucide-react'

export function SellerDetailsPage() {
  const { id } = useParams<{ id: string }>()

  const {
    data: seller,
    isLoading: isLoadingSeller,
    error: sellerError,
  } = useQuery({
    queryKey: ['seller', id],
    queryFn: () => getSeller(id!),
    enabled: !!id,
  })

  const {
    data: outlets,
    isLoading: isLoadingOutlets,
    error: outletsError,
  } = useQuery({
    queryKey: ['outlets', id],
    queryFn: () => getOutlets(id!),
    enabled: !!id,
  })

  if (isLoadingSeller || isLoadingOutlets) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-500" />
      </div>
    )
  }

  if (sellerError || outletsError) {
    return <div className="text-red-500">Error loading seller details.</div>
  }

  if (!seller) {
    return <div className="text-red-500">Seller not found</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/sellers" className="rounded-full p-2 hover:bg-gray-100">
          <ArrowLeft className="h-6 w-6" />
        </Link>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{seller.name}</h1>
          <p className="text-muted-foreground capitalize">
            {seller.type} Seller
          </p>
        </div>
      </div>

      <div className="border-t pt-6">
        <OutletList sellerId={seller.id} outlets={outlets || []} />
      </div>
    </div>
  )
}
