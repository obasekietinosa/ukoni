import * as React from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const buttonVariants = cva(
  'inline-flex items-center justify-center whitespace-nowrap rounded-full text-sm font-bold transition-all duration-300 ease-gentle focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-slate-950 disabled:pointer-events-none disabled:opacity-50 active:scale-95',
  {
    variants: {
      variant: {
        default:
          'bg-electric-mint text-midnight-slate shadow hover:bg-electric-mint/90 hover:-translate-y-0.5 hover:shadow-lg',
        destructive: 'bg-red-500 text-slate-50 shadow-sm hover:bg-red-500/90',
        outline:
          'border border-soft-pebble bg-white shadow-sm hover:bg-slate-100 hover:text-midnight-slate',
        secondary:
          'bg-slate-100 text-midnight-slate shadow-sm hover:bg-slate-100/80',
        ghost: 'hover:bg-slate-100 hover:text-midnight-slate',
        link: 'text-midnight-slate underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-9 px-4 py-2',
        sm: 'h-8 px-3 text-xs',
        lg: 'h-10 px-8',
        icon: 'h-9 w-9',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
)

export interface ButtonProps
  extends
    React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <button
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = 'Button'

export { Button }
