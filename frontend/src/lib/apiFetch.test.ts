import { describe, it, expect, vi, beforeEach } from 'vitest'
import { apiFetch, ApiError } from './apiFetch'

const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

const makeResponse = (body: unknown, status = 200): Response =>
  ({
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  }) as unknown as Response

const makeBrokenResponse = (status = 502): Response =>
  ({
    ok: false,
    status,
    json: async () => {
      throw new SyntaxError('Unexpected end of JSON input')
    },
  }) as unknown as Response

beforeEach(() => {
  mockFetch.mockReset()
})

describe('apiFetch', () => {
  it('returns parsed data on 200', async () => {
    mockFetch.mockResolvedValue(makeResponse({ id: 1, name: 'Shop' }))
    const result = await apiFetch<{ id: number; name: string }>('/shops/1')
    expect(result).toEqual({ id: 1, name: 'Shop' })
  })

  it('throws ApiError with offline message when fetch fails', async () => {
    mockFetch.mockRejectedValue(new TypeError('Failed to fetch'))
    await expect(apiFetch('/shops/1')).rejects.toMatchObject({
      name: 'ApiError',
      message: 'error.offline',
    })
  })

  it('throws ApiError with badResponse message when body is not JSON', async () => {
    mockFetch.mockResolvedValue(makeBrokenResponse(502))
    await expect(apiFetch('/shops/1')).rejects.toMatchObject({
      name: 'ApiError',
      message: 'error.badResponse',
      status: 502,
    })
  })

  it('throws ApiError with server error field on 4xx', async () => {
    mockFetch.mockResolvedValue(makeResponse({ error: 'shop not found' }, 404))
    await expect(apiFetch('/shops/99')).rejects.toMatchObject({
      name: 'ApiError',
      message: 'shop not found',
      status: 404,
    })
  })

  it('throws ApiError with fallback message when 4xx has no error field', async () => {
    mockFetch.mockResolvedValue(makeResponse({}, 400))
    await expect(apiFetch('/shops/1', undefined, 'custom fallback')).rejects.toMatchObject({
      name: 'ApiError',
      message: 'custom fallback',
    })
  })

  it('throws ApiError with unknown message when 4xx has no error field and no fallback', async () => {
    mockFetch.mockResolvedValue(makeResponse({}, 500))
    await expect(apiFetch('/shops/1')).rejects.toMatchObject({
      name: 'ApiError',
      message: 'error.unknown',
    })
  })

  it('passes options to fetch', async () => {
    mockFetch.mockResolvedValue(makeResponse({ ok: true }))
    await apiFetch('/endpoint', { method: 'POST', body: '{}' })
    expect(mockFetch).toHaveBeenCalledWith('/endpoint', { method: 'POST', body: '{}' })
  })
})

describe('ApiError', () => {
  it('has correct name and message', () => {
    const err = new ApiError('oops', 404)
    expect(err.name).toBe('ApiError')
    expect(err.message).toBe('oops')
    expect(err.status).toBe(404)
    expect(err instanceof Error).toBe(true)
  })

  it('status is optional', () => {
    const err = new ApiError('no status')
    expect(err.status).toBeUndefined()
  })
})
