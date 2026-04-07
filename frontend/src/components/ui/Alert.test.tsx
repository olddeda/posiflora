import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Alert } from './Alert'

describe('Alert', () => {
  it('renders children', () => {
    render(<Alert type="success">Saved!</Alert>)
    expect(screen.getByText('Saved!')).toBeInTheDocument()
  })

  it('applies success styles', () => {
    render(<Alert type="success">OK</Alert>)
    const el = screen.getByText('OK')
    expect(el.className).toContain('text-green-700')
    expect(el.className).toContain('bg-green-50')
  })

  it('applies error styles', () => {
    render(<Alert type="error">Fail</Alert>)
    const el = screen.getByText('Fail')
    expect(el.className).toContain('text-red-600')
    expect(el.className).toContain('bg-red-50')
  })
})
