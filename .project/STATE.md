# Estado del Proyecto: StackGenome

- **Fase actual**: 19 (Validación con Territor — COMPLETADA)
- **Resultado**: `PILOT_PASSED_WITH_FINDINGS`
- **Bloqueos**: Ninguno.
- **Siguiente fase pendiente de autorización**: Fase 20 (Correcciones post-piloto).

## Fases completadas
- [x] Fase 0-15: Core, analizadores, y pipeline de CI.
- [x] Fase 16: Auditoría Alpha.
- [x] Fase 17: CI/CD.
- [x] Fase 18A: Web Readiness Gate.
- [x] Fase 18B: Web Alpha.
- [x] Fase 19: Validación piloto con Territor (Flutter/Dart real).

## Hallazgos de Fase 19 (para resolver en Fase 20)
- F-001 HIGH: Excluir ios/Pods/ y vendored code del análisis de lenguajes.
- F-002 MEDIUM: Deduplicar nodos de plataforma por nombre.
- F-003 HIGH: Añadir recursos Flutter/Dart al catálogo.
- F-004 MEDIUM: Extraer versiones de dependencias desde pubspec.yaml.
- F-005 MEDIUM: Diferenciar dev vs prod dependencies en Dart.
- F-006 LOW: Detectar shorebird.yaml.
- F-007 LOW: Detectar firebase.json como infraestructura.

## Notas técnicas
- Análisis de Territor: 3 runs deterministas, SHA-256 idéntico.
- Privacidad: CLEAN (no rutas absolutas, no project IDs, no API keys en el JSON).
- Territor: No modificado. git status idéntico pre/post análisis.
- Reporte completo: docs/validation/TERRITOR_PILOT_VALIDATION_2026-07.md
