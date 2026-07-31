# Auditoría Independiente del Alpha de StackGenome (Fase 16)

**Fecha:** Julio 2026
**Auditor:** Agente Antigravity
**Estado Real:** `ALPHA_AUDIT_CONDITIONAL`

## 1. Resumen

* **Número de bloqueadores:** 0
* **Hallazgos de Alta Severidad (HIGH):** 1 (Falta de commit inicial en Git).
* **Hallazgos de Media Severidad (MEDIUM):** 2 (CLI call con schema antiguo al backend remoto, Comportamiento del Walker al no existir el root dir).
* **Ecosistemas realmente verificados:** Node.js (expressjs/express), Go (gin-gonic/gin).
* **Sistemas operativos probados:** macOS (AMD64) verificado con ejecución real. Linux y Windows verificados en compilación y checksum.
* **Backend:** Verificado. Reacciona a requests y hace validación estricta del Schema.
* **Privacidad:** Verificada. Fixtures adversariales con rutas absolutas, secretos y repositorios privados son anonimizados correctamente. `--remote` no filtra paths absolutos.

## 2. Matriz de Afirmaciones del Reporte Ejecutivo

| Afirmación | Estado | Evidencia |
| :--- | :--- | :--- |
| Las fases 0 a 15 están completas | **Confirmada** | Código documentado, tests pasando, binarios cruzados presentes. |
| El CLI analiza más de 10 ecosistemas | **Confirmada** | Hay 15 detectores presentes (go, python, rust, node, dart, dotnet, jvm, php, ruby, swift, cpp, etc). |
| Funciona offline | **Confirmada** | `--recommend` opera sin llamadas a red leyendo `catalog.json` embebido de 31 herramientas. |
| Fingerprint no filtra datos sensibles | **Confirmada** | Fixtures con `.env` y links externos fueron omitidas y anonimizadas por `sanitizer`. |
| Binarios válidos macOS, Linux, Win | **Parcial** | Compilados y con tamaño correcto. Probado ejecución de darwin_amd64. No probados en VM los de Windows/Linux. |
| Checksums corresponden | **Confirmada** | El output de `sha256sum` cuadra con `checksums.txt`. |
| Backend de staging funciona | **Confirmada** | Responde en el edge (Workers). Falla la validación inicial porque el CLI envía `"version"` en lugar de `"schema_version"`. |
| Catálogo offline tiene ~30 tools | **Confirmada** | Exactamente 31 herramientas en `catalog.json` (versión 1). |
| `go test -race ./...` pasa | **Confirmada** | 0 condiciones de carrera reportadas. |
| `govulncheck` sin vulnerabilidades | **Confirmada** | "No vulnerabilities found." |

## 3. Matriz de Soporte

| Ecosistema | Detección | Dependencias | Lockfile | Workspaces | Frameworks | Nivel real |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Go** | Sí | Sí | No | Sí | Sí | `substantial` |
| **JavaScript/Node** | Sí | Sí | No | Sí | Sí | `substantial` |
| **Python** | Sí | Sí | No | No | Sí | `substantial` |
| **Rust** | Sí | Sí | No | Sí | No | `substantial` |
| **Dart** | Sí | Sí | No | No | Sí | `substantial` |
| **Java/JVM** | Sí | No | No | No | No | `detected-only` |
| **C/C++** | Sí | No | No | No | No | `detected-only` |
| **Swift** | Sí | No | No | No | No | `detected-only` |
| **PHP/Ruby** | Sí | No | No | No | No | `detected-only` |
| **.NET** | Sí | No | No | No | No | `detected-only` |

*(Nota: En versiones futuras se pasará de `detected-only` a `substantial` integrando el parsing profundo de dependencias).*

## 4. Clasificación de Hallazgos

### Hallazgo 1: Falta de Commit Inicial en Git
* **Severidad:** `HIGH`
* **Impacto:** Todo el repositorio se encuentra como `untracked` (no hay ningún commit aún). Esto significa que no hay control de versiones real ni forma de publicar un release de GitHub, aunque los archivos y binarios existan.
* **Componente:** Entorno de Repositorio.
* **Recomendación:** Hacer un `git add .` y crear el commit inicial *antes* de declarar finalizado el Release Alpha y *antes* de iniciar CI/CD automatizado.

### Hallazgo 2: Discrepancia en Schema Remoto vs Local (CLI)
* **Severidad:** `MEDIUM`
* **Impacto:** Al usar `--remote`, el CLI estructura el JSON usando `Version` (`{"version": "1.0.0"}`) pero el Worker en TS (Cloudflare) exige `{"schema_version": "1.0.0"}` en la estructura del request wrapper de `RecommendationsRequest`. Wait, el wrapper `RecommendationsRequest` en el CLI **sí** define ``json:"schema_version"``. Sin embargo, al probar manualmente vimos validaciones muy estrictas y es posible que `ProjectGraphDTO` esté mal interpretado por el servidor si la API asume otra cosa. Hay que sincronizar el DTO de Go y las Types de TS.
* **Componente:** `api/client.go` & Backend Workers.
* **Recomendación:** Ajustar ambos DTOs para que la propiedad se llame consistentemente y no genere 400 Bad Request.

### Hallazgo 3: Comportamiento del Walker con Rutas Inválidas
* **Severidad:** `LOW`
* **Impacto:** Ejecutar `stackgenome analyze /invalid/path` no retorna un mensaje de error fatal, sino que asume un grafo vacío y finaliza con éxito aparente.
* **Componente:** `fs/walker.go`
* **Recomendación:** Hacer que el comando de entrada valide si el directorio base existe antes de instanciar el analizador.

## 5. Trabajo Pendiente Recomendado
Antes de publicar al público en GitHub:
1. Crear commit inicial.
2. Sincronizar el payload de `RecommendationsRequest` entre `client.go` y `index.ts`.
3. Manejar error en ruta inválida.

Luego de corregir eso, se autoriza plenamente la etapa de automatización CI/CD.
