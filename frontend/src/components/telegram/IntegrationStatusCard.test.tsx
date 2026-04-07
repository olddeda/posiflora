import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { IntegrationStatusCard } from './IntegrationStatusCard'
import type { IntegrationStatus } from '../../types/telegram'

const fullStatus: IntegrationStatus = {
  enabled: true,
  chatId: '****1234',
  lastSentAt: '2024-01-15T10:30:00Z',
  sentCount7d: 42,
  failedCount7d: 3,
}

describe('IntegrationStatusCard', () => {
  it('shows loading text when loading', () => {
    render(
      <IntegrationStatusCard
        status={null}
        loading={true}
        error={null}
      />,
    )
    expect(screen.getByText('status.loading')).toBeInTheDocument()
  })

  it('shows error text when error is set', () => {
    render(
      <IntegrationStatusCard
        status={null}
        loading={false}
        error="Сервер недоступен"
      />,
    )
    expect(screen.getByText('Сервер недоступен')).toBeInTheDocument()
  })

  it('shows not-configured text when status is null and no error', () => {
    render(
      <IntegrationStatusCard
        status={null}
        loading={false}
        error={null}
      />,
    )
    expect(screen.getByText('status.notConfigured')).toBeInTheDocument()
  })

  it('shows status data when status is provided', () => {
    render(
      <IntegrationStatusCard
        status={fullStatus}
        loading={false}
        error={null}
      />,
    )
    expect(screen.getByText('****1234')).toBeInTheDocument()
    expect(screen.getByText('42')).toBeInTheDocument()
    expect(screen.getByText('3')).toBeInTheDocument()
  })

  it('shows sent and failed counts', () => {
    render(
      <IntegrationStatusCard
        status={fullStatus}
        loading={false}
        error={null}
      />,
    )
    expect(screen.getByText('status.sent7d')).toBeInTheDocument()
    expect(screen.getByText('status.failed7d')).toBeInTheDocument()
  })

  it('does not show chatId section when chatId is empty', () => {
    render(
      <IntegrationStatusCard
        status={{ ...fullStatus, chatId: '' }}
        loading={false}
        error={null}
      />,
    )
    expect(screen.queryByText('status.chatId')).not.toBeInTheDocument()
  })

  it('error takes priority over status', () => {
    render(
      <IntegrationStatusCard
        status={fullStatus}
        loading={false}
        error="Error!"
      />,
    )
    expect(screen.getByText('Error!')).toBeInTheDocument()
    expect(screen.queryByText('42')).not.toBeInTheDocument()
  })

  it('loading takes priority over everything', () => {
    render(
      <IntegrationStatusCard
        status={fullStatus}
        loading={true}
        error="Error!"
      />,
    )
    expect(screen.getByText('status.loading')).toBeInTheDocument()
    expect(screen.queryByText('42')).not.toBeInTheDocument()
  })
})
