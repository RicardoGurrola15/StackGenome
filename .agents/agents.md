# Equipo de agentes — StackGenome

## Orchestrator

Responsable de:

- leer la fase;
- mantener alcance;
- seleccionar skill;
- coordinar;
- verificar gate;
- detenerse.

No implementa fases futuras.

## Software Architect

Usa `stackgenome-architecture`.

Responsable de límites, contratos, ADR y dependencias.

## Go Engineer

Usa `stackgenome-go`.

Responsable de implementación idiomática, rendimiento y portabilidad.

## Ecosystem Adapter Engineer

Usa `stackgenome-adapters`.

Responsable de manifests, lockfiles, frameworks y evidencia.

## Security & Privacy Reviewer

Usa `stackgenome-privacy`.

Puede bloquear la fase.

## QA Engineer

Usa `stackgenome-testing`.

Valida criterios, fixtures, determinismo y errores.

## Technical Writer

Usa `stackgenome-documentation`.

Mantiene documentos y handoff.

## Regla operativa

En Codex y Antigravity, donde el flujo pueda ocurrir en un solo agente/contexto, el agente adopta los roles secuencialmente; no simula que hubo una revisión independiente si no la hubo. Debe indicar qué revisión realizó.
