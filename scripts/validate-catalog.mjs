#!/usr/bin/env node
// scripts/validate-catalog.mjs
// Validates the catalog seed migration for completeness and correctness.
// Usage: node scripts/validate-catalog.mjs

import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __dirname = dirname(fileURLToPath(import.meta.url));
const seedPath = join(__dirname, '../backend/migrations/0002_seed_staging.sql');
const content = readFileSync(seedPath, 'utf8');

// ─── Known taxonomy ──────────────────────────────────────────────────────────
const KNOWN_TYPES = new Set(['tool', 'library', 'service', 'platform']);
const KNOWN_ECOSYSTEMS = new Set([
  'go', 'python', 'node', 'rust', 'java', 'dotnet', 'php', 'ruby',
  'swift', 'kotlin', 'cpp', 'c', 'dart', 'scala', 'elixir', 'haskell',
]);
const KNOWN_INFRA = new Set([
  'docker', 'kubernetes', 'terraform', 'github-actions', 'aws', 'gcp',
  'azure', 'linux', 'macos', 'windows',
]);

const URL_REGEX = /^https?:\/\/.+/;

// ─── Parse INSERT rows ────────────────────────────────────────────────────────
// Matches: ('id', 'type', 'canonical_name', 'summary', 'url', 'ecosystem_json', 'infra_json', 'status', 'updated_at')
const ROW_REGEX = /\('([^']+)',\s*'([^']+)',\s*'([^']+)',\s*'([^']+)',\s*'([^']*)',\s*'([^']*)',\s*'([^']*)',\s*'([^']+)',\s*'([^']*)'\)/g;

const errors = [];
const warnings = [];
const ids = new Set();

let match;
let rowCount = 0;

while ((match = ROW_REGEX.exec(content)) !== null) {
  rowCount++;
  const [, id, type, name, summary, url, ecosystemRaw, infraRaw, status] = match;

  // 1. Required fields
  if (!id)      errors.push(`Row ${rowCount}: 'id' está vacío.`);
  if (!type)    errors.push(`Row ${rowCount} (${id}): 'type' está vacío.`);
  if (!name)    errors.push(`Row ${rowCount} (${id}): 'canonical_name' está vacío.`);
  if (!summary) errors.push(`Row ${rowCount} (${id}): 'summary' está vacío.`);

  // 2. Duplicate check
  if (ids.has(id)) {
    errors.push(`Row ${rowCount}: ID duplicado detectado → '${id}'.`);
  }
  ids.add(id);

  // 3. Type taxonomy
  if (!KNOWN_TYPES.has(type)) {
    errors.push(`Row ${rowCount} (${id}): Tipo desconocido → '${type}'. Permitidos: ${[...KNOWN_TYPES].join(', ')}.`);
  }

  // 4. URL validation
  if (url && !URL_REGEX.test(url)) {
    errors.push(`Row ${rowCount} (${id}): URL inválida → '${url}'.`);
  }
  if (!url) {
    warnings.push(`Row ${rowCount} (${id}): Sin URL canónica.`);
  }

  // 5. Ecosystem taxonomy
  let ecosystems = [];
  try {
    ecosystems = JSON.parse(ecosystemRaw || '[]');
  } catch {
    errors.push(`Row ${rowCount} (${id}): JSON de ecosistema inválido → '${ecosystemRaw}'.`);
  }
  for (const eco of ecosystems) {
    if (!KNOWN_ECOSYSTEMS.has(eco)) {
      warnings.push(`Row ${rowCount} (${id}): Ecosistema no reconocido → '${eco}'.`);
    }
  }

  // 6. Infra taxonomy
  let infras = [];
  try {
    infras = JSON.parse(infraRaw || '[]');
  } catch {
    errors.push(`Row ${rowCount} (${id}): JSON de infra inválido → '${infraRaw}'.`);
  }
  for (const infra of infras) {
    if (!KNOWN_INFRA.has(infra)) {
      warnings.push(`Row ${rowCount} (${id}): Infra no reconocida → '${infra}'.`);
    }
  }

  // 7. Status check
  if (!['active', 'deprecated', 'archived'].includes(status)) {
    errors.push(`Row ${rowCount} (${id}): Status inválido → '${status}'.`);
  }

  // 8. Summary minimum length
  if (summary.length < 20) {
    warnings.push(`Row ${rowCount} (${id}): Summary muy corto (${summary.length} chars).`);
  }
}

// ─── Report ───────────────────────────────────────────────────────────────────
console.log(`\n=== StackGenome Catalog Validator ===`);
console.log(`Archivo: ${seedPath}`);
console.log(`Total de recursos: ${rowCount}`);
console.log(`IDs únicos: ${ids.size}`);

if (warnings.length > 0) {
  console.log(`\n⚠️  Warnings (${warnings.length}):`);
  warnings.forEach(w => console.log(`  - ${w}`));
}

if (errors.length > 0) {
  console.error(`\n❌ Errores (${errors.length}):`);
  errors.forEach(e => console.error(`  - ${e}`));
  process.exit(1);
} else {
  console.log(`\n✅ Catálogo válido. Sin errores críticos.`);
}
