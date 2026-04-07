import { useEffect, useRef, useMemo } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation, Trans } from 'react-i18next'
import { connectTelegram } from '../../api/telegram'
import type { IntegrationStatus } from '../../types/telegram'
import { Alert } from '../ui/Alert'
import { Card } from '../ui/Card'
import { Input } from '../ui/Input'
import { Label } from '../ui/Label'
import { Toggle } from '../ui/Toggle'

interface Props {
  shopId: number
  initialStatus: IntegrationStatus | null
  onSaved: () => void
}

const BOT_TOKEN_RE = /^\d+:[A-Za-z0-9_-]{35,}$/

export const ConnectBotCard = ({ shopId, initialStatus, onSaved }: Props) => {
  const { t } = useTranslation()

  const isConfigured = Boolean(initialStatus?.chatId)

  const schema = useMemo(
    () =>
      z.object({
        botToken: isConfigured
          ? z.union([
              z.literal(''),
              z.string().min(1).regex(BOT_TOKEN_RE, t('validation.botTokenFormat')),
            ])
          : z
              .string()
              .min(1, t('validation.required'))
              .regex(BOT_TOKEN_RE, t('validation.botTokenFormat')),
        chatId: z
          .string()
          .min(1, t('validation.required'))
          .regex(/^-?\d+$/, t('validation.chatIdFormat')),
        enabled: z.boolean(),
      }),
    [isConfigured, t],
  )

  type FormValues = z.infer<typeof schema>

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors, isSubmitting, isSubmitSuccessful },
    setError,
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { botToken: '', chatId: '', enabled: true },
  })

  const initialized = useRef(false)

  useEffect(() => {
    if (initialStatus && !initialized.current) {
      initialized.current = true
      reset({ botToken: '', chatId: initialStatus.chatId, enabled: initialStatus.enabled })
    }
  }, [initialStatus, reset])

  const onSubmit = async (data: FormValues) => {
    try {
      await connectTelegram(shopId, data)
      reset({ botToken: '', chatId: data.chatId, enabled: data.enabled })
      onSaved()
    } catch (err: unknown) {
      setError('root', {
        message: err instanceof Error ? err.message : t('connect.saveError'),
      })
    }
  }

  return (
    <Card>
      <div className="text-base font-bold text-gray-900 mb-4">
        {isConfigured ? t('connect.titleUpdate') : t('connect.title')}
      </div>

      {isSubmitSuccessful && !errors.root && (
        <div className="mb-4">
          <Alert type="success">{t('connect.saveSuccess')}</Alert>
        </div>
      )}
      {errors.root && (
        <div className="mb-4">
          <Alert type="error">{errors.root.message}</Alert>
        </div>
      )}

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="space-y-4"
      >
        <div>
          <div className="flex items-baseline justify-between">
            <Label htmlFor="botToken">
              {t('connect.botToken')}{' '}
              <a
                href="https://t.me/BotFather"
                target="_blank"
                rel="noopener noreferrer"
                className="text-xs font-normal text-blue-500 hover:underline"
              >
                {t('connect.botTokenHint')}
              </a>
            </Label>
            <span className="text-xs font-medium text-gray-400">
              {isConfigured ? t('connect.optional') : t('connect.required')}
            </span>
          </div>
          <Input
            id="botToken"
            className="mt-1.5"
            placeholder={
              isConfigured
                ? t('connect.botTokenUpdatePlaceholder')
                : t('connect.botTokenPlaceholder')
            }
            {...register('botToken')}
          />
          {errors.botToken && (
            <p className="mt-1 text-xs text-red-500">{errors.botToken.message}</p>
          )}
        </div>

        <div>
          <Label htmlFor="chatId">{t('connect.chatId')}</Label>
          <Input
            id="chatId"
            className="mt-1.5"
            placeholder={t('connect.chatIdPlaceholder')}
            {...register('chatId')}
          />
          {errors.chatId && <p className="mt-1 text-xs text-red-500">{errors.chatId.message}</p>}
        </div>

        <div className="flex items-center gap-3">
          <Controller
            name="enabled"
            control={control}
            render={({ field }) => (
              <Toggle
                checked={field.value}
                onChange={field.onChange}
                label={t('connect.enabledLabel')}
              />
            )}
          />
          <span className="text-sm text-gray-700">{t('connect.enabledLabel')}</span>
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="inline-flex items-center px-5 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white text-sm font-semibold rounded-lg transition-colors cursor-pointer"
        >
          {isSubmitting
            ? t('connect.saving')
            : isConfigured
              ? t('connect.update')
              : t('connect.save')}
        </button>
      </form>

      <div className="border-t border-gray-200 my-5" />

      <div className="bg-amber-50 border border-amber-200 rounded-lg px-4 py-3 text-xs text-amber-900 leading-relaxed">
        <strong>{t('hint.title')}</strong>
        <ol className="list-decimal pl-4 mt-1.5 space-y-1">
          <li>{t('hint.step1')}</li>
          <li>
            {t('hint.step2prefix')}{' '}
            <a
              href={`https://api.telegram.org/bot${'<TOKEN>'}/getUpdates`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:underline break-all"
            >
              <code>{`https://api.telegram.org/bot${'<TOKEN>'}/getUpdates`}</code>
            </a>
          </li>
          <li>
            <Trans
              i18nKey="hint.step3"
              components={{ code: <code className="bg-amber-100 px-1 rounded" /> }}
            />
          </li>
          <li>{t('hint.step4')}</li>
        </ol>
      </div>
    </Card>
  )
}
