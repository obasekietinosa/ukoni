import { Link } from 'react-router-dom'
import type { Plan } from '../types'
import { Calendar, ChevronRight } from 'lucide-react'

interface Props {
  plans: Plan[]
}

export function PlanList({ plans }: Props) {
  if (plans.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-slate-300 bg-slate-50 p-8 text-center animate-in fade-in-50">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-white shadow-sm mb-4">
          <Calendar className="h-6 w-6 text-slate-400" />
        </div>
        <h3 className="mt-2 text-sm font-semibold text-slate-900">No plans yet</h3>
        <p className="mt-1 text-sm text-slate-500 max-w-sm">
          Create a plan for meals, cleaning, or any household activity.
        </p>
      </div>
    )
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {plans.map((plan) => (
        <Link
          key={plan.id}
          to={`/plans/${plan.id}`}
          className="group relative flex flex-col justify-between overflow-hidden rounded-xl border border-slate-200 bg-white p-6 shadow-sm transition-all hover:border-electric-mint/50 hover:shadow-md active:scale-[0.98]"
        >
          <div>
            <h3 className="font-semibold text-lg text-slate-900 group-hover:text-electric-mint transition-colors">
              {plan.title}
            </h3>
            {plan.description && (
              <p className="mt-2 line-clamp-2 text-sm text-slate-500 leading-relaxed">
                {plan.description}
              </p>
            )}
          </div>
          <div className="mt-6 flex items-center justify-between border-t border-slate-100 pt-4">
             <span className="text-xs font-medium text-slate-400 flex items-center gap-1">
                Updated {new Date(plan.updated_at).toLocaleDateString()}
             </span>
             <ChevronRight className="h-4 w-4 text-slate-300 group-hover:text-electric-mint transition-colors" />
          </div>
        </Link>
      ))}
    </div>
  )
}
