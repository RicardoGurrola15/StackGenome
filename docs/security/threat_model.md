# Modelo de Amenazas (Threat Model) - StackGenome

Fecha de revisión: 2026-07-30
Estado: Mitigaciones aplicadas y validadas (Fase 14).

## 1. Privacidad de Datos Locales

**Amenaza:** Fuga de código fuente o secretos mediante telemetría inadvertida.
**Mitigación:** 
- Arquitectura *Local-first*. StackGenome realiza su análisis íntegramente de manera offline en la máquina del usuario.
- **Anonimización Forzada**: Antes de invocar a la API remota (`--remote`), se ejecuta `sanitizer.Anonymize()` que remueve todo tipo de propiedades confidenciales, rutas absolutas y evidencias completas, enviando únicamente una lista de identificadores genéricos de lenguajes o herramientas.

## 2. Directory Traversal y Symlink Attacks

**Amenaza:** Repositorios maliciosos (por ejemplo, descargados de internet) podrían contener enlaces simbólicos apuntando a `/etc/passwd` o fuera del directorio del proyecto, provocando que el CLI lea archivos del sistema.
**Mitigación:**
- `SafeWalker` (implementado en `internal/fs/security.go`) intercepta cualquier resolución de `Lstat`/`Stat`.
- Los tests de fuzzing (`security_fuzz_test.go`) validan que no se escape de la raíz del directorio base durante el recorrido de archivos, sin ocasionar *panics*.

## 3. Pánicos en el Parser (Denial of Service)

**Amenaza:** Ficheros de manifiestos (`package.json`, `go.mod`, etc.) malformados a propósito para causar un cuelgue, un bucle infinito o un pánico de memoria (`out of bounds`).
**Mitigación:**
- Se implementó Fuzz Testing (ej: `node_detector_fuzz_test.go`) sometiendo los analizadores principales a millones de entradas mutadas.
- Los parsers fallan de forma grácil (`err != nil`) ignorando los archivos que no pueden ser descifrados en lugar de colapsar la aplicación completa.

## 4. Dependencias Envenenadas (Supply Chain)

**Amenaza:** Vulnerabilidades introducidas mediante paquetes de terceros en Go o Node.
**Mitigación:**
- Restricción estricta del proyecto: **0 dependencias externas en el core de Go**.
- Verificación en CI/local con `govulncheck` de la librería estándar y la toolchain de Go, certificando 0 hallazgos.

## 5. Man-in-the-Middle y Backend Compromise

**Amenaza:** Intercepción de la comunicación HTTP o un backend Cloudflare vulnerado.
**Mitigación:**
- Todas las peticiones `--remote` viajan cifradas sobre TLS (HTTPS) hacia dominios `.workers.dev`.
- Incluso si el backend fuera vulnerado, el payload enviado por el CLI carece de metadatos útiles para la explotación, y el CLI en Go está aislado y parsea la respuesta JSON limitando los campos al DTO de `Recommendation`. No se ejecuta código remoto en la máquina del desarrollador.
