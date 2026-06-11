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

type InventoryProductDetail = {
  id: string
  inventory_id: string
  product_variant_id: string
  quantity: number
  unit?: string
  created_at: string
  last_updated: string
  product_name: string
  brand?: string
  variant_name: string
  size?: number
  product_unit?: string
  canonical_product_id?: string
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

type PlanItem = {
  id: string
  product?: { name: string }
  product_variant?: { variant_name: string }
  canonical_product?: { name: string }
  quantity: number
  unit: string
  note?: string
}

type Plan = {
  id: string
  title: string
  description?: string
}

type ShoppingList = {
  id: string
  name: string
  inventory_id: string
}

type ConsumptionEvent = {
  id: string
  canonical_product_id?: string
  product_variant_id?: string
  canonical_product_name?: string
  product_name?: string
  variant_name?: string
  quantity?: number
  unit?: string
  consumed_at?: string
}

type RecordConsumptionResult = {
  event?: ConsumptionEvent
  inventory_adjusted: boolean
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

  if (tool_name === 'search_canonical_products') {
    const products = data as {
      id: string
      name: string
      description?: string
    }[]
    if (!Array.isArray(products) || products.length === 0) return null

    return (
      <Card className="p-3 bg-emerald-50 border-emerald-200 text-sm">
        <div className="flex items-center gap-2 mb-2 font-medium text-emerald-800">
          <ClipboardList className="h-4 w-4" />
          <span>Found Products</span>
        </div>
        <div className="flex flex-col gap-2">
          {products.map((product) => (
            <div
              key={product.id}
              className="flex flex-col justify-center bg-white p-2 rounded border border-emerald-100"
            >
              <span className="font-medium text-slate-800">{product.name}</span>
              {product.description && (
                <span className="text-xs text-slate-500 mt-1">
                  {product.description}
                </span>
              )}
            </div>
          ))}
        </div>
      </Card>
    )
  }

  if (tool_name === 'list_inventory_products') {
    const products = data as InventoryProductDetail[]
    if (!Array.isArray(products) || products.length === 0) return null

    return (
      <Card className="p-3 bg-teal-50 border-teal-200 text-sm">
        <div className="flex items-center gap-2 mb-2 font-medium text-teal-800">
          <ClipboardList className="h-4 w-4" />
          <span>Inventory Products</span>
        </div>
        <div className="flex flex-col gap-2">
          {products.map((product) => (
            <div
              key={product.id}
              className="flex justify-between items-center bg-white p-2 rounded border border-teal-100"
            >
              <div className="flex flex-col">
                <span className="font-medium text-slate-800">
                  {product.product_name}
                </span>
                <span className="text-xs text-slate-500">
                  {product.variant_name}{' '}
                  {product.size
                    ? `(${product.size} ${product.product_unit || ''})`
                    : ''}
                </span>
              </div>
              <div className="text-teal-700 font-semibold bg-teal-100 px-2 py-1 rounded">
                {product.quantity} {product.unit}
              </div>
            </div>
          ))}
        </div>
      </Card>
    )
  }

  if (tool_name === 'list_shopping_lists') {
    const lists = data as ShoppingList[]
    if (!Array.isArray(lists) || lists.length === 0) return null

    return (
      <Card className="p-3 bg-amber-50 border-amber-200 text-sm">
        <div className="flex items-center gap-2 mb-2 font-medium text-amber-800">
          <ClipboardList className="h-4 w-4" />
          <span>Shopping Lists</span>
        </div>
        <div className="flex flex-col gap-2">
          {lists.map((list) => (
            <div
              key={list.id}
              className="flex justify-between items-center bg-white p-2 rounded border border-amber-100"
            >
              <span className="font-medium text-slate-800">{list.name}</span>
            </div>
          ))}
        </div>
      </Card>
    )
  }

  if (tool_name === 'list_plans') {
    const plans = data as Plan[]
    if (!Array.isArray(plans) || plans.length === 0) return null

    return (
      <Card className="p-3 bg-fuchsia-50 border-fuchsia-200 text-sm">
        <div className="flex items-center gap-2 mb-2 font-medium text-fuchsia-800">
          <Calendar className="h-4 w-4" />
          <span>Plans</span>
        </div>
        <div className="flex flex-col gap-2">
          {plans.map((plan) => (
            <div
              key={plan.id}
              className="flex flex-col justify-center bg-white p-2 rounded border border-fuchsia-100"
            >
              <span className="font-medium text-slate-800">{plan.title}</span>
              {plan.description && (
                <span className="text-xs text-slate-500 mt-1">
                  {plan.description}
                </span>
              )}
            </div>
          ))}
        </div>
      </Card>
    )
  }

  if (tool_name === 'list_plan_groups') {
    const groups = data as Plan[]
    if (!Array.isArray(groups) || groups.length === 0) return null

    return (
      <Card className="p-3 bg-fuchsia-50 border-fuchsia-200 text-sm">
        <div className="flex items-center gap-2 mb-2 font-medium text-fuchsia-800">
          <Calendar className="h-4 w-4" />
          <span>Plan Groups</span>
        </div>
        <div className="flex flex-col gap-2">
          {groups.map((group) => (
            <div
              key={group.id}
              className="flex flex-col justify-center bg-white p-2 rounded border border-fuchsia-100"
            >
              <span className="font-medium text-slate-800">{group.title}</span>
              {group.description && (
                <span className="text-xs text-slate-500 mt-1">
                  {group.description}
                </span>
              )}
            </div>
          ))}
        </div>
      </Card>
    )
  }

  if (tool_name === 'list_consumption_events') {
    const events = data as ConsumptionEvent[]
    if (!Array.isArray(events) || events.length === 0) return null

    return (
      <Card className="p-3 bg-cyan-50 border-cyan-200 text-sm">
        <div className="flex items-center gap-2 mb-2 font-medium text-cyan-800">
          <ClipboardList className="h-4 w-4" />
          <span>Recent Consumption</span>
        </div>
        <div className="flex flex-col gap-2">
          {events.map((event) => {
            const name =
              event.product_name ||
              event.variant_name ||
              event.canonical_product_name ||
              'Item'

            return (
              <div
                key={event.id}
                className="flex justify-between items-center bg-white p-2 rounded border border-cyan-100"
              >
                <div className="flex flex-col min-w-0">
                  <span className="font-medium text-slate-800 truncate">
                    {name}
                  </span>
                  {event.consumed_at && (
                    <span className="text-xs text-slate-500">
                      {new Date(event.consumed_at).toLocaleString()}
                    </span>
                  )}
                </div>
                <div className="text-cyan-700 font-semibold bg-cyan-100 px-2 py-1 rounded">
                  {event.quantity ?? 0} {event.unit}
                </div>
              </div>
            )
          })}
        </div>
      </Card>
    )
  }

  if (tool_name === 'record_consumption') {
    const result = data as RecordConsumptionResult
    const event = result.event

    return (
      <Card className="p-3 bg-lime-50 border-lime-200 text-sm">
        <div className="flex items-center gap-2 mb-1 font-medium text-lime-800">
          <Check className="h-4 w-4" />
          <span>Consumption Recorded</span>
        </div>
        {event && (
          <div className="text-lime-700 pl-6">
            {event.quantity ?? 0} {event.unit}
            {result.inventory_adjusted ? ' removed from inventory' : ''}
          </div>
        )}
      </Card>
    )
  }

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

  if (tool_name === 'add_plan_item') {
    const item = data as PlanItem
    const name =
      item.product?.name ||
      item.product_variant?.variant_name ||
      item.canonical_product?.name ||
      'Item'

    return (
      <Card className="p-3 bg-indigo-50 border-indigo-200 text-sm">
        <div className="flex items-center gap-2 mb-1 font-medium text-indigo-800">
          <Plus className="h-4 w-4" />
          <span>Added to Plan</span>
        </div>
        <div className="text-indigo-700 pl-6">
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

  if (
    tool_name === 'create_shopping_list' ||
    tool_name === 'create_shopping_list_from_plan' ||
    tool_name === 'create_shopping_list_from_plan_group'
  ) {
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
