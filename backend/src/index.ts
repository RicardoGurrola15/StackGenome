import type { Env } from './types/env.js';
import { routeRequest } from './router.js';
import { applyCors, handleOptions } from './middleware/cors.js';
import { checkRateLimit } from './middleware/rateLimit.js';
import { errorResponse } from './errors.js';

export default {
  async fetch(req: Request, env: Env, _ctx: ExecutionContext): Promise<Response> {
    // 1. CORS Preflight
    if (req.method.toUpperCase() === 'OPTIONS') {
      return handleOptions(req);
    }

    // 2. Rate Limiting
    const rlResponse = checkRateLimit(req, env);
    if (rlResponse) {
      return applyCors(req, rlResponse);
    }

    // 3. Routing
    try {
      const response = await routeRequest(req, env);
      return applyCors(req, response);
    } catch (e) {
      console.error('Fatal execution error:', e);
      return applyCors(req, errorResponse('INTERNAL_ERROR', 'Fatal internal error', 500));
    }
  },
};
