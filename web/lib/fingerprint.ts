import type { NodeDTO, FingerprintV1, ProjectGraphDTO } from '@/types/stackgenome'

// Fields allowed in properties when sending remote fingerprint
const ALLOWED_PROPERTY_KEYS = new Set(['primary', 'role', 'workspace'])

function sanitizeNode(node: NodeDTO): NodeDTO {
  const sanitized: NodeDTO = {
    id: node.id,
    type: node.type,
    name: node.name,
    confidence: node.confidence,
  }
  if (node.role) sanitized.role = node.role
  if (node.version) sanitized.version = node.version

  // Only include allowlisted property keys — strip paths, usernames, etc.
  if (node.properties) {
    const cleanProps: Record<string, string> = {}
    for (const [k, v] of Object.entries(node.properties)) {
      if (ALLOWED_PROPERTY_KEYS.has(k)) {
        cleanProps[k] = v
      }
    }
    if (Object.keys(cleanProps).length > 0) {
      sanitized.properties = cleanProps
    }
  }

  // Evidences are NEVER sent remotely
  return sanitized
}

export function extractFingerprint(graph: ProjectGraphDTO): FingerprintV1 {
  return {
    nodes: graph.nodes.map(sanitizeNode),
    edges: graph.edges.map(e => ({
      source_id: e.source_id,
      target_id: e.target_id,
      relation: e.relation,
    })),
  }
}

export function describeFingerprintPayload(fp: FingerprintV1): string {
  const nodeTypes = [...new Set(fp.nodes.map(n => n.type))].sort()
  return `${fp.nodes.length} nodos (${nodeTypes.join(', ')}) · ${fp.edges.length} relaciones · Sin evidencias · Sin rutas de archivos`
}
