# Product Charter — StackGenome

## 1. Visión

StackGenome será una herramienta universal que inspecciona proyectos de software localmente, construye un inventario técnico verificable y recomienda recursos reales que encajan con su stack, versiones, plataformas, objetivos y restricciones.

No es otro buscador general de repositorios. Su valor es la combinación de:

- comprensión estructurada del proyecto;
- privacidad local-first;
- catálogo federado;
- compatibilidad verificable;
- señales de seguridad y mantenimiento;
- recomendaciones pequeñas, explicables y accionables.

## 2. Usuarios

- desarrolladores que inician un proyecto;
- mantenedores de proyectos existentes;
- equipos que exploran herramientas;
- estudiantes;
- arquitectos de software;
- desarrolladores nativos, web, backend, móvil, escritorio, CLI, librerías, videojuegos, firmware e infraestructura;
- mantenedores que desean registrar sus recursos en el catálogo futuro.

## 3. Problemas

1. Los recursos útiles están dispersos entre registries, forjas, blogs y catálogos.
2. Popularidad no equivale a compatibilidad, mantenimiento o seguridad.
3. Buscar manualmente exige conocer de antemano la terminología correcta.
4. Un proyecto puede ser monorepo y combinar múltiples lenguajes, módulos y plataformas.
5. Muchos usuarios no quieren enviar su código a servicios externos.
6. Los buscadores devuelven demasiados resultados, sin priorizar los que realmente encajan.

## 4. Propuesta de valor inicial

El CLI:

1. descubre repositorios y módulos;
2. identifica lenguajes principales y satélite;
3. identifica manifests, lockfiles, frameworks, herramientas, infraestructura y editores declarados;
4. genera un `ProjectGraph`;
5. muestra al usuario el fingerprint enviable;
6. obtiene tres recomendaciones principales;
7. explica cada recomendación mediante datos estructurados, no mediante un LLM.

## 5. Principios

### Universal por arquitectura

El núcleo no conoce un único lenguaje. Los ecosistemas se añaden mediante adaptadores.

### Progresivo por profundidad

StackGenome puede reconocer muchos proyectos, pero declara honestamente el nivel de soporte de cada ecosistema:

- full;
- substantial;
- basic;
- detected-only;
- unsupported.

### Local-first

La detección y normalización ocurren localmente. El servidor recibe únicamente metadatos autorizados.

### Determinista

Las mismas entradas, catálogo y configuración deben producir el mismo orden.

### Evidence-first

Toda afirmación debe apuntar a evidencia: archivo, regla, manifest, comando local autorizado o dato de catálogo.

### Seguridad antes que popularidad

Compatibilidad, vulnerabilidades, mantenimiento y licencia pesan más que estrellas.

### Extensible sin ejecutar terceros

El MVP usa adaptadores compilados. Un marketplace de plugins ejecutables queda fuera hasta diseñar sandbox y firma.

## 6. Alcance del MVP CLI

Incluye:

- macOS, Linux y Windows;
- repositorios simples y monorepos;
- detección de lenguajes;
- manifests y lockfiles;
- módulos;
- frameworks y plataformas reconocibles;
- archivos de CI/CD e infraestructura;
- editores declarados en el repositorio;
- escaneo opcional del entorno;
- Project Graph y fingerprint;
- catálogo local semilla;
- recomendaciones deterministas;
- salida terminal y JSON;
- privacidad y redacción de datos.

## 7. Fuera del MVP

- web pública;
- app móvil o escritorio;
- login, likes, comentarios y guardados;
- crawling masivo;
- análisis semántico exhaustivo de todo el código;
- ejecución de código del repositorio inspeccionado;
- instalación automática de recomendaciones;
- LLM dentro del producto;
- plugins binarios descargados;
- evaluación legal definitiva de licencias;
- garantía de ausencia de vulnerabilidades.

## 8. Métricas de éxito del alpha

- analiza proyectos pequeños en menos de 2 segundos y repositorios medianos en menos de 10 segundos en hardware razonable;
- no sale de la raíz autorizada;
- no lee archivos ignorados o sensibles sin regla explícita;
- identifica correctamente módulos y lenguajes en el conjunto de fixtures;
- produce resultados estables;
- cero secretos en fingerprints de pruebas;
- recomendaciones relevantes evaluadas manualmente en una muestra diversa;
- binarios reproducibles para macOS, Linux y Windows.

## 9. Evolución

1. CLI local.
2. API y catálogo remoto.
3. web pública tipo buscador y reporte.
4. cuentas y comunidad.
5. aplicación cliente.
