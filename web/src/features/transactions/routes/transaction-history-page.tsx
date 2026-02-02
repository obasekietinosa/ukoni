import { useQuery } from '@tanstack/react-query'
import { getTransactions } from '../api'
import { useInventoryStore } from '@/store/inventory'

export function TransactionHistoryPage() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )

  const {
    data: transactions,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['transactions', activeInventoryId],
    queryFn: () => getTransactions(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  if (isLoading) {
    return (
      <div className="p-8 text-center text-gray-500">Loading history...</div>
    )
  }

  if (error) {
    return (
      <div className="p-8 text-center text-red-500">
        Error loading history: {(error as Error).message}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">History</h1>
        <p className="text-gray-500">View past transactions and consumption.</p>
      </div>

      <div className="overflow-hidden rounded-lg border bg-white shadow-sm">
        <table className="w-full text-left text-sm">
          <thead className="border-b bg-gray-50 font-medium text-gray-700">
            <tr>
              <th className="px-4 py-3">Date</th>
              <th className="px-4 py-3">Type</th>
              <th className="px-4 py-3 text-right">Amount</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {!transactions || transactions.length === 0 ? (
              <tr>
                <td colSpan={3} className="px-4 py-8 text-center text-gray-500">
                  No transactions found.
                </td>
              </tr>
            ) : (
              transactions.map((t) => (
                <tr key={t.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 text-gray-900">
                    {new Date(t.transaction_date).toLocaleDateString()}
                  </td>
                  <td className="px-4 py-3 text-gray-600">
                    {t.total_amount && t.total_amount > 0
                      ? 'Purchase'
                      : 'Consumption / Adjustment'}
                  </td>
                  <td className="px-4 py-3 text-right font-medium text-gray-900">
                    {t.total_amount ? `$${t.total_amount.toFixed(2)}` : '-'}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
