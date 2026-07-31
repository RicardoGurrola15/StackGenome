# ADR-0008: Procesamiento Local de Archivos JSON y Fingerprint Remoto

- **Estado**: aceptado
- **Fecha**: 2026-07-31
- **Fase**: Fase 18B

## Contexto
Los reportes JSON generados por el CLI de StackGenome pueden contener información sensible como evidencias, fragmentos de manifiestos, versiones y rutas absolutas locales del usuario (`/Users/ricardo/...`). 

## Decisión
1. **Importación 100% Local**: La interfaz web carga el archivo `project_graph.json` usando la File API del navegador (`file.text()`). El archivo jamás se sube (upload) a ningún servidor.
2. **Sanitización (Fingerprint)**: Al usuario se le ofrece la opción de "Obtener recomendaciones actualizadas". Al aceptar, el navegador aplica una extracción estricta:
   - Elimina la clave `evidences` por completo.
   - Elimina las propiedades sensibles de `properties` (paths, strings irreconocibles) y permite únicamente un *allowlist* de claves (ej. `workspace`).
   - Envía únicamente la lista de nodos (id, type, name, confidence, version) y aristas.
3. **Opt-in Explicito**: La llamada a red (`POST /api/v1/recommendations`) solo ocurre si el usuario lo solicita explícitamente haciendo clic en un botón que explica exactamente qué se envía.

## Consecuencias
- Se refuerza la promesa de privacidad local-first del producto.
- La carga del grafo es instantánea al no depender de la red para renderizar el archivo local.
- La UI web debe poder operar en modo "desconectado" o estático para el dashboard inicial.
