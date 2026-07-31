# Arquitectura técnica

## 1. Vista general

```text
Authorized Root
    │
    ▼
Repository Discovery
    │
    ▼
Filesystem Inventory ── Ignore/Safety Policy
    │
    ├── Language Detection
    ├── Manifest Registry
    ├── Module/Monorepo Detection
    ├── Framework/Platform Detection
    ├── Tooling/Editor Detection
    └── Infrastructure Detection
             │
             ▼
        Evidence Store
             │
             ▼
        ProjectGraph
             │
    ┌────────┴─────────┐
    ▼                  ▼
Fingerprint        Local Report
    │
Privacy Sanitizer
    │
User Preview/Consent
    │
Remote Catalog API (future)
    │
Deterministic Recommendation Engine
    │
Top 3 + optional --more
```

## 2. Lenguajes

### CLI

- Go 1.26.x, versión exacta fijada por `mise.toml`.
- Un solo binario por sistema y arquitectura.
- La lógica principal no debe depender de Node, Python, Java o Docker instalados.

### Backend futuro

- TypeScript.
- Cloudflare Workers.
- D1 para datos estructurados.
- R2 únicamente cuando se requieran imágenes o snapshots.
- Hono y Zod son candidatos, no dependencias aprobadas hasta la fase de backend.

## 3. Estructura prevista del código

```text
cmd/
  stackgenome/
    main.go

internal/
  app/                 orquestación de casos de uso
  cli/                 comandos, flags, renderizado
  config/              configuración local
  discovery/           repositorios, módulos y límites
  filesystem/          inventario seguro y abstracción FS
  ignore/              gitignore y exclusiones
  language/            detección y clasificación
  manifest/            registro y parseo de manifests
  analyzer/            contratos y registro de adaptadores
  analyzers/           implementaciones por ecosistema
  framework/           reglas de frameworks
  platform/            móvil, web, escritorio, backend, etc.
  tooling/             editor, CI/CD, build tools
  infrastructure/      Docker, Kubernetes, Terraform, cloud
  evidence/            evidencia, procedencia y confianza
  projectgraph/        nodos, aristas y validación
  privacy/             redacción y allowlist
  fingerprint/         documento enviable y hashing
  catalog/             catálogo local/remoto
  recommendation/      filtros, scoring y diversidad
  report/              terminal, JSON, CycloneDX
  environment/         escaneo local opt-in
  telemetry/           desactivado por defecto; futuro

pkg/
  schema/              contratos públicos solo si son estables

testdata/
  fixtures/
```

No se debe crear `pkg/` por costumbre. Solo se mueve código allí cuando exista un consumidor externo real.

## 4. Dependencias y dirección

```text
CLI/UI -> Application -> Domain
                 └────> Ports
Adapters/Infrastructure -> Ports/Domain
```

El dominio no importa:

- UI;
- HTTP;
- Cloudflare;
- bases de datos;
- comandos de sistema;
- parsers concretos.

## 5. Contrato de analizador

Concepto inicial, sujeto a especificación durante las fases correspondientes:

```go
type Analyzer interface {
    ID() string
    Version() string
    Detect(ctx context.Context, input DetectInput) ([]Candidate, error)
    Analyze(ctx context.Context, input AnalyzeInput) (AnalysisResult, error)
}
```

Cada resultado debe incluir:

- `analyzer_id`;
- versión del analizador;
- evidencias;
- confidence;
- warnings;
- capabilities;
- errores parciales.

Un analizador no debe:

- ejecutar scripts del proyecto;
- instalar paquetes;
- resolver dependencias mediante código arbitrario;
- escribir dentro del proyecto inspeccionado;
- hacer red sin autorización;
- leer fuera de la raíz.

## 6. Registro de adaptadores

El MVP usa registro compilado:

```go
registry.Register(analyzer)
```

Razones:

- auditabilidad;
- portabilidad;
- binario autocontenido;
- ausencia de ABI inestable;
- menor superficie de ataque.

La extensibilidad comunitaria inicial será mediante contribuciones al repositorio. Plugins externos quedan para una fase posterior con protocolo fuera de proceso, firma y sandbox.

