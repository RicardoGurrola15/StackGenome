import { describe, it, expect } from 'vitest';
import { scoreAndRank } from '../src/ranking/engine.js';
import type { CatalogEntry } from '../src/db/queries.js';
import type { FingerprintV1 } from '../src/types/fingerprint.js';
import { WEIGHTS } from '../src/ranking/weights.js';

describe('Ranking Engine', () => {
  const dummyCatalog: CatalogEntry[] = [
    {
      id: 'tool:go-linter',
      name: 'Go Linter',
      description: 'Linter for Go',
      url: undefined,
      targets: { languages: [], infra: [], frameworks: [] },
      ecosystem: ['go'],
    },
    {
      id: 'tool:docker-scanner',
      name: 'Docker Scanner',
      description: 'Scans Docker images',
      url: undefined,
      targets: { languages: [], infra: ['docker'], frameworks: [] },
      ecosystem: [],
    },
    {
      id: 'tool:generic-tool',
      name: 'Generic Tool',
      description: 'Works for anything',
      url: undefined,
      targets: { languages: [], infra: [], frameworks: [] },
      ecosystem: [],
    }
  ];

  it('ranks by language match correctly', () => {
    const fp: FingerprintV1 = {
      nodes: [
        { id: 'lang:go', type: 'language', name: 'Go', confidence: 1 },
      ],
      edges: []
    };

    const recs = scoreAndRank(fp, dummyCatalog, 3);
    
    // Go Linter should match language + novelty
    expect(recs[0].id).toBe('tool:go-linter');
    expect(recs[0].score).toBe(WEIGHTS.language + WEIGHTS.novelty);
    expect(recs[0].reasons).toContain('ecosistema compatible: go');
  });

  it('ranks generic tools as fallback', () => {
    const fp: FingerprintV1 = {
      nodes: [
        { id: 'lang:python', type: 'language', name: 'Python', confidence: 1 },
      ],
      edges: []
    };

    const recs = scoreAndRank(fp, dummyCatalog, 3);
    
    // Python doesn't match go-linter or docker-scanner, so generic-tool wins
    expect(recs[0].id).toBe('tool:generic-tool');
    expect(recs[0].score).toBe(WEIGHTS.novelty);
    expect(recs[0].reasons).toContain('herramienta de uso general');
  });

  it('ignores tools already present in the fingerprint', () => {
    const fp: FingerprintV1 = {
      nodes: [
        { id: 'lang:go', type: 'language', name: 'Go', confidence: 1 },
        { id: 'tool:go-linter', type: 'tool', name: 'Go Linter', confidence: 1 },
      ],
      edges: []
    };

    const recs = scoreAndRank(fp, dummyCatalog, 3);
    
    // Go Linter matches language but loses novelty
    // Let's check score
    const linterRec = recs.find(r => r.id === 'tool:go-linter');
    expect(linterRec).toBeDefined();
    expect(linterRec!.score).toBe(WEIGHTS.language); // NO novelty
  });
});
