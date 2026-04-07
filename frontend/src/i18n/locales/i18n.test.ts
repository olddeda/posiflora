import { describe, it, expect } from 'vitest'
import ru from './ru.json'
import en from './en.json'

const flatKeys = (obj: object, prefix = ''): string[] =>
  Object.entries(obj).flatMap(([key, val]) => {
    const full = prefix ? `${prefix}.${key}` : key
    return typeof val === 'object' && val !== null ? flatKeys(val as object, full) : [full]
  })

describe('i18n locale parity', () => {
  const ruKeys = flatKeys(ru).sort()
  const enKeys = flatKeys(en).sort()

  it('en has all keys that ru has', () => {
    const missing = ruKeys.filter((k) => !enKeys.includes(k))
    expect(missing).toEqual([])
  })

  it('ru has all keys that en has', () => {
    const missing = enKeys.filter((k) => !ruKeys.includes(k))
    expect(missing).toEqual([])
  })

  it('neither locale is empty', () => {
    expect(ruKeys.length).toBeGreaterThan(0)
    expect(enKeys.length).toBeGreaterThan(0)
  })
})
