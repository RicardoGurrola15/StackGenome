# Política de Privacidad — StackGenome Web

## Principio Rector

**StackGenome es local-first.** El análisis de tu repositorio ocurre en tu máquina.
El frontend web existe para visualizar y explorar el resultado de ese análisis, no para subir tu código.

---

## Modo 1: Procesamiento Local (Sin conexión a internet)

El navegador puede:
- **Abrir** un archivo `.json` exportado localmente por el CLI (`stackgenome analyze --json`).
- **Validar** su schema contra el contrato definido en `docs/api/openapi.yaml`.
- **Visualizar** el ProjectGraph completo: nodos, aristas, ecosistemas, frameworks, evidencias.
- **Generar** resúmenes locales y mostrar las recomendaciones ya incluidas en el JSON exportado.

**Lo que no ocurre:**
- No se sube el archivo automáticamente.
- No se abre ninguna conexión de red.
- No se almacena nada fuera del navegador.

---

## Modo 2: Consulta Remota (Solo con consentimiento explícito)

Cuando el usuario hace clic en "Obtener recomendaciones del catálogo en línea":

1. **Se le muestra** exactamente qué datos se enviarán (previsualización del Fingerprint antes del envío).
2. **Se construye** un Fingerprint sanitizado: solo `nodes` y `edges` con metadatos públicos.
3. **Se omite**:
   - Evidencias con rutas absolutas.
   - Propiedades marcadas como privadas.
   - El nombre del proyecto.
   - Entorno local (OS, paths, herramientas instaladas).
4. **El backend NUNCA almacena** el payload completo del Fingerprint.
5. **El backend retorna** únicamente recomendaciones de herramientas.

### Lo que se envía (ejemplo):

```json
{
  "schema_version": "1.0.0",
  "fingerprint": {
    "nodes": [
      { "id": "lang_go", "type": "language", "name": "Go", "confidence": 1.0 }
    ],
    "edges": []
  },
  "limit": 3
}
```

### Lo que NO se envía:

- Nombres de archivos específicos de tu proyecto.
- Rutas absolutas.
- Contenido de código fuente.
- Secretos o credenciales.
- Nombre del directorio raíz.
- Información personal o del sistema.

---

## Frontera Explícita Local vs. Remoto

| Acción | Local | Remoto (opt-in) |
|---|---|---|
| Visualizar el grafo completo | ✅ | — |
| Ver evidencias de detección | ✅ | — |
| Ver rutas de archivos | ✅ | ❌ Sanitizadas |
| Ver nombre del proyecto | ✅ | ❌ Omitido |
| Consultar catálogo completo | ✅ (offline) | ✅ (live) |
| Obtener recomendaciones IA-ranked | ❌ | ✅ (con confirmación) |

---

## Datos del Catálogo

El catálogo de herramientas es completamente público.
Contiene únicamente nombres, descripciones y URLs de herramientas de desarrollo open-source.
No incluye información personal ni de los usuarios.

---

## Retención de Datos en el Backend

| Dato | Almacenado | Retención |
|---|---|---|
| Fingerprint recibido | ❌ No | — |
| Recomendación devuelta | ❌ No | — |
| Métricas de uso (agregadas) | Futura decisión | TBD |
| Catálogo de herramientas | ✅ Sí (D1) | Indefinida |
