import { useMemo, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation, Trans } from 'react-i18next'
import { createOrder } from '../../api/telegram'
import type { CreateOrderResponse } from '../../types/telegram'
import { Alert } from '../ui/Alert'
import { Card } from '../ui/Card'
import { Input } from '../ui/Input'
import { Label } from '../ui/Label'

interface Props {
  shopId: number
  onCreated: () => void
}

const notifyClass = (s: CreateOrderResponse['notifyStatus']) =>
  s === 'sent' ? 'text-green-600' : s === 'failed' ? 'text-red-500' : 'text-gray-400'

export const TestOrderCard = ({ shopId, onCreated }: Props) => {
  const { t } = useTranslation()

  const schema = useMemo(
    () =>
      z.object({
        number: z.string().min(1, t('validation.required')),
        total: z.number({ error: t('validation.number') }).positive(t('validation.positive')),
        customerName: z.string().min(1, t('validation.required')),
      }),
    [t],
  )

  type FormValues = z.infer<typeof schema>

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { number: '', total: NaN, customerName: '' },
  })

  const [result, setResult] = useState<CreateOrderResponse | null>(null)

  const onSubmit = async (data: FormValues) => {
    setResult(null)
    try {
      const res = await createOrder(shopId, {
        number: data.number,
        total: data.total,
        customerName: data.customerName,
      })
      setResult(res)
      onCreated()
    } catch (err: unknown) {
      setError('root', {
        message: err instanceof Error ? err.message : t('connect.saveError'),
      })
    }
  }

  return (
    <Card>
      <div className="text-base font-bold text-gray-900 mb-1">{t('order.title')}</div>
      <p className="text-xs text-gray-400 mb-4">{t('order.description')}</p>

      {result && (
        <div className="mb-4">
          <Alert type="success">
            <Trans
              i18nKey="order.created"
              values={{ number: result.order.number }}
              components={{ strong: <strong /> }}
            />{' '}
            <span className={notifyClass(result.notifyStatus)}>
              {t(`order.notify.${result.notifyStatus}`)}
            </span>
          </Alert>
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
          <Label htmlFor="orderNumber">{t('order.number')}</Label>
          <Input
            id="orderNumber"
            className="mt-1.5"
            placeholder={t('order.numberPlaceholder')}
            {...register('number')}
          />
          {errors.number && <p className="mt-1 text-xs text-red-500">{errors.number.message}</p>}
        </div>

        <div>
          <Label htmlFor="orderTotal">{t('order.total')}</Label>
          <Input
            id="orderTotal"
            type="number"
            className="mt-1.5"
            placeholder={t('order.totalPlaceholder')}
            {...register('total', { valueAsNumber: true })}
          />
          {errors.total && <p className="mt-1 text-xs text-red-500">{errors.total.message}</p>}
        </div>

        <div>
          <Label htmlFor="orderCustomer">{t('order.customer')}</Label>
          <Input
            id="orderCustomer"
            className="mt-1.5"
            placeholder={t('order.customerPlaceholder')}
            {...register('customerName')}
          />
          {errors.customerName && (
            <p className="mt-1 text-xs text-red-500">{errors.customerName.message}</p>
          )}
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="inline-flex items-center px-5 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white text-sm font-semibold rounded-lg transition-colors cursor-pointer"
        >
          {isSubmitting ? t('order.submitting') : t('order.submit')}
        </button>
      </form>
    </Card>
  )
}
