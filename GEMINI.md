# GEMINI.md — Contexto local de StackGenome

Este contexto aplica exclusivamente al repositorio StackGenome.

## Instrucción principal

Lee y obedece `AGENTS.md`. Después consulta:

- `docs/00_INDEX.md`
- `.project/CURRENT_PHASE.md`
- `.project/STATE.md`
- `.project/HANDOFF.md`

Para Antigravity, consulta también:

- `.agents/agents.md`
- `.agents/workflows/`
- `.agents/skills/`

## Reglas esenciales

- Ejecuta solamente la fase activa.
- No avances sin autorización explícita.
- StackGenome es universal; no está diseñado alrededor de Flutter ni de un stack personal.
- CLI en Go; backend futuro en TypeScript sobre Cloudflare.
- MVP sin IA dentro del producto.
- Análisis local-first y metadata-only.
- El escaneo del entorno es opcional.
- No ejecutes código del proyecto inspeccionado.
- No instales herramientas fuera de `/Volumes/intento1` cuando sean configurables.
- Si `/Volumes/intento1` no está montado, detente.
- Prueba y documenta todos los cambios.
- Engram es memoria auxiliar; los archivos del repositorio son la autoridad.

Al terminar una fase, presenta evidencia y pregunta si se autoriza continuar. No continúes en el mismo turno.
