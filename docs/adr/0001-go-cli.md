# ADR-0001: Go para el CLI

- Estado: aceptado
- Fecha: 2026-07-24

## Contexto

El CLI debe analizar repositorios heterogéneos, distribuirse en macOS/Linux/Windows y no exigir runtimes adicionales.

## Decisión

Usar Go para el CLI y motor local.

## Consecuencias

Positivas:
- binarios autocontenidos;
- concurrencia y rendimiento;
- biblioteca estándar;
- cross-compilation;
- ecosistema de tooling.

Negativas:
- backend/web usarán TypeScript;
- algunos parsers pueden existir antes en otros lenguajes;
- no se comparte código directo con la web.

## Alternativas

- TypeScript/Node: iteración rápida, distribución más compleja.
- Rust: alto control, mayor complejidad inicial.
- Python: parsers abundantes, runtime y empaquetado menos adecuados.
