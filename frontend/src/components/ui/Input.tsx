import { forwardRef } from 'react'
import { cn } from '../../utils/cn'

export const Input = forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => (
    <input
      ref={ref}
      {...props}
      className={cn(
        'w-full px-3 py-2 text-sm border border-gray-300 rounded-lg outline-none focus:border-blue-500 focus-visible:ring-2 focus-visible:ring-blue-500/30 transition-colors',
        className,
      )}
    />
  ),
)

Input.displayName = 'Input'
