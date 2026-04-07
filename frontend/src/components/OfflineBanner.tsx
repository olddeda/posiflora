import { useTranslation } from 'react-i18next'

interface Props {
  message?: string
  onRetry?: () => void
}

export const OfflineBanner = ({ message, onRetry }: Props) => {
  const { t } = useTranslation()
  return (
    <div className="rounded-xl border border-red-200 bg-red-50 px-5 py-4 flex items-start gap-3">
      <span className="text-xl mt-0.5">🔌</span>
      <div className="flex-1">
        <p className="text-sm font-semibold text-red-700">{message ?? t('error.offline')}</p>
        {onRetry && (
          <button
            onClick={onRetry}
            className="mt-2 text-xs font-semibold text-red-600 hover:text-red-800 underline cursor-pointer"
          >
            {t('error.retry', 'Повторить')}
          </button>
        )}
      </div>
    </div>
  )
}
