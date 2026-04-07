import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TestOrderCard } from './TestOrderCard'

vi.mock('../../api/telegram', () => ({
  createOrder: vi.fn(),
}))

import { createOrder } from '../../api/telegram'
const mockCreateOrder = vi.mocked(createOrder)

const fillForm = async (
  user: ReturnType<typeof userEvent.setup>,
  number: string,
  total: string,
  customer: string,
) => {
  if (number) {
    await user.type(screen.getByLabelText(/order\.number/i), number)
  }
  if (total) {
    await user.type(screen.getByLabelText(/order\.total/i), total)
  }
  if (customer) {
    await user.type(screen.getByLabelText(/order\.customer/i), customer)
  }
}

beforeEach(() => {
  mockCreateOrder.mockReset()
})

describe('TestOrderCard', () => {
  it('renders form fields', () => {
    render(
      <TestOrderCard
        shopId={1}
        onCreated={vi.fn()}
      />,
    )
    expect(screen.getByLabelText(/order\.number/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/order\.total/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/order\.customer/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /order\.submit/i })).toBeInTheDocument()
  })

  it('shows required errors when submitting empty form', async () => {
    const user = userEvent.setup()
    render(
      <TestOrderCard
        shopId={1}
        onCreated={vi.fn()}
      />,
    )

    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      expect(screen.getAllByText('validation.required').length).toBeGreaterThanOrEqual(2)
    })
    expect(mockCreateOrder).not.toHaveBeenCalled()
  })

  it('shows validation error when total is zero or negative', async () => {
    const user = userEvent.setup()
    render(
      <TestOrderCard
        shopId={1}
        onCreated={vi.fn()}
      />,
    )

    await fillForm(user, 'A-001', '-5', 'Иван')
    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      expect(screen.getByText('validation.positive')).toBeInTheDocument()
    })
    expect(mockCreateOrder).not.toHaveBeenCalled()
  })

  it('calls createOrder with correct args on valid submit', async () => {
    const user = userEvent.setup()
    const onCreated = vi.fn()
    mockCreateOrder.mockResolvedValue({
      order: { id: 1, number: 'A-001', total: 2500, customerName: 'Иван', createdAt: '' },
      notifyStatus: 'sent',
    })

    render(
      <TestOrderCard
        shopId={1}
        onCreated={onCreated}
      />,
    )
    await fillForm(user, 'A-001', '2500', 'Иван')
    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      expect(mockCreateOrder).toHaveBeenCalledWith(1, {
        number: 'A-001',
        total: 2500,
        customerName: 'Иван',
      })
    })
    expect(onCreated).toHaveBeenCalledOnce()
  })

  it('shows success alert with "sent" status', async () => {
    const user = userEvent.setup()
    mockCreateOrder.mockResolvedValue({
      order: { id: 1, number: 'A-001', total: 1000, customerName: 'Мария', createdAt: '' },
      notifyStatus: 'sent',
    })

    render(
      <TestOrderCard
        shopId={1}
        onCreated={vi.fn()}
      />,
    )
    await fillForm(user, 'A-001', '1000', 'Мария')
    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      expect(screen.getByText('order.notify.sent')).toBeInTheDocument()
    })
  })

  it('shows "skipped" notify status when Telegram is disabled', async () => {
    const user = userEvent.setup()
    mockCreateOrder.mockResolvedValue({
      order: { id: 2, number: 'A-002', total: 500, customerName: 'Олег', createdAt: '' },
      notifyStatus: 'skipped',
    })

    render(
      <TestOrderCard
        shopId={1}
        onCreated={vi.fn()}
      />,
    )
    await fillForm(user, 'A-002', '500', 'Олег')
    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      expect(screen.getByText('order.notify.skipped')).toBeInTheDocument()
    })
  })

  it('shows "failed" notify status and applies red class', async () => {
    const user = userEvent.setup()
    mockCreateOrder.mockResolvedValue({
      order: { id: 3, number: 'A-003', total: 800, customerName: 'Анна', createdAt: '' },
      notifyStatus: 'failed',
    })

    render(
      <TestOrderCard
        shopId={1}
        onCreated={vi.fn()}
      />,
    )
    await fillForm(user, 'A-003', '800', 'Анна')
    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      const status = screen.getByText('order.notify.failed')
      expect(status).toHaveClass('text-red-500')
    })
  })

  it('shows error alert when API fails', async () => {
    const user = userEvent.setup()
    mockCreateOrder.mockRejectedValue(new Error('Server error'))

    render(
      <TestOrderCard
        shopId={1}
        onCreated={vi.fn()}
      />,
    )
    await fillForm(user, 'A-001', '1000', 'Иван')
    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      expect(screen.getByText('Server error')).toBeInTheDocument()
    })
  })

  it('does not call onCreated when API fails', async () => {
    const user = userEvent.setup()
    const onCreated = vi.fn()
    mockCreateOrder.mockRejectedValue(new Error('err'))

    render(
      <TestOrderCard
        shopId={1}
        onCreated={onCreated}
      />,
    )
    await fillForm(user, 'A-001', '100', 'X')
    await user.click(screen.getByRole('button', { name: /order\.submit/i }))

    await waitFor(() => {
      expect(screen.getByText('err')).toBeInTheDocument()
    })
    expect(onCreated).not.toHaveBeenCalled()
  })
})
