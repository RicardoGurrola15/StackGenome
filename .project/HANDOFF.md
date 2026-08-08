# Agent Handoff

## State at Phase 21A Completion

The repository is clean, all tests pass, and the scoring engine has been made context-aware.
The CI has known issues that must be fixed before the first real tag is pushed.

## For Next Agent (Phase 21B onward)

### Priority 1 — Fix CI before creating any tag
The backend CI job (`test-backend` in `ci.yml`) fails with:
```
npm error EACCES: permission denied, mkdir '/Volumes'
```
Root cause: `npm` on the GitHub Actions runner (Ubuntu) is picking up the macOS-style cache path from `package.json` or `.npmrc`. Fix: ensure `cache-dependency-path` and any `npm config` does not reference `/Volumes`.

### Priority 2 — Add web build job to ci.yml
No job currently builds or lints the Next.js frontend.

### Priority 3 — Create v0.2.0-alpha.1 tag
Once CI is green, create the tag. The `release.yml` workflow will fire automatically and produce the GitHub Release with all 5 binaries and checksums.txt.

### Priority 4 — Deploy to Cloudflare Pages
Decision: ADR-0010 chose Cloudflare Pages (not Workers Static Assets).
Command: `cd web && npm run build && wrangler pages deploy out --project-name stackgenome`

### Do NOT
- Force-push or rewrite history.
- Push `.engram/` directory.
- Push Territor analysis outputs.
- Declare v1.0.0.
