# StackGenome Project State

**Current Phase:** Phase 21A (Pre-Release Consistency Audit) — COMPLETED
**Next Phase:** Phase 21B–21G (Public Alpha Release)
**Status:** `PUBLIC_ALPHA_READY_WITH_CONDITIONS`

## Core Capabilities (Validated)

| Component | Status |
| :--- | :--- |
| CLI Core | READY FOR PUBLIC ALPHA |
| Static Analyzer | READY FOR PUBLIC ALPHA |
| ProjectGraph | READY FOR PUBLIC ALPHA |
| Privacy Model | READY WITH DOCUMENTED LIMITS |
| Offline Recommendations | READY FOR ALPHA |
| Remote Recommendations | CONFIGURED — Requires CLOUDFLARE_API_TOKEN in CI |
| Local Catalog | READY FOR ALPHA |
| Cloud Catalog (D1) | CONFIGURED — Not yet deployed to staging |
| Backend (Cloudflare Worker) | CONFIGURED — Broken in CI (npm /Volumes permissions bug) |
| Web Alpha UI | READY FOR DEPLOYMENT (Cloudflare Pages) |
| CI/CD (Go tests) | CONFIGURED — Requires push to trigger verification |
| CI/CD (Backend tests) | BROKEN — npm EACCES /Volumes bug in GitHub Actions |
| CI/CD (Web build) | NOT IMPLEMENTED |
| Cross-Platform Builds | CONFIGURED — Script tested locally, all 5 targets build |
| Distribution (GitHub Release) | NOT PUBLICLY RELEASED YET |
| Production | NOT READY |
| v1.0.0 | NOT READY |

## Conditions for Public Alpha (Remaining — Phase 21B–21F)

1. Fix backend CI (npm /Volumes EACCES bug in ci.yml).
2. Add web build job to ci.yml.
3. Create a real `v0.2.0-alpha.1` tag to test and trigger the release pipeline.
4. Deploy frontend to Cloudflare Pages.
5. Write public-quality README.md.
6. Publish release with checksums and release notes.

## Key Decisions Made in Phase 21A

- **D-019**: Cloudflare Pages chosen for web frontend hosting (ADR-0010).
- **D-020**: `requires_context` gate implemented in catalog scorer to prevent context-inappropriate recommendations (e.g. Dart Frog not recommended to pure mobile Flutter apps).
- **Proposed version**: `v0.2.0-alpha.1` (no tags existed prior; this is the first public release candidate).

## Known Limitations (Alpha)

- Catalog has ~27 entries covering Go, Python, Rust, Node.js, Dart/Flutter and cross-cutting tools. Other ecosystems may be identified at the language level but will not produce deep package-level recommendations.
- Determinism is guaranteed for identical CLI version, configuration, catalog, and input sets within the verified scenarios.
- The privacy model follows metadata-only policy. No exposure of prohibited sensitive fields was detected in tests, but this is not a security guarantee for all possible future inputs.
