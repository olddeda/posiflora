import { useEffect, useState, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { getTelegramStatus } from '../api/telegram'
import type { IntegrationStatus } from '../types/telegram'
import { IntegrationStatusCard } from '../components/telegram/IntegrationStatusCard'
import { ConnectBotCard } from '../components/telegram/ConnectBotCard'
import { TestOrderCard } from '../components/telegram/TestOrderCard'
import { OfflineBanner } from '../components/OfflineBanner'

export const TelegramIntegrationPage = () => {
  const { shopId } = useParams<{ shopId: string }>()
  const { t, i18n } = useTranslation()

  const id = Number(shopId)

  const [status, setStatus] = useState<IntegrationStatus | null>(null)
  const [statusLoading, setStatusLoading] = useState(true)
  const [statusError, setStatusError] = useState<string | null>(null)

  const fetchStatus = useCallback(async () => {
    setStatusError(null)
    try {
      setStatus(await getTelegramStatus(id))
    } catch (err: unknown) {
      setStatusError(err instanceof Error ? err.message : t('status.loading'))
    } finally {
      setStatusLoading(false)
    }
  }, [id, t])

  useEffect(() => {
    fetchStatus()
  }, [fetchStatus])

  if (!shopId || isNaN(id)) {
    return (
      <div className="max-w-xl mx-auto px-4 py-10">
        <p className="text-sm text-red-500">Invalid shop ID</p>
      </div>
    )
  }

  const otherLang = i18n.language.startsWith('ru') ? 'en' : 'ru'

  return (
    <div className="max-w-xl mx-auto px-4 py-10">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <span
            className="text-3xl"
            aria-hidden="true"
          >
            🤖
          </span>
          <div>
            <h1 className="text-xl font-bold text-gray-900">{t('header.title')}</h1>
            <p className="text-sm text-gray-500">{t('header.subtitle', { id })}</p>
          </div>
        </div>
        <button
          onClick={() => i18n.changeLanguage(otherLang)}
          className="text-xs font-semibold text-gray-400 hover:text-gray-700 border border-gray-200 rounded-md px-2.5 py-1 transition-colors cursor-pointer"
        >
          {otherLang.toUpperCase()}
        </button>
      </div>

      <div className="space-y-5">
        {statusError && (
          <OfflineBanner
            message={statusError}
            onRetry={fetchStatus}
          />
        )}
        <IntegrationStatusCard
          status={status}
          loading={statusLoading}
          error={statusError}
        />
        <ConnectBotCard
          shopId={id}
          initialStatus={status}
          onSaved={fetchStatus}
        />
        {status?.enabled && status?.chatId && (
          <TestOrderCard
            shopId={id}
            onCreated={fetchStatus}
          />
        )}
      </div>
    </div>
  )
}
