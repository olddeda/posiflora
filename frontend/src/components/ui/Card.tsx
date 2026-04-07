import type { ReactNode } from 'react'

export const Card = ({ children }: { children: ReactNode }) => (
  <div className="bg-white border border-gray-200 rounded-xl p-6">{children}</div>
)
