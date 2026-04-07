import { useTranslation } from 'react-i18next'
import type { IntegrationStatus } from '../../types/telegram'
import { Card } from '../ui/Card'
import { StatusBadge } from './StatusBadge'

interface Props {
  status: IntegrationStatus | null
  loading: boolean
  error: string | null
}

export const IntegrationStatusCard = ({ status, loading, error }: Props) => {
  const { t, i18n } = useTranslation()

  return (
    <Card>
      <div className="text-base font-bold text-gray-900 mb-4">{t('status.title')}</div>
      {loading ? (
        <p className="text-sm text-gray-400">{t('status.loading')}</p>
      ) : error ? (
        <p className="text-sm text-red-500">{error}</p>
      ) : status ? (
        <>
          <div className="flex items-center gap-3 mb-2">
            <StatusBadge enabled={status.enabled} />
            {status.chatId && (
              <span className="text-xs text-gray-500">
                {t('status.chatId')}:{' '}
                <code className="bg-gray-100 px-1.5 py-0.5 rounded">{status.chatId}</code>
              </span>
            )}
          </div>
          {status.lastSentAt && (
            <p className="text-xs text-gray-400 mb-4">
              {t('status.lastSent', {
                date: new Date(status.lastSentAt).toLocaleString(i18n.language),
              })}
            </p>
          )}
          <div className="grid grid-cols-2 gap-3 mt-4">
            <div className="bg-gray-50 border border-gray-200 rounded-lg px-4 py-3">
              <div className="text-xs text-gray-400 mb-1">{t('status.sent7d')}</div>
              <div className="text-2xl font-bold text-green-600">{status.sentCount7d}</div>
            </div>
            <div className="bg-gray-50 border border-gray-200 rounded-lg px-4 py-3">
              <div className="text-xs text-gray-400 mb-1">{t('status.failed7d')}</div>
              <div
                className={`text-2xl font-bold ${status.failedCount7d > 0 ? 'text-red-500' : 'text-gray-800'}`}
              >
                {status.failedCount7d}
              </div>
            </div>
          </div>
        </>
      ) : (
        <p className="text-sm text-gray-400">{t('status.notConfigured')}</p>
      )}
    </Card>
  )
}
