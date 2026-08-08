# Registro de decisiones

## Aceptadas

| ID | Decisión | Estado |
|---|---|---|
| D-001 | Nombre StackGenome | aceptada |
| D-002 | CLI primero, web después, app al final | aceptada |
| D-003 | CLI en Go | aceptada |
| D-004 | Arquitectura universal y modular | aceptada |
| D-005 | Profundidad progresiva por ecosistema | aceptada |
| D-006 | Local-first, metadata-only | aceptada |
| D-007 | Escaneo de entorno separado y opt-in | aceptada |
| D-008 | Sin IA dentro del MVP | aceptada |
| D-009 | Ranking determinista | aceptada |
| D-010 | Top 3 por defecto, `--more` opcional | aceptada |
| D-011 | Backend futuro Cloudflare Workers + D1 | aceptada preliminar |
| D-012 | ProjectGraph interno y Fingerprint enviable | aceptada |
| D-013 | PURL/CycloneDX como estándares de interoperabilidad | aceptada |
| D-014 | Adaptadores compilados en MVP | aceptada |
| D-015 | Documentación Git como fuente de verdad; Engram auxiliar | aceptada |
| D-016 | Deduplicación de nodos de plataforma (iOS/Android) | aceptada |
| D-017 | Migración de 'languages' a 'ecosystem' en catálogo | aceptada |
| D-018 | Exclusión estricta de rutas de vendor/build (ios/Pods) | aceptada |

## Pendientes

- licencia del repositorio;
- librería CLI: estándar vs Cobra;
- adopción exacta de go-enry;
- grado de reutilización de Syft/Trivy/OSV Scanner;
- catálogo semilla y criterios editoriales;
- identidad visual;
- dominio;
- política futura de telemetría;
- scope exacto de C/C++;
- tipo de datos de plugins/skills/agentes en el catálogo;
- si `.engram/` se versionará públicamente.

| D-019 | Despliegue web vía Cloudflare Pages (no Workers Static Assets) | aceptada — Fase 21A |
| D-020 | `requires_context` como mecanismo de gate contextual en el scorer | aceptada — Fase 21A |

## Regla

Una decisión pendiente no autoriza implementación. Las decisiones arquitectónicas importantes requieren ADR.

---

## ADR-0010 — Despliegue Web: Cloudflare Pages vs Workers Static Assets

**Fecha:** 2026-08-07
**Estado:** Aceptada

### Contexto

La interfaz web de StackGenome es una aplicación Next.js compilada como `static export` (output: 'export'). Actualmente no está desplegada. Debemos elegir entre dos plataformas de Cloudflare para alojarla:

**Opción A — Cloudflare Pages**
- Plataforma nativa para sitios estáticos.
- Integración directa con GitHub (push → build → deploy).
- Preview automático por PR.
- Dominio custom gratuito.
- Límites generosos (100K requests/día en plan Free).
- Workers Functions opcionales como complemento.

**Opción B — Cloudflare Workers Static Assets**
- Relativamente nueva (GA: 2025).
- Requiere Wrangler + configuración manual para servir assets estáticos.
- Más flexible si se quiere unificar todo en un Worker (API + frontend).
- Menor madurez de tooling y documentación respecto a Pages.
- Sin prebuilds automáticos desde GitHub sin configuración adicional.

### Decisión

**Se mantiene Cloudflare Pages** para el despliegue de la interfaz web.

### Justificación

1. La web actual es un `static export` puro que encaja perfectamente con el modelo de Pages.
2. Pages ofrece preview deployments por PR sin configuración adicional — útil para Alpha.
3. Workers Static Assets no ofrece ventaja arquitectónica significativa para nuestro caso de uso (frontend 100% estático + API separada en Worker).
4. El backend ya existe como Worker independiente. Combinarlo con el frontend en un sólo Worker añadiría complejidad sin beneficio.
5. No migramos por moda: Pages sigue siendo la opción más simple, documentada y rápida de activar.

### Consecuencias

- Fase 21D desplegará el frontend via `wrangler pages deploy` o la integración Git de Cloudflare Pages.
- El backend Worker existente permanece independiente y se comunica vía CORS.
- CORS headers ya están configurados en el Worker backend.

