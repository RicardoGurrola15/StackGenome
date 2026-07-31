import type { Env } from '../types/env.js';
import { jsonResponse } from '../errors.js';

export async function handleHealth(_req: Request, env: Env): Promise<Response> {
  return jsonResponse({
    status: 'ok',
    environment: env.ENVIRONMENT,
    ranking_version: env.RANKING_VERSION,
    timestamp: new Date().toISOString(),
    version: '0.1.0',
  });
}
