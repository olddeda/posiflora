import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { OfflineBanner } from './OfflineBanner'

describe('OfflineBanner', () => {
  it('shows custom message when provided', () => {
    render(<OfflineBanner message="Custom error" />)
    expect(screen.getByText('Custom error')).toBeInTheDocument()
  })

  it('shows i18n fallback when no message', () => {
    render(<OfflineBanner />)
    expect(screen.getByText('error.offline')).toBeInTheDocument()
  })

  it('does not render retry button when onRetry is not provided', () => {
    render(<OfflineBanner message="Error" />)
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })

  it('renders retry button when onRetry is provided', () => {
    render(
      <OfflineBanner
        message="Error"
        onRetry={vi.fn()}
      />,
    )
    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  it('calls onRetry when retry button clicked', async () => {
    const onRetry = vi.fn()
    render(
      <OfflineBanner
        message="Error"
        onRetry={onRetry}
      />,
    )
    await userEvent.click(screen.getByRole('button'))
    expect(onRetry).toHaveBeenCalledOnce()
  })

  it('shows plug emoji', () => {
    render(<OfflineBanner />)
    expect(screen.getByText('🔌')).toBeInTheDocument()
  })
})
