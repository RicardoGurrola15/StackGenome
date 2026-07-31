# Estado del Proyecto: StackGenome

- **Fase actual**: 18B (Web Alpha Completada)
- **Bloqueos**: Ninguno.
- **Siguiente fase pendiente de autorización**: Fase 19 (Validación con proyecto externo Territor).

## Fases completadas
- [x] Fase 0-15: Core, analizadores, y pipeline de CI.
- [x] Fase 16: Auditoría Alpha.
- [x] Fase 17: CI/CD.
- [x] Fase 18A: Web Readiness Gate (CI limpio, contratos definidos).
- [x] Fase 18B: Web Alpha (Next.js App Router, SSG Catalog, local File API para importación).

## Notas técnicas actuales
- Frontend implementado en `web/` usando Next.js 14. 
- Componentes clave: `FileImporter` procesa localmente. `StackGraph` usa SVG + `d3-force` para grafos livianos. 
- API calls al Worker apuntan a staging por defecto (ver `web/lib/api.ts`).
- Contrato JSON tipado contra `docs/api/openapi.yaml`.
