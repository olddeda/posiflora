import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Toggle } from './Toggle'

describe('Toggle', () => {
  it('renders with correct aria-checked when checked', () => {
    render(
      <Toggle
        checked={true}
        onChange={vi.fn()}
      />,
    )
    expect(screen.getByRole('switch')).toHaveAttribute('aria-checked', 'true')
  })

  it('renders with correct aria-checked when unchecked', () => {
    render(
      <Toggle
        checked={false}
        onChange={vi.fn()}
      />,
    )
    expect(screen.getByRole('switch')).toHaveAttribute('aria-checked', 'false')
  })

  it('calls onChange with toggled value on click', async () => {
    const onChange = vi.fn()
    render(
      <Toggle
        checked={false}
        onChange={onChange}
      />,
    )
    await userEvent.click(screen.getByRole('switch'))
    expect(onChange).toHaveBeenCalledWith(true)
  })

  it('calls onChange with false when currently checked', async () => {
    const onChange = vi.fn()
    render(
      <Toggle
        checked={true}
        onChange={onChange}
      />,
    )
    await userEvent.click(screen.getByRole('switch'))
    expect(onChange).toHaveBeenCalledWith(false)
  })

  it('renders aria-label when provided', () => {
    render(
      <Toggle
        checked={false}
        onChange={vi.fn()}
        label="Enable notifications"
      />,
    )
    expect(screen.getByRole('switch')).toHaveAttribute('aria-label', 'Enable notifications')
  })
})
