-- Migration 0003: Populate ecosystem column and add Flutter/Dart tools

-- 1. Populate ecosystem for existing generic tools (fallback to "*")
UPDATE resources SET ecosystem = '["*"]' WHERE type = 'tool' AND ecosystem IS NULL;

-- (The existing 0002_seed_staging.sql inserted raw strings but maybe they had no ecosystem JSON)
-- We'll just update all known ecosystems manually based on ID.
UPDATE resources SET ecosystem = '["go", "python", "java", "node", "dotnet", "php", "ruby"]' WHERE id = 'tool:opentelemetry';
UPDATE resources SET ecosystem = '["go"]' WHERE id = 'tool:golangci-lint';
UPDATE resources SET ecosystem = '["python"]' WHERE id = 'tool:ruff';
UPDATE resources SET ecosystem = '["node"]' WHERE id = 'tool:biome';
UPDATE resources SET ecosystem = '["go", "python", "node", "ruby", "rust"]' WHERE id = 'tool:mise';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:trivy';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:cosign';
UPDATE resources SET ecosystem = '["go", "python", "node", "rust", "java", "php", "ruby", "dotnet"]' WHERE id = 'tool:renovate';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:act';
UPDATE resources SET ecosystem = '["go", "python", "node", "rust", "java"]' WHERE id = 'tool:earthly';
UPDATE resources SET ecosystem = '["go", "java", "python", "node", "dotnet", "rust"]' WHERE id = 'tool:testcontainers';
UPDATE resources SET ecosystem = '["go", "python", "java", "node"]' WHERE id = 'tool:buf';
UPDATE resources SET ecosystem = '["go", "java", "python", "node", "dotnet", "php"]' WHERE id = 'tool:temporal';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:sops';
UPDATE resources SET ecosystem = '["go", "python", "node"]' WHERE id = 'tool:dagger';
UPDATE resources SET ecosystem = '["go", "rust", "python", "node", "ruby"]' WHERE id = 'tool:nix';
UPDATE resources SET ecosystem = '["node", "go", "python", "java"]' WHERE id = 'tool:k6';
UPDATE resources SET ecosystem = '["rust"]' WHERE id = 'tool:cargo-nextest';
UPDATE resources SET ecosystem = '["go"]' WHERE id = 'tool:sqlc';
UPDATE resources SET ecosystem = '["go", "rust", "python", "node", "ruby", "cpp"]' WHERE id = 'tool:just';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:atuin';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:checkov';
UPDATE resources SET ecosystem = '["go", "java", "node", "python", "dotnet", "ruby"]' WHERE id = 'tool:pact';
UPDATE resources SET ecosystem = '["go", "node", "python", "rust"]' WHERE id = 'tool:taskfile';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:grype';
UPDATE resources SET ecosystem = '["go", "rust", "python", "node", "java", "cpp", "ruby", "php"]' WHERE id = 'tool:typos';
UPDATE resources SET ecosystem = '["go", "node", "python", "java", "rust"]' WHERE id = 'tool:mermaid';
UPDATE resources SET ecosystem = '["go", "java", "swift", "node", "python"]' WHERE id = 'tool:pkl';
UPDATE resources SET ecosystem = '["*"]' WHERE id = 'tool:vcluster';
UPDATE resources SET ecosystem = '["go", "node", "python", "ruby", "rust"]' WHERE id = 'tool:lefthook';

-- 2. Insert new Flutter/Dart tools
INSERT INTO resources (id, type, canonical_name, summary, canonical_url, ecosystem, infra_targets, status, updated_at)
VALUES 
  ('tool:fvm', 'tool', 'FVM (Flutter Version Management)', 'Gestor de versiones de Flutter por proyecto para asegurar builds consistentes.', 'https://fvm.app', '["dart", "flutter"]', '[]', 'active', CURRENT_TIMESTAMP),
  ('tool:melos', 'tool', 'Melos', 'Herramienta CLI para gestionar monorepos (workspaces) en proyectos Dart y Flutter.', 'https://melos.invertase.dev', '["dart", "flutter"]', '[]', 'active', CURRENT_TIMESTAMP),
  ('tool:flutter_gen', 'tool', 'FlutterGen', 'Generador de código para assets, fuentes, y colores de Flutter — evita errores de tipeo.', 'https://pub.dev/packages/flutter_gen', '["dart", "flutter"]', '[]', 'active', CURRENT_TIMESTAMP),
  ('tool:very_good_cli', 'tool', 'Very Good CLI', 'Herramienta de scaffolding para generar proyectos Flutter modernos, escalables y con tests.', 'https://cli.vgv.dev', '["dart", "flutter"]', '[]', 'active', CURRENT_TIMESTAMP),
  ('tool:shorebird', 'tool', 'Shorebird', 'Actualizaciones OTA (Over-The-Air) para apps Flutter, empujando parches instantáneos.', 'https://shorebird.dev', '["dart", "flutter"]', '[]', 'active', CURRENT_TIMESTAMP),
  ('tool:patrol', 'tool', 'Patrol', 'Framework avanzado de UI testing para Flutter, con soporte para interactuar con SO nativo.', 'https://patrol.leancode.co', '["dart", "flutter"]', '[]', 'active', CURRENT_TIMESTAMP),
  ('tool:dart_frog', 'tool', 'Dart Frog', 'Framework rápido y minimalista para construir backends y APIs en Dart.', 'https://dartfrog.vgv.dev', '["dart"]', '[]', 'active', CURRENT_TIMESTAMP),
  ('tool:leak_tracker', 'tool', 'Leak Tracker', 'Herramienta de detección de memory leaks para Dart y Flutter.', 'https://pub.dev/packages/leak_tracker', '["dart", "flutter"]', '[]', 'active', CURRENT_TIMESTAMP)
ON CONFLICT(id) DO UPDATE SET 
  ecosystem = excluded.ecosystem,
  summary = excluded.summary;
