# AGENTS.md — StackGenome

Este archivo contiene las reglas obligatorias para Codex y cualquier agente que respete `AGENTS.md`.

## 1. Fuente de verdad

Lee en este orden antes de modificar archivos:

1. `docs/01_PRODUCT_CHARTER.md`
2. `docs/02_TECHNICAL_ARCHITECTURE.md`
3. `docs/07_ROADMAP_BLUEPRINT.md`
4. `.project/CURRENT_PHASE.md`
5. `.project/STATE.md`
6. El documento específico de la fase activa
7. Los ADR aplicables en `docs/adr/`

Cuando haya conflicto:

1. Instrucción explícita actual del usuario.
2. `AGENTS.md`.
3. ADR aceptado.
4. Documento técnico especializado.
5. Roadmap.
6. Resto de la documentación.

No cambies una decisión aceptada silenciosamente. Propón un ADR nuevo que sustituya al anterior y espera autorización.

## 2. Regla de una sola fase

- Trabaja únicamente en la fase indicada por `.project/CURRENT_PHASE.md`.
- No empieces la fase siguiente aunque la actual parezca terminada.
- Al completar una fase: ejecuta verificaciones, actualiza el estado, presenta evidencia y **detente**.
- Solicita autorización explícita antes de avanzar.
- No conviertas tareas “opcionales” en trabajo obligatorio sin aprobación.

## 3. Restricciones del producto

- StackGenome debe ser universal por arquitectura y progresivo por profundidad.
- No diseñes el producto alrededor de un proyecto, lenguaje, IDE o usuario concreto.
- El CLI se implementa en Go.
- El MVP no utiliza un LLM para detectar, clasificar, recomendar ni redactar resultados.
- Las recomendaciones deben ser deterministas, trazables y basadas en evidencia.
- El análisis del proyecto es local-first.
- El escaneo del entorno es separado, voluntario y desactivado por defecto.
- Nunca envíes código fuente, secretos, rutas personales ni contenido sensible al servidor.
- No ejecutes scripts encontrados dentro del proyecto analizado.
- No cargues plugins binarios de terceros dentro del proceso durante el MVP.
- Los adaptadores iniciales son módulos compilados y revisables.

## 4. Reglas del entorno local

Rutas previstas:

```text
Volumen:    /Volumes/intento1
Programas:  /Volumes/intento1/programas
Cachés:     /Volumes/intento1/caches
Repositorio:/Volumes/intento1/Repos/StackGenome
```

- No instales SDK, gestores o cachés en la partición principal si existe una alternativa configurable.
- No edites archivos globales de shell sin autorización.
- No uses `sudo` sin explicar la necesidad y obtener autorización.
- Antes de instalar algo, muestra destino, tamaño aproximado y comando.
- Si el volumen no está montado, detente; no hagas fallback a `$HOME`.
- Usa `scripts/activate-env.sh` para cargar las rutas del proyecto.
- El repositorio puede estar en otra ruta en CI; el código nunca debe depender de una ruta absoluta.

## 5. Diseño y código

- Código, identificadores, APIs, commits y nombres técnicos: inglés.
- Documentación explicativa del proyecto: español, salvo documentación pública que posteriormente se decida traducir.
- Mantén dependencias mínimas.
- Prefiere biblioteca estándar de Go cuando sea razonable.
- Toda dependencia nueva requiere justificación, licencia y riesgo de mantenimiento.
- Evita abstracciones prematuras; crea interfaces en límites reales.
- Mantén `internal/` para implementación no pública.
- Errores con contexto; nunca ocultes errores.
- No registres secretos, rutas personales ni contenido de archivos analizados.
- El orden de resultados debe ser estable para las mismas entradas.
- Toda detección debe guardar evidencia y nivel de confianza.

## 6. Pruebas y verificación

Cada cambio funcional debe incluir pruebas adecuadas.

Antes de cerrar una fase, cuando existan esos comandos:

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
```

Añade pruebas de integración con fixtures pequeñas y sintéticas. No copies repositorios de terceros completos al historial. No uses datos personales.

Una fase no se considera completa por “compilar”. Debe cumplir todos sus criterios de aceptación.

## 7. Git y cambios

- No reescribas historia.
- No fuerces push.
- No elimines archivos ajenos al objetivo.
- Mantén commits pequeños y descriptivos cuando el usuario autorice commits.
- No hagas commit automático salvo instrucción explícita.
- Antes de cerrar, muestra archivos modificados, pruebas ejecutadas y riesgos pendientes.

## 8. Documentación viva

Actualiza:

- `.project/STATE.md`: estado técnico estable.
- `.project/CURRENT_PHASE.md`: progreso de la fase activa.
- `.project/HANDOFF.md`: resumen de sesión y siguiente acción.
- `docs/13_DECISION_LOG.md`: decisiones menores.
- `docs/adr/`: decisiones arquitectónicas importantes.

No marques tareas como terminadas sin evidencia.

## 9. Engram

Engram es memoria auxiliar, no la fuente de verdad.

Al iniciar:
- Busca recuerdos de `StackGenome` si las herramientas de Engram están disponibles.
- Contrasta lo recuperado con los archivos versionados.

Al cerrar:
- Guarda decisiones, descubrimientos, errores resueltos y siguiente paso.
- No guardes secretos ni fragmentos extensos de código.
- Ejecuta `engram sync` solo después de revisar que el contenido exportado es apropiado para el repositorio.

Si Engram no está disponible, continúa usando `.project/HANDOFF.md`.

## 10. Formato obligatorio al terminar una fase

Entrega:

1. Fase ejecutada.
2. Resultado y archivos modificados.
3. Criterios de aceptación, uno por uno.
4. Comandos y pruebas con resultado.
5. Decisiones tomadas.
6. Riesgos o deuda pendiente.
7. Estado de Engram/documentación.
8. Pregunta final: **“¿Autorizas avanzar a la fase siguiente?”**

Después, detente.
