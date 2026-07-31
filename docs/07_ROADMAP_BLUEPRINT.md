# Roadmap / Blueprint

## Regla

Cada fase requiere autorización independiente. El agente completa una fase, verifica, documenta y se detiene.

## Fase 0 — Validación documental y entorno

Objetivo:
- confirmar documentos, rutas, Git, herramientas y memoria.

Entregables:
- repositorio inicializado si no lo está;
- `mise` y Go instalados en `/Volumes/intento1` solo con autorización;
- Engram/GentleAI diagnosticados;
- estado actualizado;
- ningún código de producto.

Criterios:
- volumen validado;
- rutas de caché verificadas;
- `go version`;
- documentación legible;
- propuesta de cambios de setup presentada.

## Fase 1 — Fundación del repositorio Go

Objetivo:
- crear módulo, estructura mínima, CLI que muestre versión y CI básica.

Entregables:
- `go.mod`;
- `cmd/stackgenome`;
- estructura `internal`;
- tests mínimos;
- lint/vet;
- build multiplataforma inicial.

Criterios:
- compila;
- tests pasan;
- sin dependencia innecesaria;
- versión inyectable;
- no lógica de análisis.

## Fase 2 — Dominio ProjectGraph y Evidence

Objetivo:
- implementar modelos, invariantes, validación y serialización estable.

Entregables:
- modelos;
- schema;
- golden tests;
- versionado.

Criterios:
- rutas relativas;
- ids únicos;
- orden determinista;
- JSON sin campos sensibles.

## Fase 3 — Descubrimiento seguro

Objetivo:
- recorrer raíz, detectar repositorios/módulos candidatos y aplicar exclusiones.

Entregables:
- FS abstraction;
- ignore engine;
- symlink safety;
- límites;
- fixtures.

Criterios:
- no sale de raíz;
- manejo de ciclos;
- cancelación;
- rendimiento básico.

## Fase 4 — Detección universal de lenguajes

Objetivo:
- identificar lenguajes principales, satélite, generados y vendorizados.

Entregables:
- detector;
- integración evaluada con go-enry o alternativa;
- clasificación;
- fixtures multi-lenguaje.

Criterios:
- resultados estables;
- exclusiones correctas;
- cobertura documentada.

## Fase 5 — Registry de manifests y package managers

Objetivo:
- arquitectura común para manifests y lockfiles.

Entregables:
- contratos;
- registry;
- detection-only para un conjunto amplio;
- errores parciales.

Criterios:
- nuevo ecosistema añadible sin modificar núcleo;
- manifest corrupto no aborta;
- evidencias completas.

## Fase 6 — Adaptadores profundos, ola A

Ecosistemas:
- JavaScript/TypeScript;
- Python;
- Go;
- Rust.

Objetivo:
- dependencias directas, versiones, workspaces y frameworks principales.

Criterios:
- fixtures por manager;
- no ejecutar scripts;
- PURL cuando aplique;
- cobertura declarada.

## Fase 7 — Adaptadores profundos, ola B

Ecosistemas:
- JVM/Gradle/Maven;
- Swift/Objective-C/SPM/CocoaPods;
- .NET;
- Dart/Flutter;
- PHP;
- Ruby;
- C/C++ con CMake/Conan/vcpkg en alcance acordado.

Objetivo:
- ampliar profundidad sin degradar núcleo.

Criterios:
- soporte honesto por nivel;
- fixtures;
- manejo de monorepo.

## Fase 8 — Frameworks, plataformas, infraestructura y editores

Objetivo:
- detectar web, backend, móvil, escritorio, contenedores, IaC, CI/CD y declaraciones de IDE.

Incluye:
- React/Next/Vue/Angular;
- Spring/ASP.NET/Django/FastAPI;
- Android/iOS/Flutter/Electron/Tauri;
- Docker/Kubernetes/Terraform;
- GitHub Actions;
- `.vscode`, `.idea`, devcontainers y configuraciones declarativas compatibles.

Criterios:
- no inferir herramientas personales desde archivos ambiguos;
- evidencia y confidence.

## Fase 9 — Escaneo opcional del entorno

Objetivo:
- versiones instaladas, OS/arquitectura y editor/plugins autorizados.

Criterios:
- off por defecto;
- preview;
- allowlist;
- timeout;
- no rutas personales;
- compatible con macOS primero y diseño portable.

## Fase 10 — Fingerprint y privacidad

Objetivo:
- proyección metadata-only, sanitizer y consentimiento.

Criterios:
- secret fixtures;
- campos allowlisted;
- preview legible;
- schema versionado;
- no nombre de proyecto por defecto.

## Fase 11 — Catálogo local y recomendaciones

Objetivo:
- validar valor sin backend.

Entregables:
- catálogo semilla versionado;
- filtros;
- scoring;
- top 3;
- objetivos.

Criterios:
- determinismo;
- explicación estructurada;
- recursos emergentes no excluidos por popularidad;
- evaluación manual diversa.

## Fase 12 — Backend Cloudflare

Objetivo:
- API mínima con D1 y catálogo remoto.

Criterios:
- límites free verificados;
- no persistir fingerprint;
- validación;
- rate limiting;
- ranking versionado;
- pruebas locales y staging.

## Fase 13 — UX final del CLI y conexión remota

Objetivo:
- comandos estables, errores, completions, JSON y consentimiento.

Criterios:
- terminal accesible;
- modo offline;
- remote opt-in;
- documentación de comandos.

## Fase 14 — Calidad, seguridad y rendimiento

Objetivo:
- hardening.

Incluye:
- fuzzing de parsers;
- race;
- benchmarks;
- `govulncheck`;
- SBOM;
- threat review;
- pruebas en tres OS.

## Fase 15 — Alpha

Objetivo:
- probar en repositorios diversos y preparar publicación.

Entregables:
- binarios;
- checksums;
- guía de contribución;
- política de seguridad;
- casos de prueba;
- feedback.

Criterios:
- no defectos bloqueantes;
- licencias revisadas;
- limitaciones publicadas;
- decisión explícita antes de web.

## Gates globales

No avanzar si:

- hay filtración de datos;
- no hay tests;
- el modelo requiere cambios incompatibles no documentados;
- el scope crece sin aprobación;
- el agente no puede demostrar qué hizo.
