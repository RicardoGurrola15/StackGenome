# Índice y autoridad documental

Versión documental: **0.1.0**  
Fecha base: **2026-07-24**

## Orden de lectura

| Orden | Documento | Propósito |
|---:|---|---|
| 1 | `01_PRODUCT_CHARTER.md` | Visión, usuarios, alcance y principios |
| 2 | `02_TECHNICAL_ARCHITECTURE.md` | Arquitectura del CLI y límites de componentes |
| 3 | `03_PROJECT_GRAPH_SPEC.md` | Modelo normalizado, evidencia y confianza |
| 4 | `04_CLI_SPECIFICATION.md` | Comandos, UX, salidas y códigos de error |
| 5 | `05_BACKEND_AND_DATA.md` | Backend futuro, API, catálogo e indexación |
| 6 | `06_PRIVACY_SECURITY.md` | Privacidad, amenazas y controles |
| 7 | `07_ROADMAP_BLUEPRINT.md` | Fases, metas y criterios de salida |
| 8 | `08_AGENT_EXECUTION_PROTOCOL.md` | Cómo deben trabajar Codex y Antigravity |
| 9 | `09_PROMPT_LIBRARY.md` | Prompts por fase |
| 10 | `10_CONTEXT_MEMORY_ENGRAM.md` | Continuidad entre agentes y memoria |
| 11 | `11_DEVELOPMENT_ENVIRONMENT_MACOS.md` | Herramientas en `/Volumes/intento1` |
| 12 | `12_TESTING_QUALITY.md` | Estrategia de calidad |
| 13 | `13_DECISION_LOG.md` | Decisiones menores y temas abiertos |
| 14 | `14_REFERENCES.md` | Referencias oficiales |

## Documentos operativos

- `.project/CURRENT_PHASE.md`: fase autorizada.
- `.project/STATE.md`: estado estable del proyecto.
- `.project/HANDOFF.md`: último punto de continuidad.
- `.project/BACKLOG.md`: ideas no autorizadas para la fase actual.
- `AGENTS.md`: reglas obligatorias.
- `GEMINI.md`: contexto local para Gemini/Antigravity.
- `.agents/`: skills y workflows locales.

## Regla de autoridad

Los documentos descriptivos no sustituyen el estado operativo. Un agente debe consultar siempre `.project/CURRENT_PHASE.md` antes de trabajar.

## Actualización

Toda fase debe revisar si alteró:

- contratos;
- modelos;
- comandos;
- dependencias;
- decisiones arquitectónicas;
- amenazas;
- criterios de pruebas.

Si los alteró, debe actualizar la documentación en el mismo cambio.
