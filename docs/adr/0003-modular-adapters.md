# ADR-0003: Adaptadores modulares compilados

- Estado: aceptado
- Fecha: 2026-07-24

## Decisión

El núcleo define contratos; los ecosistemas se implementan mediante adaptadores compilados en el binario durante el MVP.

## Consecuencias

- universalidad sin `switch` monolítico;
- contribuciones revisables;
- sin ejecución de plugins no confiables;
- cada release incorpora adaptadores;
- plugins fuera de proceso se evaluarán después.
