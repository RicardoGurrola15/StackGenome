# StackGenome

**StackGenome** analyzes software projects and builds a normalized technology graph. It identifies languages, dependencies, tools, platforms, and infrastructure — then recommends relevant, well-maintained tools your project doesn't yet use.

> **Status: Public Alpha (v0.2.0-alpha.1)**
> The core CLI and analysis engine are stable for evaluation. Expect rough edges in output formatting and catalog coverage. Not recommended for production automation yet.

---

## What it does

```
📁 Your project
   ↓
🔬 StackGenome CLI (static analysis only — no code execution)
   ↓
🧬 ProjectGraph (languages, dependencies, tools, platforms, infra)
   ↓
✨ Recommendations from embedded catalog (offline by default)
```

StackGenome reads manifest files and configuration files — it does **not** execute your code, install dependencies, or modify your project in any way.

---

## Quick Start

### Install

Download the binary for your platform from [Releases](https://github.com/RicardoGurrola15/StackGenome/releases):

| Platform | File |
| :--- | :--- |
| macOS (Apple Silicon) | `stackgenome_darwin_arm64` |
| macOS (Intel) | `stackgenome_darwin_amd64` |
| Linux x86-64 | `stackgenome_linux_amd64` |
| Linux ARM64 | `stackgenome_linux_arm64` |
| Windows x86-64 | `stackgenome_windows_amd64.exe` |

```bash
# macOS / Linux
chmod +x stackgenome_darwin_arm64
mv stackgenome_darwin_arm64 /usr/local/bin/stackgenome
```

Verify the checksum from `checksums.txt` before running any downloaded binary.

### Usage

```bash
# Analyze the current directory
stackgenome analyze .

# Analyze a specific project
stackgenome analyze /path/to/your/project

# Get tool recommendations (offline, no network required)
stackgenome analyze -recommend /path/to/your/project

# Export full graph as JSON
stackgenome analyze -json . > report.json

# Check version
stackgenome version
```

---

## Supported Ecosystems

Analysis depth varies by ecosystem:

| Ecosystem | Language Detection | Dependency Extraction | Lock File Resolution |
| :--- | :---: | :---: | :---: |
| Go | ✅ | ✅ (`go.mod`, `go.work`) | ✅ |
| Node.js / TypeScript | ✅ | ✅ (`package.json`) | ✅ (`package-lock.json`) |
| Python | ✅ | ✅ (`requirements.txt`, `pyproject.toml`) | — |
| Rust | ✅ | ✅ (`Cargo.toml`) | ✅ (`Cargo.lock`) |
| Dart / Flutter | ✅ | ✅ (`pubspec.yaml`) | ✅ (`pubspec.lock`) |
| Java / JVM | ✅ | Partial | — |
| Swift / Objective-C | ✅ | — | — |
| C / C++ | Detected | — | — |

Projects using multiple ecosystems are fully supported.

---

## Privacy

StackGenome is built **local-first**:

- All analysis runs on your machine.
- No code, secrets, or file contents are read — only manifest and configuration files.
- The local catalog (`-recommend`) works entirely offline; no network calls are made.
- The optional `-remote` flag sends only a sanitized metadata fingerprint (no paths, versions, or credentials) to the Cloudflare API for updated recommendations. This requires your explicit opt-in.

The privacy model has been validated against a set of known sensitive fields. This is not a security guarantee for all possible inputs.

---

## Limitations (Alpha)

- **Catalog size**: ~27 curated tools. Ecosystems like iOS Native, Android Native, and DevOps-specific tooling have limited recommendations.
- **No runtime analysis**: StackGenome only reads static files. It cannot detect dynamically loaded dependencies.
- **Windows ARM64**: Not currently included in release binaries.
- **Determinism**: Results are deterministic for identical CLI version, catalog, and input files. Behavior may change between releases.

---

## Web Viewer

A browser-based graph viewer for `report.json` files is available at:
> *(URL to be published — Cloudflare Pages deployment in progress)*

The viewer runs entirely in your browser. JSON files are never uploaded.

---

## Documentation

| Document | Description |
| :--- | :--- |
| [CHANGELOG.md](CHANGELOG.md) | Release history |
| [CONTRIBUTING.md](CONTRIBUTING.md) | How to contribute |
| [SECURITY.md](SECURITY.md) | Security policy and reporting |
| [docs/](docs/) | Full technical documentation |

---

## For Developers & Agents

Read in order:
1. `AGENTS.md`
2. `docs/00_INDEX.md`
3. `.project/CURRENT_PHASE.md`

Do not implement phases beyond the currently authorized one.

---

## License

[MIT](LICENSE)
