import { env, createExecutionContext, waitOnExecutionContext } from 'cloudflare:test';
import { describe, it, expect, beforeAll } from 'vitest';
import worker from '../src/index.js';
import type { Env } from '../src/types/env.js';

describe('StackGenome API', () => {
  beforeAll(async () => {
    const typedEnv = env as unknown as Env;
    await typedEnv.DB.prepare(`
      CREATE TABLE IF NOT EXISTS resources (
          id             TEXT PRIMARY KEY,
          type           TEXT NOT NULL,
          canonical_name TEXT NOT NULL,
          summary        TEXT NOT NULL,
          canonical_url  TEXT,
          ecosystem      TEXT NOT NULL DEFAULT '[]',
          infra_targets  TEXT NOT NULL DEFAULT '[]',
          status         TEXT NOT NULL DEFAULT 'active',
          updated_at     TEXT NOT NULL
      );
    `).run();

    await typedEnv.DB.prepare(`
      INSERT INTO resources (id, type, canonical_name, summary, ecosystem, infra_targets, updated_at) 
      VALUES ('tool:test', 'tool', 'Test Tool', 'A test tool', '["go"]', '[]', '2026-07-30')
    `).run();
  });

  it('GET /v1/health returns 200 OK', async () => {
    const request = new Request('http://example.com/v1/health');
    const ctx = createExecutionContext();
    const response = await worker.fetch(request, env as any, ctx);
    await waitOnExecutionContext(ctx);
    
    expect(response.status).toBe(200);
    const data = await response.json() as any;
    expect(data.status).toBe('ok');
    expect(data.version).toBe('0.1.0');
  });

  it('GET /v1/schema/fingerprint returns 200 OK', async () => {
    const request = new Request('http://example.com/v1/schema/fingerprint');
    const ctx = createExecutionContext();
    const response = await worker.fetch(request, env as any, ctx);
    await waitOnExecutionContext(ctx);
    
    expect(response.status).toBe(200);
    const data = await response.json() as any;
    expect(data.schema_version).toBe('1.0.0');
  });

  it('POST /v1/recommendations validates payload size', async () => {
    // Create a payload larger than 32KB
    const bigPayload = 'a'.repeat(33000);
    const request = new Request('http://example.com/v1/recommendations', {
      method: 'POST',
      headers: { 'content-length': '33000' },
      body: bigPayload
    });
    const ctx = createExecutionContext();
    
    // We override max payload for the test environment just in case
    const testEnv = { ...env, MAX_PAYLOAD_BYTES: '32768' };
    const response = await worker.fetch(request, testEnv as any, ctx);
    await waitOnExecutionContext(ctx);
    
    expect(response.status).toBe(413);
    const data = await response.json() as any;
    expect(data.error.code).toBe('PAYLOAD_TOO_LARGE');
  });

  it('POST /v1/recommendations works with valid fingerprint', async () => {
    const payload = {
      schema_version: '1.0.0',
      fingerprint: {
        nodes: [{ id: 'lang:go', type: 'language', name: 'Go', confidence: 1 }],
        edges: []
      },
      limit: 1
    };
    
    const request = new Request('http://example.com/v1/recommendations', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
    const ctx = createExecutionContext();
    
    const testEnv = { ...env, MAX_PAYLOAD_BYTES: '32768', RANKING_VERSION: 'v1' };
    const response = await worker.fetch(request, testEnv as any, ctx);
    await waitOnExecutionContext(ctx);
    
    expect(response.status).toBe(200);
    const data = await response.json() as any;
    expect(data.ranking_version).toBe('v1');
    expect(data.recommendations).toBeInstanceOf(Array);
    
    // Should contain our seeded 'tool:test'
    expect(data.recommendations.length).toBeGreaterThan(0);
    expect(data.recommendations[0].id).toBe('tool:test');
  });
});
