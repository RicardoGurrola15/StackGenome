import type { FingerprintV1, RecommendationResult } from '../types/fingerprint.js';
import { WEIGHTS } from './weights.js';
import type { CatalogEntry } from '../db/queries.js';

export interface EntryTarget {
  languages: string[];
  frameworks: string[];
  infra: string[];
}

export function scoreAndRank(fingerprint: FingerprintV1, resources: CatalogEntry[], limit: number): RecommendationResult[] {
  const ctx = buildContext(fingerprint);
  const scored = [];

  for (const res of resources) {
    if (!entryIsRelevant(res, ctx)) continue;

    const { score, reasons } = computeScore(res, ctx);
    if (score > 0) {
      const result: RecommendationResult = {
        id: res.id,
        name: res.name,
        description: res.description,
        score: Math.round(score * 100) / 100, // round to 2 decimals
        reasons,
      };
      if (res.url) result.url = res.url;
      scored.push(result);
    }
  }

  // Deterministic sort: score DESC, then ID ASC
  scored.sort((a, b) => {
    if (a.score !== b.score) return b.score - a.score;
    return a.id.localeCompare(b.id);
  });

  return scored.slice(0, limit);
}

interface ProjectContext {
  languages: Set<string>;
  infra: Set<string>;
  frameworks: Set<string>;
  nodeIds: Set<string>;
}

function buildContext(fingerprint: FingerprintV1): ProjectContext {
  const ctx: ProjectContext = {
    languages: new Set(),
    infra: new Set(),
    frameworks: new Set(),
    nodeIds: new Set(),
  };

  for (const n of fingerprint.nodes) {
    ctx.nodeIds.add(n.id);
    const nameLower = n.name.toLowerCase();
    switch (n.type) {
      case 'language':
        ctx.languages.add(nameLower);
        if (nameLower.includes('/')) {
          const parts = nameLower.split('/');
          for (const p of parts) {
            ctx.languages.add(p.trim());
          }
        }
        break;
      case 'infrastructure':
        ctx.infra.add(nameLower);
        break;
      case 'framework':
        ctx.frameworks.add(nameLower);
        break;
    }
  }

  return ctx;
}

function entryIsRelevant(e: CatalogEntry, ctx: ProjectContext): boolean {
  if (e.ecosystem.length === 0 && e.targets.infra.length === 0 && e.targets.frameworks.length === 0) {
    return true;
  }
  if (e.ecosystem.length === 1 && e.ecosystem[0] === '*') {
    return true;
  }
  for (const eco of e.ecosystem) {
    if (ctx.languages.has(eco.toLowerCase())) return true;
  }
  for (const infra of e.targets.infra) {
    if (ctx.infra.has(infra.toLowerCase())) return true;
  }
  for (const fw of e.targets.frameworks) {
    if (ctx.frameworks.has(fw.toLowerCase())) return true;
  }
  return false;
}

function computeScore(e: CatalogEntry, ctx: ProjectContext): { score: number; reasons: string[] } {
  let score = 0;
  const reasons: string[] = [];

  const matchedLangs = e.ecosystem.filter(l => l !== '*' && ctx.languages.has(l.toLowerCase()));
  if (matchedLangs.length > 0) {
    score += WEIGHTS.language;
    reasons.push(`ecosistema compatible: ${matchedLangs.join(', ')}`);
  }

  const matchedInfra = e.targets.infra.filter(i => ctx.infra.has(i.toLowerCase()));
  if (matchedInfra.length > 0) {
    score += WEIGHTS.infra;
    reasons.push(`infraestructura compatible: ${matchedInfra.join(', ')}`);
  }

  const matchedFw = e.targets.frameworks.filter(f => ctx.frameworks.has(f.toLowerCase()));
  if (matchedFw.length > 0) {
    score += WEIGHTS.framework;
    reasons.push(`framework compatible: ${matchedFw.join(', ')}`);
  }

  if (!ctx.nodeIds.has(e.id)) {
    score += WEIGHTS.novelty;
    reasons.push('herramienta no detectada aún en el proyecto');
  }

  if ((e.ecosystem.length === 0 || (e.ecosystem.length === 1 && e.ecosystem[0] === '*')) && e.targets.infra.length === 0 && e.targets.frameworks.length === 0) {
    if (score <= WEIGHTS.novelty) {
      score = WEIGHTS.novelty;
      reasons.push('herramienta de uso general');
    }
  }

  return { score, reasons };
}
