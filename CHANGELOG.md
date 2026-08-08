# Changelog

All notable changes to StackGenome are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased — v0.2.0-alpha.1]

### Added
- **Context-aware recommendations**: Introduced `requires_context` field in the catalog. Tools like Dart Frog (a Dart backend framework) now require a `has_backend` signal to be recommended, preventing irrelevant suggestions in pure mobile Flutter projects.
- **Project Signals**: `BuildContext` now derives `has_backend`, `has_mobile`, and `has_firebase` signals from the project graph for use in context-gating.
- **Improved CLI output**: Languages are now deduplicated and organized into "Primary" vs "Also Detected" sections. Platforms, Infrastructure, and Tools now appear as separate labeled sections.
- **Better error messages**: Invalid directories now produce a clear, actionable error (not a panic).
- **Improved `--help`**: The `analyze` command now shows usage, flag descriptions, and concrete examples.
- **`-remote` fallback**: If the remote Cloudflare API is unavailable, the CLI now falls back to the local catalog gracefully.
- **Web build CI job**: Added `test-web` job to `ci.yml` that verifies `tsc` and `npm run build` pass.
- **Fixed backend CI**: Resolved `npm EACCES /Volumes` error in GitHub Actions by overriding `npm_config_cache`.
- **Catalog audit documents**: `docs/audits/CATALOG_ALPHA_READINESS_2026-08.md` and `docs/audits/CI_CD_REALITY_CHECK_2026-08.md`.
- **Public Alpha Smoke Test**: `docs/testing/PUBLIC_ALPHA_SMOKE_TEST.md`.
- **ADR-0010**: Documented decision to use Cloudflare Pages (vs Workers Static Assets) for frontend hosting.
- **Version injection fixed** in `build-release.sh`: now uses `git describe --tags` as source of truth.

### Changed
- `README.md`: Corrected absolute claims ("100% offline" → accurate technical description).
- `README.md`: Updated phase label from "Alpha (Fase 15)" to "Public Alpha / Pre-Release (Fase 21)".
- Catalog: Dart Frog entry updated with `category: "backend"` and `requires_context: ["has_backend"]`.

### Fixed
- Language deduplication in CLI output: multiple C/C++ and JavaScript nodes no longer appear as separate lines.
- `PrimaryLanguage` tie-breaking: when multiple languages share `1.0` confidence, the one with the most evidence files wins — ensuring Dart/Flutter is correctly identified as primary in Flutter projects.
- Desktop boilerplate excluded from walker: `windows/runner/`, `linux/flutter/` etc. are now in `vendoredPathPrefixes`.
- Evidences sorted by path in `ToDTO()` for deterministic JSON output.

---

## [alpha-20] — 2026-08-07 (internal)

### Added
- Glassmorphism Web UI with interactive D3 graph (zoom, drag, clustering).
- Firebase and Shorebird detection for Flutter projects.
- `pubspec.lock` version resolution merged into `pubspec.yaml` dependency nodes.
- Cloudflare D1 migration: ecosystem and Flutter-specific columns.
- Expanded local catalog with FVM, Melos, FlutterGen, Shorebird, Patrol, Leak Tracker, Dart Frog.
- `minScoreThreshold` (0.30) and `primaryLanguageBoost` (0.15) in scoring engine.

### Fixed
- `DeduplicatePlatforms` now sorts node IDs before iteration (deterministic).
- Excluded `ios/Pods/` and related CocoaPods paths from walker (F-001).
- Platform nodes properly deduplicated (F-002).
- Dart dependency versions and scopes extracted correctly (F-004, F-005).

---

## [alpha-18] — 2026-07-31 (internal)

### Added
- Static export Next.js Web Alpha.
- `ResourceModal` component for catalog detail view.
- D3.js stack graph with zoom/drag.
