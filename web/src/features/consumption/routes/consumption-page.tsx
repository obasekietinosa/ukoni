import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useInventoryStore } from '@/store/inventory'
import { getConsumptionEvents } from '../api'
import { ConsumptionHistoryList } from '../components/consumption-history-list'
import { RecordConsumptionDialog } from '../components/record-consumption-dialog'

export function ConsumptionPage() {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const [recordDialogOpen, setRecordDialogOpen] = useState(false)

  const {
    data: events,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ['consumption-events', activeInventoryId],
    queryFn: () => getConsumptionEvents(activeInventoryId!),
    enabled: !!activeInventoryId,
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Consumption</h1>
          <p className="text-gray-500">Track your product usage history.</p>
        </div>
        <Button onClick={() => setRecordDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Record Consumption
        </Button>
      </div>

      <ConsumptionHistoryList events={events || []} isLoading={isLoading} />

      <RecordConsumptionDialog
        open={recordDialogOpen}
        onOpenChange={setRecordDialogOpen}
        onSuccess={() => refetch()}
      />
    </div>
  )
}
