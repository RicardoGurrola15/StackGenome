# Guía de Contribución para StackGenome

¡Gracias por tu interés en contribuir a StackGenome!

## Entorno de Desarrollo Local

El proyecto está diseñado para funcionar de forma completamente local usando `mise` y `Go`.
Para configurar tu entorno:

1. Clona el repositorio.
2. Asegúrate de tener `mise` instalado.
3. Ejecuta `mise install` o usa el script interactivo:
   ```bash
   source ./scripts/activate-env.sh
   ```
4. Ejecuta los tests para validar que tu entorno está listo:
   ```bash
   go test ./...
   ```

## Cómo añadir nuevos detectores de Lenguajes o Frameworks

El corazón de StackGenome es su analizador estático sin ejecución.
Si quieres añadir soporte para un nuevo lenguaje (por ejemplo, Elixir o Haskell):

1. Navega a `internal/detectors/language/` y crea un archivo (ej: `elixir_detector.go`).
2. Implementa la interfaz `FileDetector`:
   - `Handles(filename string) bool`: Retorna `true` para archivos relevantes como `mix.exs`.
   - `Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error)`: Parsea el archivo estáticamente y extrae nodos de tipo `language` y `dependency`.
3. Registra tu nuevo detector en `internal/detectors/registry.go` (o deja que el sistema lo registre por reflexión si se expande la API).
4. **Muy importante**: No uses dependencias de terceros si es posible usar la librería estándar. Tampoco añadas binarios nativos del lenguaje a las dependencias. El parsing debe ser 100% estático.

## Política de Privacidad (Zero-Telemetry)

Cualquier PR que añada recolección de datos, telemetría o intente saltarse el `sanitizer.Anonymize` será rechazado.

¡Esperamos tus Pull Requests!
