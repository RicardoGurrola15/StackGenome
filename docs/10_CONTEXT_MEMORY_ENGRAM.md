# Contexto, memoria y Engram

## 1. Fuente de verdad

La continuidad oficial vive en Git:

- `AGENTS.md`;
- documentación;
- ADR;
- `.project/`.

Engram acelera recuperación, pero puede estar ausente, desactualizado o contener observaciones antiguas.

## 2. Separación del proyecto

Nombre canónico de memoria:

```text
StackGenome
```

Engram agrupa proyectos por nombre/remote. Si aparecen variantes, revisar:

```bash
engram projects list
engram projects consolidate
```

No consolidar automáticamente sin revisar.

## 3. Configuración

Codex, si no está configurado:

```bash
engram setup codex
```

Antigravity CLI, si se utiliza y no está configurado:

```bash
engram setup antigravity-cli
```

Estos comandos pueden modificar configuración del agente a nivel de usuario. El agente debe diagnosticar y pedir autorización antes de ejecutarlos.

Para un workspace compatible con VS Code se incluye `.vscode/mcp.json`:

```json
{
  "servers": {
    "engram": {
      "command": "engram",
      "args": ["mcp"]
    }
  }
}
```

El binario puede ser global, pero las observaciones deben identificarse con StackGenome.

## 4. Exportación local

Después de revisar memoria:

```bash
engram sync
```

Esto puede exportar memoria a `.engram/` para compartir mediante Git. Antes de commit:

- revisar secretos;
- revisar datos personales;
- eliminar ruido;
- confirmar que pertenece a StackGenome.

En otra máquina:

```bash
engram sync --import
```

## 5. Qué guardar

- decisiones aprobadas;
- razones de arquitectura;
- bugs y causa raíz;
- comandos que resolvieron un problema;
- limitaciones;
- resultados de benchmarks;
- fase completada;
- siguiente paso;
- nombres de archivos relevantes.

## 6. Qué no guardar

- tokens;
- API keys;
- `.env`;
- datos personales;
- contenido extenso de código;
- payloads privados;
- rutas de otros proyectos;
- información clínica o ajena al proyecto;
- conversaciones no técnicas.

## 7. Inicio de sesión del agente

1. leer archivos;
2. buscar en Engram:
   - `StackGenome current phase`;
   - `StackGenome architecture decisions`;
   - tema específico;
3. contrastar fechas y ADR;
4. ignorar memoria contradictoria;
5. registrar conflicto si importa.

## 8. Cierre de sesión

Guardar una observación compacta:

```text
Project: StackGenome
Phase:
Completed:
Verified:
Decisions:
Known issues:
Next authorized action:
Files:
```

Actualizar `.project/HANDOFF.md` aunque Engram funcione.

## 9. GentleAI

GentleAI puede instalar skills y flujos con alcance workspace:

```bash
gentle-ai install --scope=workspace
```

No ejecutar automáticamente. Antes:

- crear branch o snapshot;
- revisar qué archivos generará;
- preservar `AGENTS.md`, `GEMINI.md`, `.agents/` y `.project/`;
- evitar duplicar reglas;
- decidir si se usarán sus flujos SDD o los de StackGenome.

Codex y Antigravity pueden operar como agentes de una sola conversación/fase; el protocolo de StackGenome sigue siendo obligatorio.

## 10. Fallback

Si Engram o GentleAI fallan:

- no bloquear el proyecto;
- usar `HANDOFF.md`;
- registrar el problema;
- continuar solo dentro de la fase autorizada.
