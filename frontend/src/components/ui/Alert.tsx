import type { ReactNode } from 'react'

const styles = {
  success: 'bg-green-50 border-green-200 text-green-700',
  error: 'bg-red-50 border-red-200 text-red-600',
}

export const Alert = ({ type, children }: { type: 'success' | 'error'; children: ReactNode }) => (
  <div className={`text-sm px-3.5 py-2.5 rounded-lg border ${styles[type]}`}>{children}</div>
)
