export interface TelegramIntegration {
  id: number
  shopId: number
  chatId: string
  enabled: boolean
  createdAt: string
  updatedAt: string
}

export interface IntegrationStatus {
  enabled: boolean
  chatId: string
  lastSentAt: string | null
  sentCount7d: number
  failedCount7d: number
}

export interface ConnectPayload {
  botToken: string
  chatId: string
  enabled: boolean
}

export interface CreateOrderPayload {
  number: string
  total: number
  customerName: string
}

export interface CreateOrderResponse {
  order: {
    id: number
    number: string
    total: number
    customerName: string
    createdAt: string
  }
  notifyStatus: 'sent' | 'failed' | 'skipped'
}

export interface ApiError {
  error: string
}
