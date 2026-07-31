# Skill: Go Engineering

## Reglas

- Go idiomático;
- biblioteca estándar primero;
- context para cancelación;
- errores con `%w`;
- no panics para entradas del usuario;
- orden determinista;
- concurrencia acotada;
- APIs pequeñas;
- no globals mutables;
- tests table-driven;
- fuzz para parsers;
- `gofmt`, `go vet`, race.

## Dependencias

Justificar licencia, mantenimiento y tamaño.

## Portabilidad

No asumir POSIX en dominio. Aislar diferencias de OS.
