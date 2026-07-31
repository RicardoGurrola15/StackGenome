import type { Env } from '../types/env.js';
import type { RecommendationsRequest, RecommendationsResponse } from '../types/fingerprint.js';
import { errorResponse, jsonResponse } from '../errors.js';
import { scoreAndRank } from '../ranking/engine.js';
import { queryResources } from '../db/queries.js';

export async function handleRecommendations(req: Request, env: Env): Promise<Response> {
  // 1. Content-Length guard (before parsing)
  const maxBytes = parseInt(env.MAX_PAYLOAD_BYTES, 10);
  const contentLength = req.headers.get('content-length');
  if (contentLength !== null && parseInt(contentLength, 10) > maxBytes) {
    return errorResponse('PAYLOAD_TOO_LARGE', `Payload exceeds maximum of ${maxBytes} bytes.`, 413);
  }

  // 2. Parse body with size enforcement
  let body: unknown;
  try {
    const text = await req.text();
    if (text.length > maxBytes) {
      return errorResponse('PAYLOAD_TOO_LARGE', `Payload exceeds maximum of ${maxBytes} bytes.`, 413);
    }
    body = JSON.parse(text);
  } catch {
    return errorResponse('INVALID_PAYLOAD', 'Request body must be valid JSON.', 400);
  }

  // 3. Validate schema_version
  const parsed = body as Partial<RecommendationsRequest>;
  if (parsed.schema_version !== '1.0.0') {
    return errorResponse(
      'INVALID_SCHEMA_VERSION',
      `Unsupported schema_version "${parsed.schema_version ?? ''}". Expected "1.0.0".`,
      422
    );
  }

  // 4. Validate fingerprint structure
  if (!parsed.fingerprint || !Array.isArray(parsed.fingerprint.nodes) || parsed.fingerprint.nodes.length === 0) {
    return errorResponse('INVALID_PAYLOAD', 'fingerprint.nodes must be a non-empty array.', 422);
  }

  const limit = Math.min(parsed.limit ?? 3, 10); // max 10, default 3

  // 5. Load resources from D1 (prepared statement)
  const resources = await queryResources(env.DB);

  // 6. Apply ranking engine
  const recommendations = scoreAndRank(parsed.fingerprint, resources, limit);

  // 7. Build response — fingerprint is NOT stored or logged
  const response: RecommendationsResponse = {
    catalog_snapshot_id: 'staging-v1',
    ranking_version: env.RANKING_VERSION,
    recommendations,
    warnings: [],
    generated_at: new Date().toISOString(),
  };

  return jsonResponse(response);
}
