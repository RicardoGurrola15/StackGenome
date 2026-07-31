import type { Env } from '../types/env.js';
import { jsonResponse } from '../errors.js';

// Describes the fingerprint schema accepted by POST /v1/recommendations
const FINGERPRINT_SCHEMA = {
  schema_version: '1.0.0',
  description: 'StackGenome ProjectGraph fingerprint — metadata-only projection.',
  fields: {
    nodes: {
      type: 'array',
      required: true,
      items: {
        id: 'string (required)',
        type: 'string (required) — language | framework | infrastructure | cicd | editor | platform | dependency | workspace',
        name: 'string (required)',
        version: 'string (optional)',
        confidence: 'number 0.0–1.0 (required)',
        properties: 'object (optional) — allowlisted keys only: type, framework, scope, workspace, purl_type',
      },
    },
    edges: {
      type: 'array',
      required: false,
      items: {
        source_id: 'string',
        target_id: 'string',
        relation: 'string — contains | depends_on | uses',
      },
    },
  },
  not_accepted: [
    'evidences (filesystem paths)',
    'absolute paths',
    'secret values',
    'source code fragments',
    'personal identifiers',
  ],
  max_payload_bytes: 32768,
};

export async function handleSchemaFingerprint(_req: Request, _env: Env): Promise<Response> {
  return jsonResponse(FINGERPRINT_SCHEMA);
}
