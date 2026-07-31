# Prompt inicial para Codex o Antigravity

Copia desde “INICIO DEL PROMPT” hasta “FIN DEL PROMPT”.

---

## INICIO DEL PROMPT

Estás trabajando exclusivamente en el proyecto **StackGenome**, ubicado normalmente en:

`/Volumes/intento1/Repos/StackGenome`

Tu primera tarea no es construir todo el producto. Debes iniciar y completar únicamente la fase que esté declarada en `.project/CURRENT_PHASE.md`.

Antes de actuar:

1. Confirma que el repositorio abierto es StackGenome.
2. Lee íntegramente:
   - `AGENTS.md`
   - `GEMINI.md` si tu herramienta lo utiliza
   - `docs/00_INDEX.md`
   - `docs/01_PRODUCT_CHARTER.md`
   - `docs/02_TECHNICAL_ARCHITECTURE.md`
   - `docs/07_ROADMAP_BLUEPRINT.md`
   - `docs/08_AGENT_EXECUTION_PROTOCOL.md`
   - `.project/CURRENT_PHASE.md`
   - `.project/STATE.md`
   - `.project/HANDOFF.md`
3. Revisa los ADR aceptados en `docs/adr/`.
4. Si Engram está disponible, busca memoria del proyecto `StackGenome`; úsala solo para complementar y contrástala con los archivos.
5. Verifica que `/Volumes/intento1` esté montado. No instales ni escribas cachés en la partición principal como fallback.
6. Inspecciona el estado actual del repositorio y Git sin modificar nada.
7. Presenta un plan breve para **la fase activa solamente**, incluyendo:
   - archivos que planeas crear o modificar;
   - comandos que planeas ejecutar;
   - dependencias o programas que planeas instalar;
   - ubicación exacta donde se instalarán;
   - criterios de aceptación que verificarás.

Reglas obligatorias:

- No avances a otra fase.
- No cambies decisiones arquitectónicas aceptadas sin proponer un ADR.
- No uses IA dentro del producto MVP.
- No diseñes para un stack personal; el analizador es universal y modular.
- No ejecutes archivos o scripts del proyecto que StackGenome esté analizando.
- No envíes código fuente ni información sensible.
- No uses `sudo`, no edites configuración global y no realices commits sin autorización.
- Si necesitas instalar o configurar GentleAI/Engram para este workspace, primero detecta el estado actual, explica qué archivos globales o locales tocaría y solicita autorización antes de ejecutar cambios.
- Tras completar la fase: ejecuta todas las pruebas, actualiza `.project/STATE.md`, `.project/CURRENT_PHASE.md` y `.project/HANDOFF.md`, registra decisiones y guarda una memoria breve en Engram si está disponible.
- Después presenta evidencia y pregunta exactamente: **“¿Autorizas avanzar a la fase siguiente?”**
- Detente y espera la respuesta. No inicies la siguiente fase en el mismo turno.

Comienza ahora con la lectura y diagnóstico de la fase activa.

## FIN DEL PROMPT
