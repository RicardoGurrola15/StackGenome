# Privacidad y seguridad

## 1. Modelo de amenaza

Activos:

- código y propiedad intelectual;
- secretos;
- estructura interna;
- dependencias;
- información del entorno;
- identidad del usuario;
- integridad del catálogo;
- confianza de recomendaciones.

Adversarios:

- repositorio malicioso;
- symlinks;
- archivos gigantes;
- manifests diseñados para explotar parsers;
- recursos catalogados maliciosos;
- submissions falsas;
- fuente externa comprometida;
- atacante de red;
- dependencia del propio CLI comprometida.

## 2. Política de lectura

Predeterminado:

- nombres de archivo;
- metadatos;
- manifests y lockfiles conocidos;
- archivos de configuración allowlisted;
- contenido limitado necesario para detectar.

Excluir:

- `.env*`;
- claves;
- certificados;
- `.git/objects`;
- caches;
- build outputs;
- vendor cuando no se necesite;
- archivos binarios;
- archivos sobre límite;
- rutas ignoradas.

La exclusión no debe depender únicamente de `.gitignore`.

## 3. Symlinks y traversal

- resolver rutas de forma segura;
- nunca salir de la raíz;
- no seguir symlink externo;
- detectar ciclos;
- normalizar antes de validar;
- pruebas específicas para `../`, Unicode y case sensitivity.

## 4. No ejecución

StackGenome no ejecuta:

- package scripts;
- Gradle tasks;
- Maven plugins;
- `go generate`;
- `build.rs`;
- setup.py;
- shell hooks;
- binaries del proyecto.

Si un analizador necesita una herramienta externa, debe ser un modo explícito, aislado y posterior.

## 5. Secret scanning preventivo

Antes de serializar Fingerprint:

1. allowlist de campos;
2. normalización;
3. detección de patrones sensibles;
4. rechazo ante coincidencias;
5. preview;
6. consentimiento.

No confiar en redacción posterior de un payload ya construido.

## 6. Entorno local

El modo `--environment`:

- separado;
- desactivado;
- lista cada comando antes de ejecutar;
- usa comandos de versión conocidos;
- timeout;
- sin shell interpolation;
- captura solo stdout esperado;
- sanitiza rutas.

## 7. Red

- TLS;
- host allowlist;
- timeout;
- tamaño máximo;
- no redirects inesperados;
- user agent identificable;
- retry acotado;
- ningún upload automático.

## 8. Catálogo

Cada recurso muestra:

- fuente;
- fecha de actualización;
- versión;
- licencia;
- señales de mantenimiento;
- señales de seguridad;
- nivel de verificación.

“No se conocen vulnerabilidades” no significa “es seguro”.

## 9. Dependencias del CLI

- revisión de licencia;
- pin;
- checksum vía módulos;
- Dependabot/Renovate a decidir;
- `govulncheck`;
- SBOM del release;
- firma y checksums;
- provenance de build posterior.

## 10. Privacidad operativa

No registrar:

- contenido de manifests completos;
- rutas absolutas;
- username;
- hostname;
- tokens;
- remote privado;
- fingerprint completo en backend.

## 11. Reporte de vulnerabilidades

Antes de publicación:

- `SECURITY.md`;
- canal privado;
- política de versiones soportadas;
- proceso de revocación del catálogo;
- advisory y release corregido.

## 12. Criterio de bloqueo

Una fase se bloquea si:

- el fingerprint filtra un secreto en fixtures;
- existe path traversal;
- se ejecuta código no autorizado;
- la serialización no es determinista;
- una dependencia crítica no puede auditarse;
- el agente intenta hacer fallback de herramientas a la partición principal.
