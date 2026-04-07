import i18n from '../i18n'

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status?: number,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

const safeJson = async (res: Response): Promise<unknown> => {
  try {
    return await res.json()
  } catch {
    throw new ApiError(i18n.t('error.badResponse'), res.status)
  }
}

export const apiFetch = async <T>(
  url: string,
  options?: RequestInit,
  fallback?: string,
): Promise<T> => {
  let res: Response
  try {
    res = await fetch(url, options)
  } catch {
    throw new ApiError(i18n.t('error.offline'))
  }

  const data = await safeJson(res)

  if (!res.ok) {
    const msg =
      data && typeof data === 'object' && 'error' in data
        ? String((data as { error: string }).error)
        : (fallback ?? i18n.t('error.unknown'))
    throw new ApiError(msg, res.status)
  }

  return data as T
}
