# Arquitectura de Integración Web — StackGenome

## Visión General

El frontend web de StackGenome opera como una **capa de visualización** sobre el grafo generado por el CLI.
No reemplaza al CLI: lo complementa.

```
┌─────────────────────┐     export --json     ┌───────────────────────┐
│  StackGenome CLI    │ ──────────────────────► │  project_graph.json   │
│  (Go, local)        │                         │  (Archivo local)      │
└─────────────────────┘                         └───────────┬───────────┘
                                                            │ file open
                                                            ▼
                                               ┌────────────────────────┐
                                               │    Web UI (Browser)     │
                                               │  • Visualización grafo  │
                                               │  • Resúmenes locales    │
                                               │  • Recomendaciones      │
                                               └───────────┬────────────┘
                                                           │ opt-in
                                                           │ fingerprint
                                                           ▼
                                               ┌────────────────────────┐
                                               │  Cloudflare Worker     │
                                               │  (Staging API)         │
                                               │  • Ranking remoto      │
                                               │  • Catálogo live       │
                                               └────────────────────────┘
```

---

## Flujo de Datos Detallado

### Flujo A: Solo Local (Sin red)

1. El usuario ejecuta `stackgenome analyze --json > graph.json`.
2. El usuario abre el archivo en el navegador desde la UI ("Abrir archivo").
3. El browser valida el schema contra `openapi.yaml#/components/schemas/ProjectGraphDTO`.
4. El grafo se renderiza localmente con los nodos, aristas, y recomendaciones offline ya incluidas.

### Flujo B: Consulta Remota (Opt-in)

1. El usuario elige "Obtener recomendaciones actualizadas".
2. La UI extrae el Fingerprint sanitizado del grafo cargado (solo `nodes` + `edges`, sin evidencias ni paths).
3. Se muestra un preview del payload antes de enviarlo.
4. El usuario confirma.
5. Se envía a `POST /api/v1/recommendations` en staging.
6. El Worker valida, rankea y retorna las recomendaciones.
7. La UI muestra los resultados.

---

## Arquitectura del Frontend Propuesta (Fase 18B)

| Capa | Tecnología Propuesta | Justificación |
|---|---|---|
| **Framework** | Next.js (App Router) | SSG para la landing, CSR para la UI de grafo |
| **Visualización** | React Force Graph | Visualización de grafos tipo D3 sin deps pesadas |
| **Validación de schema** | `ajv` (JSON Schema) | Validar el archivo abierto contra el openapi |
| **Estado** | `zustand` | Estado liviano para el grafo activo |
| **Estilos** | Vanilla CSS + CSS Variables | Alineado con las directrices del proyecto |
| **Deployment** | Cloudflare Pages | Mismo ecosistema que el worker backend |

---

## Criterios para Iniciar la Web Alpha (Gate de Fase 18A)

- [x] CI/CD verde en `main` (el golden test fue corregido).
- [x] Backend Staging responde (verificado en Fase 16).
- [x] Contratos versionados en `docs/api/openapi.yaml`.
- [x] CORS restringido a dominios autorizados.
- [x] Política de privacidad definida (`docs/16_WEB_PRIVACY.md`).
- [x] Matriz de entornos definida (`docs/15_ENVIRONMENTS.md`).
- [x] Sin secretos versionados en Git.
- [x] Catálogo staging valida correctamente (30 herramientas, 0 errores).
- [x] Producción bloqueada (Worker/D1 no creados).
