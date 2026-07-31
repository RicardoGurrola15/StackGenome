import type { Env } from '../types/env.js';
import { errorResponse, jsonResponse } from '../errors.js';
import { searchResources } from '../db/queries.js';

export async function handleSearch(req: Request, env: Env): Promise<Response> {
  const url = new URL(req.url);
  const q = (url.searchParams.get('q') ?? '').trim();
  const type = (url.searchParams.get('type') ?? '').trim();
  const ecosystem = (url.searchParams.get('ecosystem') ?? '').trim();

  if (q.length === 0 && type.length === 0 && ecosystem.length === 0) {
    return errorResponse('INVALID_PAYLOAD', 'At least one query parameter is required: q, type, or ecosystem.', 400);
  }
  if (q.length > 200) {
    return errorResponse('INVALID_PAYLOAD', 'Query parameter "q" must not exceed 200 characters.', 400);
  }

  const filters: { q?: string; type?: string; ecosystem?: string } = {};
  if (q) filters.q = q;
  if (type) filters.type = type;
  if (ecosystem) filters.ecosystem = ecosystem;

  const results = await searchResources(env.DB, filters);

  return jsonResponse({ results, count: results.length });
}
