# Biblioteca de prompts por fase

## Uso

1. Usa primero `START_PROMPT.md`.
2. Cuando autorices una fase, copia el prompt correspondiente.
3. Reemplaza únicamente los campos entre `<...>`.
4. El agente debe detenerse al finalizar.

## Prompt base para cualquier fase

```text
Trabaja en StackGenome y ejecuta exclusivamente la Fase <NOMBRE/NUMERO> autorizada.

Antes de modificar:
- lee AGENTS.md;
- lee docs/00_INDEX.md;
- lee docs/07_ROADMAP_BLUEPRINT.md;
- lee .project/CURRENT_PHASE.md, STATE.md y HANDOFF.md;
- revisa los ADR aplicables;
- recupera memoria de Engram si está disponible y contrástala con el repositorio;
- confirma que /Volumes/intento1 está montado para el entorno local.

Presenta un plan breve limitado a esta fase. Después impleméntalo sin iniciar fases posteriores.

Al terminar:
- ejecuta pruebas;
- revisa seguridad y privacidad;
- actualiza documentación y .project;
- registra decisiones;
- guarda una memoria breve si Engram está disponible;
- presenta evidencia;
- pregunta “¿Autorizas avanzar a la fase siguiente?”;
- detente.
```

## Fase 0

```text
Ejecuta exclusivamente la Fase 0: validación documental y entorno.

No escribas código de producto. Diagnostica:
- ruta y estado Git;
- montaje de /Volumes/intento1;
- scripts de entorno;
- instalación actual de mise y Go;
- rutas reales de programas y cachés;
- disponibilidad de Codex, Antigravity, GentleAI y Engram;
- integración de Engram por proyecto.

No cambies configuraciones globales. Si falta una herramienta, muestra el comando, destino y archivos afectados y pide autorización antes de instalar.

Valida todos los documentos y señala contradicciones. Actualiza solo archivos de estado/documentación autorizados. Finaliza con evidencia y detente.
```

## Fase 1

```text
Ejecuta exclusivamente la Fase 1: fundación del repositorio Go.

Crea el módulo y un CLI mínimo que solo soporte ayuda y versión. Define estructura limpia, pruebas mínimas y comandos reproducibles. No implementes descubrimiento, analizadores ni backend.

Revisa cada dependencia; prefiere estándar. Ejecuta format, vet, test y build. Actualiza estado y detente.
```

## Fase 2

```text
Ejecuta exclusivamente la Fase 2: ProjectGraph y Evidence.

Implementa modelos, invariantes, ids, rutas relativas, serialización estable y tests golden. No implementes recorrido real ni parsers de ecosistemas.

Asegura que el modelo no acepte rutas absolutas o campos sensibles. Documenta cualquier ajuste del schema. Detente.
```

## Fase 3

```text
Ejecuta exclusivamente la Fase 3: descubrimiento seguro.

Implementa abstracción de filesystem, recorrido, exclusiones, límites y seguridad de symlinks usando fixtures sintéticas. No detectes aún frameworks ni dependencias.

Prueba traversal, ciclos, permisos, archivos grandes y cancelación. Detente.
```

## Fase 4

```text
Ejecuta exclusivamente la Fase 4: detección universal de lenguajes.

Evalúa formalmente go-enry frente a una implementación mínima. Documenta licencia, tamaño y exactitud antes de adoptar. Clasifica primary, satellite, generated, vendored, markup y config.

Añade fixtures diversos y resultados deterministas. No inicies manifests. Detente.
```

## Fase 5

```text
Ejecuta exclusivamente la Fase 5: manifest registry.

Diseña contratos y detección de manifests/lockfiles para un conjunto amplio, inicialmente detection-only. Un manifest corrupto debe producir warning, no abortar el análisis.

Demuestra que un nuevo adaptador puede añadirse sin editar el núcleo. Detente.
```

## Fase 6

```text
Ejecuta exclusivamente la Fase 6: adaptadores profundos ola A.

Implementa JS/TS, Python, Go y Rust en incrementos pequeños. Extrae dependencias directas, versiones, workspaces y frameworks aprobados sin ejecutar scripts.

Usa PURL cuando corresponda y documenta cobertura por manager. Detente.
```

## Fase 7

```text
Ejecuta exclusivamente la Fase 7: adaptadores profundos ola B.

Implementa los ecosistemas autorizados de JVM, Apple, .NET, Dart, PHP, Ruby y C/C++ sin fingir paridad. Declara soporte full/substantial/basic.

Divide el trabajo en subpasos dentro de la fase, pero no avances a la fase 8. Detente.
```

## Fase 8

```text
Ejecuta exclusivamente la Fase 8: frameworks, plataformas, infraestructura y editores declarados.

Usa evidencia explícita y reglas versionadas. No leas plugins instalados en la máquina; eso pertenece a la fase 9.

Añade confidence y evita afirmar ausencias como defectos. Detente.
```

## Fase 9

```text
Ejecuta exclusivamente la Fase 9: entorno opcional.

Diseña y desarrolla un scanner off-by-default. Cada comando debe estar allowlisted, sin shell interpolation, con timeout y sanitización.

Implementa primero macOS sin romper la interfaz portable. No envíes datos. Detente.
```

## Fase 10

```text
Ejecuta exclusivamente la Fase 10: Fingerprint y privacidad.

Implementa proyección allowlist, secret checks, preview y schema. Construye pruebas adversariales. El nombre del proyecto y rutas absolutas se excluyen.

No conectes backend. Detente.
```

## Fase 11

```text
Ejecuta exclusivamente la Fase 11: catálogo local y recomendaciones.

Crea un catálogo semilla pequeño y revisable. Implementa filtros, ranking determinista, diversidad y razones estructuradas. No uses IA.

Evalúa manualmente repositorios de tipos distintos. Entrega top 3 y --more. Detente.
```

## Fase 12

```text
Ejecuta exclusivamente la Fase 12: backend Cloudflare.

Antes de código, verifica límites y precios actuales en documentación oficial. Propón schema D1, API y privacidad. Implementa el mínimo aprobado.

No agregues login, comunidad, imágenes ni web. No persistas fingerprints por defecto. Detente.
```

## Fase 13

```text
Ejecuta exclusivamente la Fase 13: UX y conexión remota.

Estabiliza comandos, códigos de salida, consentimiento, offline, JSON y completions. No cambies ranking salvo bug probado.

Prueba terminales y errores. Detente.
```

## Fase 14

```text
Ejecuta exclusivamente la Fase 14: hardening.

Ejecuta fuzzing, race, benchmarks, govulncheck, SBOM y revisión de amenazas. Corrige únicamente hallazgos dentro del alcance o documenta bloqueos.

No publiques todavía. Detente.
```

## Fase 15

```text
Ejecuta exclusivamente la Fase 15: alpha.

Prepara releases reproducibles, checksums, documentación, SECURITY y CONTRIBUTING. Ejecuta matriz de CI y valida licencias.

No inicies la web. Entrega recomendación go/no-go y detente.
```
