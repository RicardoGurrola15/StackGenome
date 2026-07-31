import type { Env } from './types/env.js';
import { handleHealth } from './handlers/health.js';
import { handleSchemaFingerprint } from './handlers/schema.js';
import { handleRecommendations } from './handlers/recommendations.js';
import { handleResourceById } from './handlers/resources.js';
import { handleSearch } from './handlers/search.js';
import { errorResponse } from './errors.js';

export async function routeRequest(req: Request, env: Env): Promise<Response> {
  const url = new URL(req.url);
  const path = url.pathname;
  const method = req.method.toUpperCase();

  try {
    if (path === '/v1/health' && method === 'GET') {
      return await handleHealth(req, env);
    }
    
    if (path === '/v1/schema/fingerprint' && method === 'GET') {
      return await handleSchemaFingerprint(req, env);
    }

    if (path === '/v1/recommendations' && method === 'POST') {
      return await handleRecommendations(req, env);
    }

    if (path === '/v1/search' && method === 'GET') {
      return await handleSearch(req, env);
    }

    const resourceMatch = path.match(/^\/v1\/resources\/(.+)$/);
    if (resourceMatch && method === 'GET') {
      return await handleResourceById(req, env, resourceMatch[1]);
    }

    return errorResponse('NOT_FOUND', `Route ${method} ${path} not found.`, 404);
  } catch (err) {
    console.error('Unhandled route error:', err);
    return errorResponse('INTERNAL_ERROR', 'An unexpected internal error occurred.', 500);
  }
}
