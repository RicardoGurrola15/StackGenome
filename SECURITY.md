# Política de Seguridad de StackGenome

## Versiones Soportadas

Actualmente, solo la rama `main` y el release `Alpha` más reciente reciben parches de seguridad.

## Reporte de Vulnerabilidades

Si descubres una vulnerabilidad de seguridad en StackGenome (especialmente fugas de datos locales hacia orígenes de red no deseados o fallas en el `SafeWalker` que permitan *directory traversal*), te pedimos que la reportes de inmediato.

Por favor, no abras un issue público en GitHub de inmediato. Envía un correo electrónico a `security@stackgenome.com` con los detalles. Intentaremos responder y publicar un parche lo más rápido posible.

## Modelo de Amenazas

Como parte de nuestro compromiso con la seguridad, mantenemos un [Threat Model](docs/security/threat_model.md) vivo. StackGenome:
1. Opera **offline** por defecto.
2. Cuando se usa `--remote`, **anonimiza** obligatoriamente todos los payloads, borrando rutas locales, nombres de paquetes sensibles y metadatos comprometedores.
3. No requiere ejecución de código arbitrario ni binarios (parsing 100% estático).
