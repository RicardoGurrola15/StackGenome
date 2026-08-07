# StackGenome Project State

**Current Phase:** Phase 21 (Next Phase) - Public Release & Maintenance
**Status:** Alpha 20 Completed

## Core Capabilities Validated
- Deterministic language, framework, and infrastructure detection (via Territor pilot).
- Platform deduplication (Android/iOS correctly merged).
- Lockfile version resolution and scope propagation (pubspec.lock -> pubspec.yaml).
- Advanced tool and infra privacy (app_id/project_id excluded).
- Catalog recommendation engine with `minScoreThreshold` and `primaryLanguageBoost`.
- Web UI alpha with interactive D3 Glassmorphism stack graph.

## Unresolved Issues / Known Limitations
- The catalog is still limited to ~30 entries, mostly Go/JS/Python/Dart.
- The web viewer is purely client-side static export (requires providing the JSON file manually).
