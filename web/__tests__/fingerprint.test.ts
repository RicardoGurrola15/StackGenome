import { describe, it, expect } from 'vitest'
import { extractFingerprint, describeFingerprintPayload } from '../lib/fingerprint'
import type { ProjectGraphDTO } from '../types/stackgenome'

describe('fingerprint lib', () => {
  it('should extract fingerprint without evidences or sensitive properties', () => {
    const graph: ProjectGraphDTO = {
      schema_version: '1.0.0',
      nodes: [
        {
          id: 'n1',
          type: 'framework',
          name: 'Next.js',
          confidence: 1,
          properties: {
            workspace: 'app',
            path: '/users/secret/path', // sensitive
            token: '12345' // sensitive
          },
          evidences: [
            { kind: 'manifest', path: '/users/secret/package.json', value: '14.0.0' }
          ]
        }
      ],
      edges: []
    }

    const fp = extractFingerprint(graph)
    
    // Evidences should be stripped
    expect((fp.nodes[0] as any).evidences).toBeUndefined()
    
    // Sensitive properties should be stripped
    expect(fp.nodes[0].properties).toBeDefined()
    expect(fp.nodes[0].properties?.workspace).toBe('app')
    expect(fp.nodes[0].properties?.path).toBeUndefined()
    expect(fp.nodes[0].properties?.token).toBeUndefined()
  })

  it('describeFingerprintPayload should output a readable string', () => {
    const fp = {
      nodes: [
        { id: '1', type: 'primary', name: 'a', confidence: 1 },
        { id: '2', type: 'framework', name: 'b', confidence: 1 }
      ],
      edges: [{ source_id: '1', target_id: '2', relation: 'uses' }]
    }
    const desc = describeFingerprintPayload(fp)
    expect(desc).toContain('2 nodos (framework, primary)')
    expect(desc).toContain('1 relaciones')
    expect(desc).toContain('Sin evidencias')
  })
})
