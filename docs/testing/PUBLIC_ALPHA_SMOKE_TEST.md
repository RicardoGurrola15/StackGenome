# Public Alpha Smoke Test

**Version:** v0.2.0-alpha.1 (candidate)
**Audience:** First-time user with no prior knowledge of StackGenome development.
**Goal:** Validate that StackGenome is installable and usable by an external developer using only the published artefacts.

> [!IMPORTANT]
> This test must be performed using a downloaded binary, not via `go run` or a source checkout.

---

## Step 1 — Obtain the Artefact

Download the binary for your platform from the GitHub Releases page:

```
https://github.com/RicardoGurrola15/StackGenome/releases
```

Select the correct artefact:

| Platform | File |
| :--- | :--- |
| macOS (Apple Silicon) | `stackgenome_darwin_arm64` |
| macOS (Intel) | `stackgenome_darwin_amd64` |
| Linux (x86-64) | `stackgenome_linux_amd64` |
| Linux (ARM64) | `stackgenome_linux_arm64` |
| Windows (x86-64) | `stackgenome_windows_amd64.exe` |

## Step 2 — Verify Checksum

Download `checksums.txt` from the same release and verify your binary:

```bash
# macOS / Linux
shasum -a 256 stackgenome_darwin_arm64
# Compare with the corresponding line in checksums.txt
```

> [!CAUTION]
> Do not use a binary whose checksum does not match.

## Step 3 — Make Executable & Test Invocation

```bash
chmod +x stackgenome_darwin_arm64
./stackgenome_darwin_arm64 --version
# Expected output: stackgenome v0.2.0-alpha.1

./stackgenome_darwin_arm64 --help
# Expected: usage text with available subcommands (analyze, version, completion)
```

## Step 4 — Analyze a Real Repository

Point the CLI at any software project directory on your machine:

```bash
./stackgenome_darwin_arm64 analyze /path/to/your/project
```

Expected output format:

```
🔬 StackGenome Analysis Complete
=================================

📦 Languages: ...
✨ Top Recommendations: ...
```

## Step 5 — Obtain Human-Readable Summary

The default output (without flags) produces a human-readable report:

```bash
./stackgenome_darwin_arm64 analyze /path/to/your/project
```

## Step 6 — Generate JSON Output

```bash
./stackgenome_darwin_arm64 analyze -json /path/to/your/project > project_graph.json
```

Verify the file is valid JSON with at least `nodes` and `edges` keys:

```bash
python3 -c "import json; d=json.load(open('project_graph.json')); print(len(d['nodes']), 'nodes')"
```

## Step 7 — Offline Recommendations

No network connection required:

```bash
./stackgenome_darwin_arm64 analyze -recommend /path/to/your/project
```

Expected: recommendations appear without any HTTP calls (verifiable via `--verbose` or network monitor).

## Step 8 — Open JSON in the Web Viewer

Navigate to the StackGenome web viewer (Cloudflare Pages URL — to be published).

1. Click **"Load JSON"** or drag the `project_graph.json` file.
2. The interactive stack graph should render.
3. Verify that no network upload occurs (all processing is client-side).

> [!NOTE]
> Until the web viewer is publicly deployed, this step can be tested locally by running `npm run dev` inside the `/web` directory.

## Step 9 — Delete the Report

```bash
rm project_graph.json
```

Confirm no data was cached or uploaded by StackGenome.

## Step 10 — Test Invalid Input

```bash
./stackgenome_darwin_arm64 analyze /nonexistent/path
```

Expected: a clear, non-panicking error message such as:

```
Error: directory does not exist: /nonexistent/path
```

The CLI must exit with a non-zero status code and must NOT panic.

---

## Pass / Fail Criteria

| Step | Pass Condition |
| :--- | :--- |
| 1 | Binary downloaded successfully |
| 2 | Checksum matches |
| 3 | `--version` prints version; `--help` shows usage |
| 4 | Analysis runs without crash |
| 5 | Human-readable output produced |
| 6 | Valid JSON with `nodes` key |
| 7 | Recommendations appear; no external network calls |
| 8 | Web graph renders from local file |
| 9 | File deleted; no side effects |
| 10 | Error is graceful; exit code ≠ 0 |

## Known Limitations in Alpha

- Language support is limited to the ecosystems listed in the support matrix.
- The web viewer must be visited separately — the CLI does not launch a browser automatically.
- Windows arm64 is not included in the current release.
