// Types mirroring pkg/schema/v1/schema.go and backend/src/types/fingerprint.ts
// Single source of truth: docs/api/openapi.yaml

export interface EvidenceDTO {
  kind: string
  path?: string
  selector?: string
  value?: string
  sensitivity?: string
}

export interface NodeDTO {
  id: string
  type: string
  role?: string
  name: string
  version?: string
  confidence: number
  properties?: Record<string, string>
  evidences?: EvidenceDTO[]
}

export interface EdgeDTO {
  source_id: string
  target_id: string
  relation: string
}

export interface EnvironmentDTO {
  os?: string
  arch?: string
  tools?: Array<{ name: string; version?: string }>
}

export interface RecommendationResult {
  id: string
  name: string
  description: string
  score: number
  reasons: string[]
  url?: string
}

export interface ProjectGraphDTO {
  schema_version: string
  nodes: NodeDTO[]
  edges: EdgeDTO[]
  environment?: EnvironmentDTO
  recommendations?: RecommendationResult[]
}

export interface FingerprintV1 {
  nodes: NodeDTO[]
  edges: EdgeDTO[]
}

export interface RecommendationRequest {
  schema_version: string
  fingerprint: FingerprintV1
  limit?: number
}

export interface RecommendationResponse {
  catalog_snapshot_id: string
  ranking_version: string
  recommendations: RecommendationResult[]
  warnings: string[]
  generated_at: string
}

export interface CatalogResource {
  id: string
  type: string
  canonical_name: string
  summary: string
  canonical_url?: string
  ecosystem: string[]
  infra_targets: string[]
  license?: string
  status: 'active' | 'deprecated' | 'archived'
  updated_at: string
}

export interface ErrorResponse {
  error: string
  message: string
}
