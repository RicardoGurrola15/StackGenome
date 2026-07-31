# Estado del proyecto

Actualizado: 2026-07-30

## Producto

- Nombre: StackGenome
- Etapa: MVP en progreso
- CLI: Go — funcional (Fases 0-13 completas)
- Backend: Cloudflare Workers + D1 (Fase 12, completada)
- Web: futura
- App: futura

## Estado Global

`COMPLETED`

## Fases Completadas

- **Fase 0-11**: (Entorno, Core Analyzer, 20+ Detectores, Sanitización, Scoring Local, Catálogo Embebido).
- **Fase 12**: Backend Cloudflare Workers + D1 (Desplegado en staging: `stackgenome-api-staging`).
- **Fase 13**: UX final del CLI y conexión remota.
- **Fase 14**: Calidad, seguridad y rendimiento.
- **Fase 15**: Alpha (Binarios, Licencia MIT, Comunidad, Limitaciones).
- **Fase 16**: Auditoría independiente (Completada, estado ALPHA_AUDIT_CONDITIONAL superado).
- **Fase 17**: CI/CD Automático (Completada).

## Fase Activa

**Fase 17 completada**

## Pendiente Inmediato

La siguiente prioridad del proyecto (Fase 18) es iniciar el desarrollo de la **Web UI (Dashboard)** para explorar de forma visual el grafo generado por el CLI.

## Riesgos actuales

- Límites Free de Cloudflare Workers (100K req/día, 10ms CPU/req) — adecuados para staging.
- Elegir licencia antes de publicación abierta.
