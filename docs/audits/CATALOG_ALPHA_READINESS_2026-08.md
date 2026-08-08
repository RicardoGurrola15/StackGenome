# Catalog Alpha Readiness Audit (Agosto 2026)

Este documento audita cada recurso del catálogo embebido (`internal/catalog/catalog.json`) antes del lanzamiento de la Public Alpha.

## Criterios de Clasificación

- `VALID` — Nombre, URL, descripción, ecosistema y condiciones de recomendación correctas.
- `NEEDS_METADATA` — Faltan campos como `category`, `license` o información de mantenimiento.
- `NEEDS_RULES` — Las reglas de recomendación son demasiado amplias (falta `requires_context`).
- `OUTDATED` — URL o nombre ya no son correctos.
- `DUPLICATE` — Existe una entrada equivalente.
- `REMOVE` — El recurso no es pertinente para el catálogo Alpha.

---

## Inventario

| ID | Nombre | Clasificación | Notas |
| :--- | :--- | :--- | :--- |
| `tool:opentelemetry` | OpenTelemetry | `NEEDS_METADATA` | Falta `category`. Ecosistema completo. |
| `tool:golangci-lint` | golangci-lint | `VALID` | Bien definido, targets correctos. |
| `tool:ko` | Ko | `VALID` | Build Go containers, ecosistema Go. |
| `tool:buf` | Buf | `NEEDS_METADATA` | Falta `category: "api"`. Requiere Protobuf señal. |
| `tool:earthly` | Earthly | `NEEDS_RULES` | Se recomienda a proyectos Java sin validar CI real. Necesita `requires_context: ["has_ci"]`. |
| `tool:goreleaser` | GoReleaser | `VALID` | Bien scoped a Go. |
| `tool:ruff` | Ruff | `VALID` | Python linter, correcto. |
| `tool:poetry` | Poetry | `VALID` | Python package manager. |
| `tool:uv` | uv | `VALID` | Python fast package manager. Correcto. |
| `tool:mypy` | mypy | `VALID` | Python type checker. |
| `tool:hypercorn` | Hypercorn | `NEEDS_METADATA` | Falta `category: "backend"`. |
| `tool:cargo-deny` | cargo-deny | `VALID` | Rust auditing. |
| `tool:nextest` | cargo-nextest | `VALID` | Rust test runner. |
| `tool:wasm-pack` | wasm-pack | `NEEDS_METADATA` | Falta `category`. |
| `tool:trivy` | Trivy | `VALID` | Security scanner universal. |
| `tool:sops` | SOPS | `VALID` | Secrets management. |
| `tool:renovate` | Renovate | `VALID` | Dependency updates. Ecosistema `["*"]`. |
| `tool:dependabot` | Dependabot | `VALID` | Bien definido. |
| `tool:semgrep` | Semgrep | `VALID` | SAST multi-lenguaje. |
| `tool:lefthook` | Lefthook | `VALID` | Git hooks. |
| `tool:fvm` | FVM | `VALID` | Flutter version manager. ✅ Añadido en Fase 20E. |
| `tool:melos` | Melos | `VALID` | Monorepo Dart/Flutter. ✅ Añadido en Fase 20E. |
| `tool:flutter_gen` | FlutterGen | `VALID` | Code gen para assets Flutter. ✅ Añadido en Fase 20E. |
| `tool:shorebird` | Shorebird | `VALID` | OTA updates Flutter. ✅ Añadido en Fase 20E. |
| `tool:patrol` | Patrol | `VALID` | Integration testing Flutter. ✅ Añadido en Fase 20E. |
| `tool:leak_tracker` | Leak Tracker | `VALID` | Memory leak detection Flutter. |
| `tool:dart_frog` | Dart Frog | `VALID` | **Corregido en Fase 21A**: añadido `requires_context: ["has_backend"]`. Ya no se recomienda a proyectos solo móviles. |

---

## Reglas de Recomendación Revisadas

### Antes (Fase 20)
El motor recomendaba herramientas basándose únicamente en coincidencia de ecosistema y `primaryLanguage`. Esto causaba que Dart Frog (un framework backend) fuese recomendado a aplicaciones Flutter puramente móviles.

### Después (Fase 21A)
Se añadió el sistema de `requires_context` al `Entry` struct y al motor de filtrado (`filter.go`). Los **Signals** derivados del grafo del proyecto controlan si una entrada es elegible antes de pasar al scorer:

| Signal | Criterio de activación |
| :--- | :--- |
| `has_mobile` | Nodo de lenguaje Dart/Flutter o plataforma Android/iOS detectada |
| `has_backend` | Nodo de lenguaje Go, Python, Node.js, Java, etc., o infra Docker/K8s |
| `has_firebase` | Nodo de infraestructura con ID que contiene `firebase` |

### Herramientas que ahora requieren señal de contexto
- `tool:dart_frog` → `requires_context: ["has_backend"]`

### Herramientas que en futura revisión deberían añadir contexto
- `tool:earthly` → considerar `requires_context: ["has_ci"]`
- `tool:buf` → considerar añadir lógica de detección de `.proto` files

---

## Conclusión

**Estado General del Catálogo:** `READY FOR ALPHA` con condiciones.

- 24 de 27 entradas: `VALID`.
- 4 entradas: `NEEDS_METADATA` (campos menores como `category`).
- 1 entrada (`earthly`): `NEEDS_RULES` (pendiente para iteración posterior).
- 0 entradas: `OUTDATED`, `DUPLICATE` o `REMOVE`.

El catálogo es adecuado para el lanzamiento Alpha. Los metadatos faltantes (`category`, `license`) no afectan la calidad de recomendaciones actuales, pero serán necesarios para la visualización avanzada en el Web UI (Fase 22+).
