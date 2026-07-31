import type { Env } from '../types/env.js';
import { errorResponse, jsonResponse } from '../errors.js';
import { queryResourceById } from '../db/queries.js';

export async function handleResourceById(_req: Request, env: Env, id: string): Promise<Response> {
  if (!id) {
    return errorResponse('NOT_FOUND', 'Resource id is required.', 404);
  }

  const resource = await queryResourceById(env.DB, id);
  if (!resource) {
    return errorResponse('NOT_FOUND', `Resource "${id}" not found.`, 404);
  }

  return jsonResponse(resource);
}
