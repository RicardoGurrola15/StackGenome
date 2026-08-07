import type { EntryTarget } from '../ranking/engine.js';

export interface ResourceRow {
  id: string;
  type: string;
  canonical_name: string;
  summary: string;
  canonical_url: string | null;
  ecosystem: string; // JSON
  infra_targets: string; // JSON
}

export interface CatalogEntry {
  id: string;
  name: string;
  description: string;
  url: string | undefined;
  targets: EntryTarget;
  ecosystem: string[];
}

export async function queryResources(db: D1Database): Promise<CatalogEntry[]> {
  const { results } = await db.prepare(
    `SELECT id, type, canonical_name, summary, canonical_url, ecosystem, infra_targets 
     FROM resources 
     WHERE status = 'active'`
  ).all<ResourceRow>();

  if (!results) return [];

  return results.map(row => ({
    id: row.id,
    name: row.canonical_name,
    description: row.summary,
    url: row.canonical_url || undefined,
    ecosystem: JSON.parse(row.ecosystem || '[]'),
    targets: {
      languages: [], // we will deprecate languages in favor of ecosystem
      infra: JSON.parse(row.infra_targets || '[]'),
      frameworks: [], // Simplified for MVP
    }
  }));
}

export async function queryResourceById(db: D1Database, id: string): Promise<ResourceRow | null> {
  const result = await db.prepare(
    `SELECT id, type, canonical_name, summary, canonical_url, ecosystem, infra_targets 
     FROM resources 
     WHERE id = ?`
  ).bind(id).first<ResourceRow>();
  
  return result;
}

export async function searchResources(db: D1Database, filters: { q?: string; type?: string; ecosystem?: string }): Promise<ResourceRow[]> {
  let query = `
    SELECT r.id, r.type, r.canonical_name, r.summary, r.canonical_url, r.ecosystem, r.infra_targets
    FROM resources r
  `;
  const conditions: string[] = [];
  const params: string[] = [];

  if (filters.q) {
    query += ` JOIN resources_fts fts ON r.id = fts.id`;
    conditions.push(`fts.resources_fts MATCH ?`);
    params.push(filters.q);
  }

  if (filters.type) {
    conditions.push(`r.type = ?`);
    params.push(filters.type);
  }

  if (filters.ecosystem) {
    // Basic JSON array search (D1/SQLite json_each could be used, but LIKE is enough for exact tag match in MVP)
    conditions.push(`r.ecosystem LIKE ?`);
    params.push(`%${filters.ecosystem}%`);
  }

  if (conditions.length > 0) {
    query += ` WHERE ` + conditions.join(' AND ');
  }

  query += ` LIMIT 50`;

  const { results } = await db.prepare(query).bind(...params).all<ResourceRow>();
  return results ?? [];
}
