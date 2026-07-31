# Handoff

## Última sesión

Se completó la **Fase 19 — Validación Piloto con Territor**.
Resultado: `PILOT_PASSED_WITH_FINDINGS`.

## Completado

- **Fase 0-15**: Creación del analizador, tests, backend, sanitización, hardenización.
- **Fase 16**: Auditoría independiente Alpha.
- **Fase 17**: Automatización CI/CD.
- **Fase 18A**: Web Readiness Gate.
- **Fase 18B**: Web Alpha (Next.js + d3-force + FileImporter local).
- **Fase 19**: Validación con Territor (Flutter/Dart real).
  - 3 ejecuciones deterministas sobre repositorio real.
  - 7 hallazgos identificados (0 BLOCKER, 2 HIGH, 3 MEDIUM, 2 LOW).
  - Privacidad confirmada CLEAN.
  - Territor no modificado.
  - Reporte: `docs/validation/TERRITOR_PILOT_VALIDATION_2026-07.md`.

## Hallazgos pendientes para Fase 20

| ID | Severidad | Descripción |
|---|---|---|
| F-001 | HIGH | Excluir ios/Pods/, android/.gradle/, build/ del análisis de lenguajes |
| F-002 | MEDIUM | Deduplicar nodos de plataforma por nombre+tipo |
| F-003 | HIGH | Añadir recursos Flutter/Dart al catálogo (fvm, melos, flutter_gen, etc.) |
| F-004 | MEDIUM | Extraer versiones de dependencias en pubspec.yaml |
| F-005 | MEDIUM | Diferenciar dev_dependencies de prod en Dart |
| F-006 | LOW | Detectar shorebird.yaml como tooling |
| F-007 | LOW | Detectar firebase.json como infraestructura |

## Siguiente acción

Pendiente autorización explícita para **Fase 20 — Correcciones post-piloto**.
