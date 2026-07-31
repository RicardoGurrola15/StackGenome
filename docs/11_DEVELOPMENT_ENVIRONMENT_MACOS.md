# Entorno de desarrollo en macOS

## 1. Rutas

```text
/Volumes/intento1/
├── programas/
│   ├── mise/
│   └── go-workspace/
├── caches/
│   ├── mise/
│   └── go/
└── Repos/
    └── StackGenome/
```

Estas rutas son de desarrollo, no del producto.

## 2. Estrategia

Se utiliza `mise` para instalar la versión fijada de Go en el volumen externo.

Variables cargadas por `scripts/activate-env.sh`:

- `MISE_INSTALL_PATH`;
- `MISE_DATA_DIR`;
- `MISE_CACHE_DIR`;
- `MISE_INSTALLS_DIR`;
- `GOPATH`;
- `GOMODCACHE`;
- `GOCACHE`;
- `GOTMPDIR`;
- `XDG_CACHE_HOME`.

No se define globalmente `TMPDIR` para evitar afectar otras aplicaciones.

## 3. Activación

```bash
cd /Volumes/intento1/Repos/StackGenome
source scripts/activate-env.sh
```

La activación:

- comprueba el volumen;
- crea directorios controlados;
- añade el binario y shims a PATH;
- no edita shell rc.

## 4. Bootstrap

```bash
./scripts/bootstrap-macos.sh
```

El script:

1. verifica macOS;
2. verifica el volumen;
3. descarga `mise` a `/Volumes/intento1/programas/mise/bin/mise` si falta;
4. instala Go definido en `mise.toml`;
5. muestra rutas;
6. no usa sudo.

La descarga requiere revisar el script y confiar en `https://mise.run`. Para máxima seguridad se puede descargar/verificar manualmente siguiendo la documentación oficial.

## 5. Versión inicial

El documento fija Go `1.26.5` por reproducibilidad al 24 de julio de 2026. Actualizar solo mediante cambio separado con pruebas.

## 6. Persistencia opcional

Para no ejecutar `source` manualmente se podría añadir una línea a `~/.zshrc`, pero eso es configuración global y requiere autorización:

```bash
source /Volumes/intento1/Repos/StackGenome/scripts/activate-env.sh
```

No se recomienda si el volumen no siempre está conectado. Alternativa: alias explícito o terminal del workspace.

## 7. IDE

Configurar el IDE para abrir:

```text
/Volumes/intento1/Repos/StackGenome
```

El proceso del IDE debe heredar las variables si se necesitan las herramientas de `mise`. Se puede iniciar desde una terminal activada.

## 8. Verificación

```bash
./scripts/verify-environment.sh
```

Debe mostrar todas las rutas bajo `/Volumes/intento1`.

## 9. CI

CI no usa estas rutas. Usa `mise.toml` o setup oficial de Go. Nunca codificar rutas absolutas en Go, tests o scripts compartidos de build.
