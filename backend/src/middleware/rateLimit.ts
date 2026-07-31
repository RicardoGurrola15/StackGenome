import type { Env } from '../types/env.js';
import { errorResponse } from '../errors.js';

// In-memory rate limiter. Ephemeral, resets on Worker restart/eviction.
// Sufficient for staging. For production, Cloudflare Rate Limiting rules or Durable Objects would be used.
const reqMap = new Map<string, { count: number; windowStart: number }>();

export function checkRateLimit(req: Request, env: Env): Response | null {
  const ip = req.headers.get('CF-Connecting-IP') || '127.0.0.1';
  const rpm = parseInt(env.RATE_LIMIT_RPM || '10', 10);
  const now = Date.now();
  const windowMs = 60 * 1000;

  let record = reqMap.get(ip);
  if (!record || now - record.windowStart > windowMs) {
    record = { count: 0, windowStart: now };
  }

  record.count++;
  reqMap.set(ip, record);

  // cleanup old entries occasionally
  if (Math.random() < 0.05) {
    for (const [key, val] of reqMap.entries()) {
      if (now - val.windowStart > windowMs) {
        reqMap.delete(key);
      }
    }
  }

  if (record.count > rpm) {
    return errorResponse('RATE_LIMITED', `Too many requests. Limit is ${rpm} per minute.`, 429);
  }

  return null;
}
