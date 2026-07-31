# ADR-0002: Local-first y metadata-only

- Estado: aceptado
- Fecha: 2026-07-24

## Decisión

El análisis ocurre localmente. El payload remoto es una proyección allowlist visible para el usuario. No se envía código fuente.

## Consecuencias

- mayor confianza;
- backend económico;
- más responsabilidad en el CLI;
- necesidad de sanitizer y pruebas adversariales;
- algunas recomendaciones tendrán menor contexto, lo cual es aceptable.
