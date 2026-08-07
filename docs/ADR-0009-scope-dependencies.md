# ADR 0009: Scope de Dependencias y Deduplicación de Plataformas

## Estado
Aceptado

## Contexto
Durante la Fase 20, enfrentamos un problema con el proyecto piloto `Territor` (Flutter/Dart). El analizador extraía dependencias del `pubspec.lock`, lo cual inundaba el `ProjectGraph` con cientos de nodos de dependencias transitivas. Asimismo, las plataformas generadas (`ios/`, `android/`, `macos/`) duplicaban las evidencias de lenguajes nativos, y los directorios de compilación (como `ios/Pods/` o `.dart_tool/`) contaminaban el análisis con código ajeno.

## Decisión
Se introducen tres estrategias fundamentales para la calidad del análisis:

1.  **Scope de Dependencias**:
    Se introduce el campo `scope` en los `NodeDTO` de tipo dependencia. El analizador ahora distingue entre dependencias `runtime` (incluidas en el binario final) y `development` (utilizadas en compilación/tests). Esto permite a los consumidores (como el Web o CLI) priorizar qué dependencias mostrar o evaluar.

2.  **Deduplicación de Plataformas y Nodos Genéricos**:
    Los analizadores ahora pueden generar evidencias redundantes desde múltiples archivos. Se ha implementado un paso post-procesamiento `DeduplicatePlatforms` en el motor principal que agrupa evidencias bajo una misma clave única (`platform:iOS/macOS Native`) si detecta fragmentación, reduciendo nodos duplicados y consolidando la información de un ecosistema nativo o framework.

3.  **Filtrado Estricto de Vendor/Build**:
    La abstracción del `walker.go` (FS abstraction) se hizo "path-aware", descartando automáticamente directorios conocidos de librerías y dependencias instaladas localmente (por ejemplo, `ios/Pods`, `.dart_tool`, `android/.gradle`). Esto garantiza que el grafo refleje verdaderamente la intención arquitectónica del desarrollador y no la implementación técnica de dependencias de terceros.

## Consecuencias
*   **Positivas:** El `ProjectGraph` resultante es mucho más conciso, fiel al código del desarrollador y no contiene "falsos positivos" introducidos por librerías compiladas o dependencias externas.
*   **Negativas:** Obliga a mantener una lista de exclusiones (`ignoreEngine`) actualizada según los ecosistemas soportados, incrementando el coste de mantenimiento cuando se añaden nuevos frameworks.
