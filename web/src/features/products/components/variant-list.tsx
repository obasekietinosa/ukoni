import { useQuery } from '@tanstack/react-query'
import { getVariants } from '../api'

interface Props {
  productId: string
}

export function VariantList({ productId }: Props) {
  const {
    data: variants,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['variants', productId],
    queryFn: () => getVariants(productId),
  })

  if (isLoading)
    return <div className="text-xs text-gray-500">Loading variants...</div>
  if (error)
    return <div className="text-xs text-red-500">Failed to load variants</div>

  if (!variants || variants.length === 0) {
    return <div className="text-xs text-gray-500 italic">No variants yet.</div>
  }

  return (
    <div className="space-y-2 mt-2">
      {variants.map((variant) => (
        <div
          key={variant.id}
          className="flex items-center justify-between rounded bg-gray-50 p-2 text-sm"
        >
          <div>
            <span className="font-medium">{variant.variant_name}</span>
            {(variant.size || variant.unit) && (
              <span className="text-gray-500 ml-2">
                ({variant.size} {variant.unit})
              </span>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}
