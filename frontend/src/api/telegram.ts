import { apiFetch } from '../lib/apiFetch'
import type {
  TelegramIntegration,
  IntegrationStatus,
  ConnectPayload,
  CreateOrderPayload,
  CreateOrderResponse,
} from '../types/telegram'

export type {
  TelegramIntegration,
  IntegrationStatus,
  ConnectPayload,
  CreateOrderPayload,
  CreateOrderResponse,
}

const BASE = import.meta.env.VITE_API_URL ?? ''

export const connectTelegram = (shopId: number, payload: ConnectPayload) =>
  apiFetch<TelegramIntegration>(`${BASE}/shops/${shopId}/telegram/connect`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

export const getTelegramStatus = (shopId: number) =>
  apiFetch<IntegrationStatus>(`${BASE}/shops/${shopId}/telegram/status`)

export const createOrder = (shopId: number, payload: CreateOrderPayload) =>
  apiFetch<CreateOrderResponse>(`${BASE}/shops/${shopId}/orders`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
