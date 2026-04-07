import { useTranslation } from 'react-i18next'

interface Props {
  enabled: boolean
}

export const StatusBadge = ({ enabled }: Props) => {
  const { t } = useTranslation()
  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold ${
        enabled ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'
      }`}
    >
      <span className={`w-1.5 h-1.5 rounded-full ${enabled ? 'bg-green-600' : 'bg-gray-400'}`} />
      {enabled ? t('status.active') : t('status.inactive')}
    </span>
  )
}
