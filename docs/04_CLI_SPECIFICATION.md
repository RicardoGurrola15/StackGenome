# Especificación del CLI

Ejecutable: `stackgenome`

## 1. Principios de UX

- seguro por defecto;
- sin envío implícito;
- salida breve por defecto;
- JSON estable para automatización;
- errores accionables;
- no requiere cuenta para análisis local;
- no instala recomendaciones.

## 2. Comandos previstos

### `stackgenome inspect [path]`

Analiza localmente y genera ProjectGraph.

Flags previstos:

```text
--format terminal|json|cyclonedx
--output <file>
--depth <n>
--max-file-size <bytes>
--include <glob>
--exclude <glob>
--no-cache
--verbose
--debug
```

### `stackgenome fingerprint [path]`

Genera la proyección sanitizada.

```text
--preview
--output <file>
--include-project-name
--environment <none|declared|local>
```

`--environment local` requiere confirmación interactiva o `--yes` explícito.

### `stackgenome recommend [path]`

```text
--goal explore|security|testing|performance|observability|deployment|editor|architecture
--catalog local|remote
--more
--limit <n>
--json
```

Predeterminado: tres recomendaciones.

### `stackgenome doctor`

Comprueba:

- versión;
- permisos;
- conectividad opcional;
- configuración;
- catálogo;
- rutas de caché;
- toolchain local;
- Engram no forma parte del producto y no se comprueba aquí.

### `stackgenome catalog`

Futuro:

- status;
- update;
- search;
- sources.

### `stackgenome completion`

Genera completions para zsh, bash, fish y PowerShell.

## 3. Interacción de consentimiento

Antes de enviar:

```text
StackGenome enviará:
- ecosystems y versiones
- dependencias públicas y versiones
- plataformas
- categorías de herramientas
- objetivo seleccionado

StackGenome NO enviará:
- código
- rutas absolutas
- secretos
- contenido del README
- nombre del proyecto, salvo autorización

Ver fingerprint completo? [Y/n]
Enviar? [y/N]
```

Predeterminado: no.

## 4. Salida recomendada

```text
1. resource/name
   Fit: 94
   Compatible: yes
   Security: no known critical advisories
   Maintenance: active
   License: Apache-2.0
   Why: matches Go 1.26 and your CLI/testing goal
   URL: ...

2. ...

3. ...
```

Debe distinguir:

- dato verificado;
- dato no disponible;
- inferencia;
- advertencia.

## 5. Códigos de salida

| Código | Significado |
|---:|---|
| 0 | éxito |
| 1 | error general |
| 2 | argumentos/configuración inválidos |
| 3 | raíz no válida o no accesible |
| 4 | violación de política de seguridad |
| 5 | análisis parcial con modo estricto |
| 6 | error de red |
| 7 | catálogo incompatible/no disponible |
| 8 | fingerprint rechazado |
| 9 | cancelado por usuario |

## 6. Configuración

Orden de precedencia:

1. flags;
2. variables `STACKGENOME_*`;
3. config de proyecto `.stackgenome.toml`;
4. config de usuario;
5. defaults.

La configuración de proyecto no puede relajar silenciosamente restricciones de seguridad.

## 7. Telemetría

- desactivada;
- no se implementa en MVP;
- cualquier futuro analytics será opt-in y documentado.

## 8. Accesibilidad

- modo sin color;
- mensajes no dependientes solo de color;
- salida legible por screen readers;
- `--quiet`;
- JSON para herramientas.
