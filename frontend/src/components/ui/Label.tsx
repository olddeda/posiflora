import type { ReactNode } from 'react'

interface Props {
  htmlFor?: string
  children: ReactNode
}

export const Label = ({ htmlFor, children }: Props) => (
  <label
    htmlFor={htmlFor}
    className="block text-sm font-semibold text-gray-800"
  >
    {children}
  </label>
)
