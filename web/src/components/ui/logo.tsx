import type { SVGProps } from 'react'
import { cn } from '@/lib/utils'

export function LogoMark(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 165 170"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect x="0" y="0" width="65" height="170" rx="32.5" fill="#1E293B" />
      <rect
        x="50"
        y="0"
        width="65"
        height="170"
        rx="32.5"
        fill="#64D2B1"
        fillOpacity="0.8"
      />
      <rect
        x="100"
        y="0"
        width="65"
        height="170"
        rx="32.5"
        fill="#A7F3D0"
        fillOpacity="0.7"
      />
    </svg>
  )
}

interface LogoProps extends React.HTMLAttributes<HTMLDivElement> {
  iconClassName?: string
}

export function Logo({ className, iconClassName, ...props }: LogoProps) {
  return (
    <div className={cn('flex items-center gap-3', className)} {...props}>
      <LogoMark className={cn('h-10 w-auto', iconClassName)} />
      <span className="font-sans text-2xl font-bold tracking-tight text-midnight-slate">
        Ukoni
      </span>
    </div>
  )
}
