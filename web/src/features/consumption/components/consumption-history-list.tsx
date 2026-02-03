import type { ConsumptionEventDetail } from '../types'

interface Props {
  events: ConsumptionEventDetail[]
  isLoading?: boolean
}

export function ConsumptionHistoryList({ events, isLoading }: Props) {
  if (isLoading) {
    return (
      <div className="p-4 text-center text-gray-500">Loading history...</div>
    )
  }

  if (events.length === 0) {
    return (
      <div className="p-8 text-center text-gray-500 border rounded-lg bg-gray-50">
        No consumption history found.
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-lg border bg-white shadow-sm">
      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm">
          <thead className="border-b bg-gray-50 font-medium text-gray-700">
            <tr>
              <th className="px-4 py-3">Date</th>
              <th className="px-4 py-3">Product</th>
              <th className="px-4 py-3">Variant / Brand</th>
              <th className="px-4 py-3 text-right">Quantity</th>
              <th className="px-4 py-3">Source</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {events.map((event) => (
              <tr key={event.id} className="transition-colors hover:bg-gray-50">
                <td className="px-4 py-3 font-medium text-gray-900">
                  {new Date(event.consumed_at).toLocaleString(undefined, {
                    month: 'short',
                    day: 'numeric',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </td>
                <td className="px-4 py-3 text-gray-900">
                  {event.canonical_product_name || 'Unknown Product'}
                </td>
                <td className="px-4 py-3 text-gray-600">
                  {event.variant_name ? (
                    <div className="flex flex-col">
                      <span>{event.variant_name}</span>
                      <span className="text-xs text-gray-500">
                        {event.brand}{' '}
                        {event.size
                          ? `- ${event.size}${event.product_unit}`
                          : ''}
                      </span>
                    </div>
                  ) : (
                    <span className="text-gray-400">-</span>
                  )}
                </td>
                <td className="px-4 py-3 text-right text-gray-900">
                  {event.quantity} {event.unit}
                </td>
                <td className="px-4 py-3 text-gray-600 capitalize">
                  {event.source}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
