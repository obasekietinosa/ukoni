import { Link } from 'react-router-dom'
import type { ShoppingList } from '../types'
import { ChevronRight } from 'lucide-react'

interface Props {
  list: ShoppingList
}

export function ShoppingListCard({ list }: Props) {
  return (
    <Link to={`/shopping-lists/${list.id}`} className="block group">
      <div className="rounded-lg border p-4 hover:border-blue-500 hover:bg-blue-50 transition-all flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-lg group-hover:text-blue-700">
            {list.name}
          </h3>
          <p className="text-sm text-gray-500">
            Updated {new Date(list.last_updated_at).toLocaleDateString()}
          </p>
        </div>
        <ChevronRight className="h-5 w-5 text-gray-400 group-hover:text-blue-500" />
      </div>
    </Link>
  )
}
