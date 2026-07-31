// Fingerprint v1 — matches ProjectGraphDTO from the Go CLI (pkg/schema/v1/schema.go)
export interface NodeDTO {
  id: string;
  type: string;
  role?: string;
  name: string;
  version?: string;
  confidence: number;
  properties?: Record<string, string>;
  // evidences intentionally omitted — not accepted from remote clients
}

export interface EdgeDTO {
  source_id: string;
  target_id: string;
  relation: string;
}

export interface FingerprintV1 {
  nodes: NodeDTO[];
  edges: EdgeDTO[];
}

export interface RecommendationsRequest {
  schema_version: string;
  fingerprint: FingerprintV1;
  limit?: number;
}

export interface RecommendationResult {
  id: string;
  name: string;
  description: string;
  score: number;
  reasons: string[];
  url?: string;
}

export interface RecommendationsResponse {
  catalog_snapshot_id: string;
  ranking_version: string;
  recommendations: RecommendationResult[];
  warnings: string[];
  generated_at: string;
}
