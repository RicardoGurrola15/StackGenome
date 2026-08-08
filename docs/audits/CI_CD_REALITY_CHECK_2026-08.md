# CI/CD Reality Check (Auditoría Agosto 2026)

Esta auditoría resuelve la discrepancia sobre el estado del CI/CD reportado en Fases anteriores. Se verificaron los workflows presentes en `.github/workflows/` y se contrastó su definición con el historial de ejecución.

## Resumen de Workflows Existentes

### 1. `ci.yml`
* **Propósito**: Ejecutar linters y tests en cada Push y Pull Request a la rama `main`.
* **Triggers**: `push`, `pull_request` (sobre `main`).
* **Jobs**:
  * `test-go`: 
    * Instala Go 1.22.x.
    * Ejecuta: `go mod verify`, `go fmt`, `go vet`, `go test -v -race ./...`.
    * Instala y ejecuta: `govulncheck`.
    * **Estado**: `CONFIGURED` pero frecuentemente **`BROKEN`** debido a fallos recientes.
  * `test-backend`:
    * Instala Node 20.
    * Ejecuta: `npm ci`, `npx tsc --noEmit`, `npm test` (vitest).
    * **Estado**: `CONFIGURED` pero actualmente **`BROKEN`** (Error de permisos intentando crear `/Volumes` durante la instalación de dependencias en GitHub Actions).

### 2. `deploy-backend.yml`
* **Propósito**: Desplegar el backend Cloudflare Worker (D1 + API) al entorno staging.
* **Triggers**: `push` a `main` cuando se modifican archivos en `backend/**`.
* **Jobs**:
  * `deploy`:
    * Instala Node 20, ejecuta `npm ci`.
    * Ejecuta `cloudflare/wrangler-action@v3` para hacer deploy con environment `staging`.
    * Secretos usados: `CLOUDFLARE_API_TOKEN`.
    * **Estado**: `CONFIGURED`, pero dependiente de la correcta inyección del token y arreglos en npm. **`BROKEN`** temporalmente por el mismo fallo de npm de `/Volumes`.

### 3. `release.yml`
* **Propósito**: Automatizar la compilación cruzada y publicación de GitHub Releases con los binarios.
* **Triggers**: Push a tags que coinciden con `v*` (ej. `v0.1.0`).
* **Jobs**:
  * `release`:
    * Ejecuta el script local `scripts/build-release.sh`.
    * Construye binarios para: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`.
    * Genera sumas de comprobación (`checksums.txt` vía `sha256sum`/`shasum`).
    * Crea un GitHub Release y sube los archivos de `dist/`.
    * **Estado**: `CONFIGURED` y completamente funcional a nivel de script. Sin embargo, su estado histórico es **`NOT_EXECUTED_SUCCESSFULLY`** (nunca se ha lanzado un tag `v*` real para dispararlo).

## Capacidades: Expectativa vs Realidad

| Funcionalidad | Expectativa (Fase 17) | Estado Real | Detalles |
| :--- | :--- | :--- | :--- |
| **Go Tests (`go test`, `race`, `vet`)** | Terminado | `CONFIGURED` | Existe en `ci.yml`. |
| **`govulncheck`** | Terminado | `CONFIGURED` | Existe en `ci.yml`. |
| **Backend Tests & `tsc`** | Terminado | `BROKEN` | Fallos de permisos (`/Volumes`) en NPM en Actions. |
| **Web Tests & Build** | Ausente | `NOT_IMPLEMENTED` | No existe job de CI para verificar la UI web. |
| **Cross-Platform Compilation** | Terminado | `CONFIGURED` | Script `build-release.sh` correcto. Funciona en local. |
| **GitHub Releases Automatizados** | Terminado | `CONFIGURED` | Existe `release.yml`, pero nunca se ha probado en CI porque no se han creado tags. |
| **Despliegue Staging (Backend)** | Terminado | `CONFIGURED` | Configurado vía wrangler-action. |
| **Despliegue Web (Frontend)** | Ausente | `NOT_IMPLEMENTED` | No hay workflow para Cloudflare Pages / Assets. |

## Conclusión

El reporte de la Fase 17 fue *técnicamente veraz* al decir que se "crearon" (escribieron) los archivos de GitHub Actions. Sin embargo, el reporte de la Fase 20 fue *prácticamente veraz* al insinuar que todavía hace falta configurar/validar el CI/CD, dado que los workflows actuales están fallando (npm) y faltan piezas clave (web build, web tests, release efectivo con tags). 

**Regla aplicada:** Ya existen los archivos fundamentales para backend, Go, y releases. **NO SE RECONSTRUIRÁN**. Durante la Fase 21 (Release Público), simplemente se arreglarán los errores (ej. permisos de NPM), se probará un tag real y se añadirá el pipeline para la Web.
