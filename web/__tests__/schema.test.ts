import { describe, it, expect } from 'vitest'
import { parseGraphFile, checkFileSize } from '../lib/schema'
import type { ProjectGraphDTO } from '../types/stackgenome'

describe('schema lib', () => {
  it('checkFileSize should return true for files under 5MB', () => {
    const file = new File(['a'], 'test.json', { type: 'application/json' })
    expect(checkFileSize(file)).toBe(true)
  })

  it('parseGraphFile should fail on invalid JSON', () => {
    const res = parseGraphFile('not json')
    expect(res.ok).toBe(false)
    if (!res.ok) {
      expect(res.error).toContain('JSON válido')
    }
  })

  it('parseGraphFile should fail on missing schema_version', () => {
    const res = parseGraphFile(JSON.stringify({ nodes: [], edges: [] }))
    expect(res.ok).toBe(false)
  })

  it('parseGraphFile should succeed on valid DTO', () => {
    const valid: ProjectGraphDTO = {
      schema_version: '1.0.0',
      nodes: [
        { id: 'n1', type: 'primary', name: 'Go', confidence: 1 },
      ],
      edges: [],
    }
    const res = parseGraphFile(JSON.stringify(valid))
    expect(res.ok).toBe(true)
    if (res.ok) {
      expect(res.data.nodes[0].name).toBe('Go')
    }
  })
})