## 7. Detección de lenguajes

Estrategia combinada:

1. extensiones;
2. nombres especiales;
3. heurísticas tipo Linguist/go-enry;
4. shebang;
5. contenido limitado cuando sea necesario;
6. exclusión de binarios, dependencias vendorizadas y archivos generados;
7. ponderación por bytes y por archivos;
8. contexto de manifests.

Clasificaciones:

- `primary`;
- `satellite`;
- `generated`;
- `vendored`;
- `markup`;
- `data`;
- `configuration`;
- `unknown`.

## 8. Módulos y monorepos

Un módulo es una unidad con:

- raíz;
- tipo;
- manifests;
- lenguajes;
- framework;
- outputs;
- relaciones.

Detectores de workspace, entre otros:

- npm/pnpm/yarn/bun workspaces;
- Nx/Turborepo;
- Go workspaces;
- Cargo workspaces;
- Gradle multi-project;
- Maven modules;
- .NET solutions;
- Swift packages;
- Dart workspaces;
- CMake subprojects;
- Bazel.

La detección debe manejar repositorios anidados sin doble conteo.

## 9. Evidencia y confianza

Toda detección produce evidencia estructurada:

```json
{
  "kind": "manifest",
  "path": "apps/web/package.json",
  "selector": "dependencies.next",
  "value": "16.1.0",
  "sensitivity": "public-metadata"
}
```

La confianza no es una sensación del agente. Se deriva de reglas documentadas:

- 1.00: declaración explícita y parseada;
- 0.90: lockfile o configuración inequívoca;
- 0.75: combinación de archivos consistente;
- 0.50: heurística;
- menor de 0.50: no se presenta como hecho.

## 10. Fallos parciales

Un manifest corrupto no aborta todo el análisis. El resultado incluye warning y continúa. Solo son fatales:

- raíz inválida;
- falta de permisos sobre la raíz;
- violación de límites;
- configuración inválida esencial;
- cancelación.

## 11. Rendimiento

- recorrido único cuando sea posible;
- límites de profundidad y tamaño configurables;
- workers acotados;
- context cancellation;
- no seguir symlinks fuera de la raíz;
- caché local opcional por hash de metadatos;
- salida streaming para progreso, no para resultados incompletos.

## 12. Recomendaciones

Pipeline:

1. recuperar candidatos;
2. filtros duros;
3. score;
4. penalizaciones;
5. diversidad;
6. top 3;
7. razones estructuradas.

Filtros duros:

- recurso disponible;
- ecosistema;
- plataforma;
- rango de versiones;
- licencia/política;
- vulnerabilidad crítica según política;
- incompatibilidad conocida.

Score inicial conceptual:

- compatibilidad: 30%;
- seguridad: 25%;
- mantenimiento: 20%;
- adecuación al objetivo: 10%;
- documentación: 5%;
- popularidad externa: 5%;
- valoración interna: 5%.

Los pesos serán configuración versionada y probada.

## 13. Integraciones externas

Se pueden reutilizar conceptos o datos de:

- go-enry/Linguist;
- PURL;
- CycloneDX;
- Syft;
- OSV;
- OpenSSF Scorecard;
- registries oficiales.

No se incorporará un proyecto completo como dependencia sin revisar:

- licencia;
- tamaño;
- modelo de actualización;
- API estable;
- impacto de binario;
- capacidad offline.

## 14. Observabilidad

En el CLI:

- `--verbose`;
- `--debug`;
- logs sin datos sensibles;
- trace local opcional;
- ningún analytics por defecto.

En backend futuro:

- request id;
- métricas agregadas;
- no almacenar fingerprints completos por defecto;
- retención documentada.

## 15. Compatibilidad

Objetivos de release:

- macOS arm64 y amd64;
- Linux amd64 y arm64;
- Windows amd64 y arm64 cuando el toolchain lo permita;
- formatos de salida versionados.

La ruta `/Volumes/intento1` existe solo en desarrollo local y jamás debe aparecer en lógica de producto ni fixtures públicos.
