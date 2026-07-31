import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getCatalog, getRecommendations } from '../lib/api'

// Mock global fetch
const fetchMock = vi.fn()
global.fetch = fetchMock

describe('api lib', () => {
  beforeEach(() => {
    fetchMock.mockReset()
  })

  it('getCatalog should return resources on success', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ resources: [{ id: 'test', name: 'Test' }] })
    })

    const res = await getCatalog()
    expect(res.ok).toBe(true)
    if (res.ok) {
      expect(res.data[0].id).toBe('test')
    }
  })

  it('getCatalog should handle HTTP errors', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: false,
      status: 500
    })

    const res = await getCatalog()
    expect(res.ok).toBe(false)
    if (!res.ok) {
      expect(res.status).toBe(500)
    }
  })

  it('getRecommendations should handle network errors', async () => {
    fetchMock.mockRejectedValueOnce(new Error('Network error'))

    const req = { schema_version: '1.0.0', fingerprint: { nodes: [], edges: [] } }
    const res = await getRecommendations(req)
    
    expect(res.ok).toBe(false)
    if (!res.ok) {
      expect(res.error).toContain('No se pudo conectar')
    }
  })
})
