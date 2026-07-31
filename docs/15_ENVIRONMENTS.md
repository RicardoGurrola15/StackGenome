# Matriz de Entornos — StackGenome

## Resumen

StackGenome opera con cuatro entornos bien diferenciados.
Ningún entorno inferior puede afectar a uno superior sin autorización explícita y proceso de escalada.

---

## Entorno: Local (Development)

| Parámetro | Valor |
|---|---|
| **URL API** | `http://localhost:8787` (via `wrangler dev`) |
| **Worker** | Local en memoria (no desplegado) |
| **D1** | SQLite local gestionado por Wrangler |
| **CLI Version** | `dev` (sin `-ldflags`) |
| **Variables** | `.dev.vars` (ignorado por `.gitignore`) |
| **Secretos** | `.dev.vars` (nunca en Git) |
| **CORS** | `localhost:3000`, `localhost:5173` |
| **Despliegue** | Manual vía `wrangler dev` |
| **Datos permitidos** | Datos ficticios o proyectos locales propios |

---

## Entorno: Staging

| Parámetro | Valor |
|---|---|
| **URL API** | `https://stackgenome-api-staging.stackgenome.workers.dev` |
| **Worker** | `stackgenome-api-staging` |
| **D1** | `stackgenome-catalog-staging` (binding: `DB`) |
| **CLI Version** | Versión del branch `main` |
| **Variables** | `wrangler.toml [env.staging]` |
| **Secretos** | GitHub Secret: `CLOUDFLARE_API_TOKEN` |
| **CORS** | `localhost:3000`, `localhost:5173`, `*.stackgenome.pages.dev`, `staging.stackgenome.com` |
| **Despliegue** | Automático desde `main` vía `deploy-backend.yml` (solo en cambios a `/backend/**`) |
| **Retención** | Indefinida; datos de catálogo de herramientas semilla |
| **Datos permitidos** | Catálogo de herramientas públicas; Fingerprints efímeros (no almacenados) |

> [!IMPORTANT]
> El staging es el único entorno cloud actualmente autorizado.

---

## Entorno: Preview (Futuro / Opcional)

| Parámetro | Valor |
|---|---|
| **URL API** | `https://<branch>-stackgenome-api.stackgenome.workers.dev` (Cloudflare Pages Preview) |
| **Worker** | Copia efímera del worker de staging |
| **D1** | Compartida con staging (lectura) o copia efímera |
| **Secretos** | Mismos que staging; nunca credenciales de producción |
| **Despliegue** | Automático por PR; eliminado al cerrar PR |

---

## Entorno: Production (BLOQUEADO — No creado)

| Parámetro | Valor |
|---|---|
| **URL API** | `https://api.stackgenome.com` *(no creado)* |
| **Worker** | *No creado* |
| **D1** | *No creado* |
| **Secretos** | *No configurados* |
| **Despliegue** | **Requiere autorización explícita antes de crear** |

> [!CAUTION]
> Producción **no existe** y no debe ser creado sin:
> - revisión de seguridad completa,
> - dominio personalizado configurado,
> - plan de rollback documentado,
> - autorización explícita del dueño del proyecto.

---

## Secretos Requeridos por Entorno

| Secreto | Entorno | Almacenamiento |
|---|---|---|
| `CLOUDFLARE_API_TOKEN` | Staging | GitHub Secret (Actions) |
| `.dev.vars` | Local | Archivo local, nunca en Git |

---

## Estrategia de Sincronización de Datos

- Las migraciones D1 son **unidireccionales** y **numeradas** (`0001_`, `0002_`).
- Se aplican manualmente con `wrangler d1 migrations apply DB --remote` previa autorización.
- **Las migraciones de staging nunca se aplican a producción directamente.**
