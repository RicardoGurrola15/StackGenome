import fs from 'node:fs';

const catalogRaw = fs.readFileSync('../internal/catalog/catalog.json', 'utf8');
const catalog = JSON.parse(catalogRaw);

const items = catalog.entries;

const timestamp = new Date().toISOString();
const sqlLines = [];
sqlLines.push('-- Migration 0002: Seed staging resources');
sqlLines.push('');
sqlLines.push(`INSERT INTO resources (id, type, canonical_name, summary, canonical_url, ecosystem, infra_targets, status, updated_at) VALUES`);

for (let i = 0; i < items.length; i++) {
  const item = items[i];
  const isLast = i === items.length - 1;
  const ecosystem = JSON.stringify(item.targets.languages);
  const infra_targets = JSON.stringify(item.targets.infra);
  
  // escape single quotes
  const safeName = item.name.replace(/'/g, "''");
  const safeDesc = item.description.replace(/'/g, "''");
  
  sqlLines.push(`('${item.id}', 'tool', '${safeName}', '${safeDesc}', '${item.url}', '${ecosystem}', '${infra_targets}', 'active', '${timestamp}')${isLast ? ';' : ','}`);
}

sqlLines.push('');
fs.writeFileSync('migrations/0002_seed_staging.sql', sqlLines.join('\n'));
console.log('Successfully generated migrations/0002_seed_staging.sql');
