import type { ActionResult } from '../api/chat'
import { Card } from '@/components/ui/card'
import { Check, Plus, ClipboardList, Calendar } from 'lucide-react'

// Define explicit types based on backend models
type CanonicalProductResult = {
  name: string
  status: 'success' | 'error'
  error?: string
  data?: {
    id: string
    name: string
    description?: string
  }
}

type ShoppingListItem = {
  id: string
  product?: { name: string }
  product_variant?: { variant_name: string }
  canonical_product?: { name: string }
  quantity: number
  unit: string
  notes?: string
}

type Plan = {
  id: string
  title: string
  description?: string
}

export function ActionRenderer({ actions }: { actions: ActionResult[] }) {
  if (!actions || actions.length === 0) return null

  return (
    <div className="flex flex-col gap-2 mt-2 w-full">
      {actions.map((action, i) => (
        <ActionItem key={i} action={action} />
      ))}
    </div>
  )
}

function ActionItem({ action }: { action: ActionResult }) {
  const { tool_name, data } = action

  if (tool_name === 'add_canonical_products') {
    const results = data as CanonicalProductResult[]
    if (!Array.isArray(results)) return null // Should be array

    // Filter successes
    const added = results.filter((r) => r.status === 'success')
    if (added.length === 0) return null

    return (
      <Card className="p-3 bg-green-50 border-green-200 text-sm">
        <div className="flex items-center gap-2 mb-1 font-medium text-green-800">
          <Check className="h-4 w-4" />
          <span>Added to Inventory</span>
        </div>
        <ul className="list-disc list-inside text-green-700 pl-1">
          {added.map((item, idx) => (
            <li key={idx}>{item.name}</li>
          ))}
        </ul>
      </Card>
    )
  }

  if (tool_name === 'add_shopping_list_item') {
    const item = data as ShoppingListItem
    const name =
      item.product?.name ||
      item.product_variant?.variant_name ||
      item.canonical_product?.name ||
      'Item'

    return (
      <Card className="p-3 bg-blue-50 border-blue-200 text-sm">
        <div className="flex items-center gap-2 mb-1 font-medium text-blue-800">
          <Plus className="h-4 w-4" />
          <span>Added to Shopping List</span>
        </div>
        <div className="text-blue-700 pl-6">
          {name} ({item.quantity} {item.unit})
        </div>
      </Card>
    )
  }

  if (tool_name === 'create_plan') {
    const plan = data as Plan
    return (
      <Card className="p-3 bg-purple-50 border-purple-200 text-sm">
        <div className="flex items-center gap-2 mb-1 font-medium text-purple-800">
          <Calendar className="h-4 w-4" />
          <span>Plan Created</span>
        </div>
        <div className="text-purple-700 pl-6">{plan.title}</div>
      </Card>
    )
  }

  if (tool_name === 'create_shopping_list') {
    const list = data as { name: string }
    return (
      <Card className="p-3 bg-orange-50 border-orange-200 text-sm">
        <div className="flex items-center gap-2 mb-1 font-medium text-orange-800">
          <ClipboardList className="h-4 w-4" />
          <span>List Created</span>
        </div>
        <div className="text-orange-700 pl-6">{list.name}</div>
      </Card>
    )
  }

  return null
}
