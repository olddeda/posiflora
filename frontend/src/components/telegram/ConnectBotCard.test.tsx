import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ConnectBotCard } from './ConnectBotCard'

vi.mock('../../api/telegram', () => ({
  connectTelegram: vi.fn(),
}))

import { connectTelegram } from '../../api/telegram'
const mockConnect = vi.mocked(connectTelegram)

const VALID_TOKEN = '123456789:AABBCCDDEEFFaabbccddeeff-1234567890a'
const VALID_CHAT_ID = '987654321'

const fillAndSubmit = async (
  user: ReturnType<typeof userEvent.setup>,
  token: string,
  chatId: string,
) => {
  const tokenInput = screen.getByLabelText(/connect\.botToken/i)
  const chatInput = screen.getByLabelText(/connect\.chatId/i)
  const btn = screen.getByRole('button', { name: /connect\.save/i })

  if (token) {
    await user.clear(tokenInput)
    await user.type(tokenInput, token)
  }
  if (chatId) {
    await user.clear(chatInput)
    await user.type(chatInput, chatId)
  }
  await user.click(btn)
}

beforeEach(() => {
  mockConnect.mockReset()
})

describe('ConnectBotCard — not configured', () => {
  it('renders connect title when no initial status', () => {
    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={null}
        onSaved={vi.fn()}
      />,
    )
    expect(screen.getByText('connect.title')).toBeInTheDocument()
  })

  it('shows validation error when botToken is empty on submit', async () => {
    const user = userEvent.setup()
    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={null}
        onSaved={vi.fn()}
      />,
    )

    await user.click(screen.getByRole('button', { name: /connect\.save/i }))

    await waitFor(() => {
      expect(screen.getAllByText('validation.required').length).toBeGreaterThanOrEqual(1)
    })
    expect(mockConnect).not.toHaveBeenCalled()
  })

  it('shows validation error for invalid botToken format', async () => {
    const user = userEvent.setup()
    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={null}
        onSaved={vi.fn()}
      />,
    )

    await user.type(screen.getByLabelText(/connect\.botToken/i), 'bad-token')
    await user.type(screen.getByLabelText(/connect\.chatId/i), VALID_CHAT_ID)
    await user.click(screen.getByRole('button', { name: /connect\.save/i }))

    await waitFor(() => {
      expect(screen.getByText('validation.botTokenFormat')).toBeInTheDocument()
    })
    expect(mockConnect).not.toHaveBeenCalled()
  })

  it('shows validation error for invalid chatId (non-numeric)', async () => {
    const user = userEvent.setup()
    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={null}
        onSaved={vi.fn()}
      />,
    )

    await user.type(screen.getByLabelText(/connect\.botToken/i), VALID_TOKEN)
    await user.type(screen.getByLabelText(/connect\.chatId/i), 'not-a-number')
    await user.click(screen.getByRole('button', { name: /connect\.save/i }))

    await waitFor(() => {
      expect(screen.getByText('validation.chatIdFormat')).toBeInTheDocument()
    })
  })

  it('calls connectTelegram with correct args on valid submit', async () => {
    const user = userEvent.setup()
    const onSaved = vi.fn()
    mockConnect.mockResolvedValue({} as never)

    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={null}
        onSaved={onSaved}
      />,
    )
    await fillAndSubmit(user, VALID_TOKEN, VALID_CHAT_ID)

    await waitFor(() => {
      expect(mockConnect).toHaveBeenCalledWith(1, {
        botToken: VALID_TOKEN,
        chatId: VALID_CHAT_ID,
        enabled: true,
      })
    })
  })

  it('shows success alert and calls onSaved after successful submit', async () => {
    const user = userEvent.setup()
    const onSaved = vi.fn()
    mockConnect.mockResolvedValue({} as never)

    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={null}
        onSaved={onSaved}
      />,
    )
    await fillAndSubmit(user, VALID_TOKEN, VALID_CHAT_ID)

    await waitFor(() => {
      expect(screen.getByText('connect.saveSuccess')).toBeInTheDocument()
    })
    expect(onSaved).toHaveBeenCalledOnce()
  })

  it('shows error alert when API call fails', async () => {
    const user = userEvent.setup()
    mockConnect.mockRejectedValue(new Error('Connection failed'))

    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={null}
        onSaved={vi.fn()}
      />,
    )
    await fillAndSubmit(user, VALID_TOKEN, VALID_CHAT_ID)

    await waitFor(() => {
      expect(screen.getByText('Connection failed')).toBeInTheDocument()
    })
  })
})

describe('ConnectBotCard — already configured', () => {
  const existingStatus = {
    enabled: true,
    chatId: '987654321',
    lastSentAt: null,
    sentCount7d: 0,
    failedCount7d: 0,
  }

  it('renders update title when status exists', () => {
    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={existingStatus}
        onSaved={vi.fn()}
      />,
    )
    expect(screen.getByText('connect.titleUpdate')).toBeInTheDocument()
  })

  it('pre-fills chatId from existing status', () => {
    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={existingStatus}
        onSaved={vi.fn()}
      />,
    )
    expect(screen.getByLabelText(/connect\.chatId/i)).toHaveValue('987654321')
  })

  it('allows submit with empty botToken when already configured', async () => {
    const user = userEvent.setup()
    mockConnect.mockResolvedValue({} as never)

    render(
      <ConnectBotCard
        shopId={1}
        initialStatus={existingStatus}
        onSaved={vi.fn()}
      />,
    )
    await user.click(screen.getByRole('button', { name: /connect\.update/i }))

    await waitFor(() => {
      expect(mockConnect).toHaveBeenCalledWith(
        1,
        expect.objectContaining({ botToken: '', chatId: '987654321' }),
      )
    })
  })
})
