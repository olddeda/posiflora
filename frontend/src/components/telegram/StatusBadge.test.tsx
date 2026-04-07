import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { StatusBadge } from './StatusBadge'

describe('StatusBadge', () => {
  it('shows active key when enabled', () => {
    render(<StatusBadge enabled={true} />)
    expect(screen.getByText('status.active')).toBeInTheDocument()
  })

  it('shows inactive key when disabled', () => {
    render(<StatusBadge enabled={false} />)
    expect(screen.getByText('status.inactive')).toBeInTheDocument()
  })

  it('applies green styles when enabled', () => {
    const { container } = render(<StatusBadge enabled={true} />)
    expect(container.firstChild).toHaveClass('bg-green-100', 'text-green-700')
  })

  it('applies gray styles when disabled', () => {
    const { container } = render(<StatusBadge enabled={false} />)
    expect(container.firstChild).toHaveClass('bg-gray-100', 'text-gray-500')
  })
})
