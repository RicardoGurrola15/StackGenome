# Estrategia de pruebas y calidad

## 1. Pirámide

- unitarias: parsers, reglas, score;
- componentes: analizadores con FS virtual;
- integración: fixtures;
- end-to-end: binario sobre repos sintéticos;
- compatibilidad: OS/arquitectura;
- seguridad: fuzzing y adversarial.

## 2. Fixtures

Cada fixture:

- mínima;
- sintética;
- licencia clara;
- sin secretos reales;
- propósito documentado;
- expected output estable.

Estructura prevista:

```text
testdata/fixtures/<ecosystem>/<case>/
```

## 3. Golden tests

Adecuados para:

- ProjectGraph;
- Fingerprint;
- CLI JSON;
- reportes.

Los snapshots deben revisarse; no actualizarse automáticamente para “hacer pasar” tests.

## 4. Parsers

Para cada parser:

- archivo válido;
- vacío;
- truncado;
- tipos incorrectos;
- Unicode;
- tamaño máximo;
- campos desconocidos;
- versiones/rangos extraños.

## 5. Seguridad

- path traversal;
- symlink externo;
- ciclo;
- race replacement;
- zip bomb si se añaden archivos;
- secret patterns;
- shell injection;
- host allowlist.

## 6. Determinismo

Ejecutar múltiples veces y comparar bytes de JSON canónico. No depender de orden de mapas o filesystem.

## 7. Rendimiento

Benchmarks:

- 1k, 10k y 100k archivos sintéticos;
- monorepo;
- manifest grande;
- pocos archivos grandes;
- cancelación.

Los objetivos se calibran con mediciones, no con optimización prematura.

## 8. Comandos

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
go test -fuzz=<target>
govulncheck ./...
```

`govulncheck` se instala en herramientas del proyecto y no se asume global.

## 9. CI

Matriz prevista:

- macOS;
- Linux;
- Windows;
- Go fijado;
- tests;
- race donde esté soportado;
- build;
- lint;
- SBOM/release en tags.

## 10. Definition of Done

- criterios de fase satisfechos;
- pruebas;
- docs;
- no warnings ocultos;
- diff revisado;
- privacidad revisada;
- estado actualizado;
- autorización solicitada.
