# 6. Rechazo de dependencia externa go-enry

Date: 2026-07-24

## Status

Accepted

## Context

El roadmap (Fase 4) sugiere "integración evaluada con go-enry o alternativa" para detectar lenguajes analizando las extensiones y el contenido de todos los archivos del repositorio, similar a lo que hace GitHub Linguist. Sin embargo, una de las reglas primarias del proyecto StackGenome (Fase 1) es "Cero dependencias externas" o minimizarlas al extremo, construyendo el core solo con la Standard Library de Go. Incorporar `go-enry` implicaría importar un analizador pesado y múltiples dependencias transitivas (CGO o diccionarios grandes).

## Decision

Se rechaza la integración de `go-enry`. 
En su lugar, StackGenome implementa una detección basada principalmente en *manifests* (que aportan alta señal sobre el ecosistema) y se complementa con un pequeño detector heurístico nativo (`ExtensionDetector`) que mapea un conjunto restringido de extensiones de archivo a sus lenguajes, sin añadir ninguna dependencia externa al binario.

## Consequences

- Mantenemos el ejecutable pequeño y libre de dependencias.
- StackGenome no detectará con precisión del 100% todos los lenguajes oscuros del mundo.
- Se simplifica el escaneo (menor I/O) al no tener que abrir cada archivo para análisis léxico.
