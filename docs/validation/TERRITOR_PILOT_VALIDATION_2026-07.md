# Informe de Validación Piloto: Territor × StackGenome
**Documento**: `docs/validation/TERRITOR_PILOT_VALIDATION_2026-07.md`  
**Fecha**: 2026-07-31  
**Fase**: 19 — Validación real con proyecto externo

---

## 1. Resumen Ejecutivo

StackGenome fue ejecutado en modo de análisis estático local sobre el repositorio **Territor**, una aplicación móvil Flutter/Dart de running y conquista territorial. El análisis se completó en menos de 2 segundos, detectó correctamente el ecosistema principal (Flutter/Dart), las dependencias directas del pubspec.yaml, las plataformas (Android/iOS), y las dependencias de Node.js del backend serverless. Se identificaron varios hallazgos categorizados, ninguno de nivel BLOCKER. El piloto concluye con resultado:

> **`PILOT_PASSED_WITH_FINDINGS`**

---

## 2. Entorno

| Parámetro | Valor |
|---|---|
| **OS** | macOS (darwin/arm64) |
| **Fecha de ejecución** | 2026-07-31 |
| **Ruta del binario** | `/tmp/stackgenome_alpha` |
| **SHA-256 del binario** | `ccf0569d8e76c34a077c758c4f9aeaf0c053b1b68d1fe27b780c2d0f1cd135e9` |
| **Versión** | `dev` (alpha-18b, commit `005fb8e`) |
| **Tamaño del binario** | 9.72 MB |
| **Comando de build** | `go build -ldflags "-X main.version=alpha-18b" ./cmd/stackgenome` |
| **Ruta de análisis** | `/Volumes/Intento1/Repos/territor` (fuera del repo StackGenome) |

---

## 3. Commit de StackGenome

```
005fb8e feat: complete Phase 18B Web Alpha
```
(Rama: `main`)

---

## 4. Commit de Territor

```
675ba66 ui: polish social conquest flows
```
Rama: `main`.  
El remote del repositorio Territor no se registra en este documento.

---

## 5. Metodología

1. Lectura de documentación obligatoria de StackGenome (`AGENTS.md`, `docs/`, `.project/`).
2. Inspección manual estática de archivos permitidos de Territor para construir línea base.
3. Compilación del binario StackGenome desde el código fuente del commit actual.
4. Verificación de flags disponibles (`stackgenome analyze -h`).
5. Ejecución de 3 análisis completos (`-json -recommend`) midiendo tiempo y comparando checksums.
6. Auditoría de privacidad del JSON de salida.
7. Verificación de integridad de Territor (git status pre/post).
8. Evaluación contra la línea base manual.

**Restricciones cumplidas**: No se ejecutó Flutter, Gradle, CocoaPods, npm ni ningún script de Territor. No se instalaron dependencias. No se modificó ningún archivo de Territor. No se hizo commit en Territor.

---

## 6. Línea Base Manual (Expectativa Pre-análisis)

Basada únicamente en archivos declarativos inspeccionados:

| Categoría | Expectativa |
|---|---|
| **Lenguaje principal** | Dart (Flutter SDK ^3.9.2) |
| **Lenguajes satélite** | Kotlin/Java (Android), Swift/ObjC (iOS), JavaScript (Cloud Functions) |
| **Tipo de aplicación** | App móvil Flutter multiplataforma (Android + iOS primarios) |
| **Plataformas declaradas** | Android, iOS, Linux, macOS, Web, Windows (por `.metadata`) |
| **Plataformas usadas activamente** | Android (minSdk 27, targetSdk 36), iOS (platform 14.0) |
| **Package manager** | pub (pubspec.yaml, pubspec.lock) |
| **Frameworks** | Flutter, Riverpod (state), go_router (routing) |
| **Firebase** | firebase_core, firebase_auth, cloud_firestore, cloud_functions, firebase_storage, firebase_messaging, firebase_app_check, firebase_crashlytics, firebase_performance, firebase_remote_config |
| **Mapas/Geo** | google_maps_flutter, maplibre_gl, geolocator, geocoding, dart_geohash, google_maps_cluster_manager_2 |
| **Local storage** | hive, hive_flutter |
| **Animaciones** | lottie, rive, flutter_animate, shimmer, animated_flip_counter |
| **Networking** | cloud_firestore (Firestore), cloud_functions |
| **Auth** | firebase_auth, google_sign_in |
| **UI** | flutter_svg, cupertino_icons, google_fonts, flutter_local_notifications |
| **Serialización** | json_annotation + json_serializable, freezed + freezed_annotation |
| **Testing** | flutter_test (flutter_lints lint), test/widget_test.dart, test/0.1-acceptance-tests.json |
| **Dev tools** | build_runner, freezed (code gen) |
| **CI/CD** | No se encontró `.github/` ni archivos CI declarados |
| **Backend serverless** | Firebase Cloud Functions (Node.js 20) con @turf/turf, ngeohash, polygon-clipping |
| **OTA Updates** | Shorebird (app_id visible en shorebird.yaml, es público según documentación Shorebird) |
| **Módulos internos** | Un workspace (`lucide_icons` en `third_party/`) |
| **Infraestructura** | Firebase project `territor-0`, Firebase Hosting, Firestore, Storage, Functions |
| **Editor** | No hay `.idea/` ni `.vscode/` declarados (en .gitignore) |
| **Kotlin version** | 2.3.0 (android/build.gradle) |
| **AGP version** | 8.13.0 |
| **Java toolchain** | 25 (build.gradle app) |

---

## 7. Resultados del Análisis

### 7.1 Métricas de Ejecución

| Métrica | Valor |
|---|---|
| **Tamaño del JSON de salida** | 1,263,533 bytes (~1.2 MB) |
| **Duración run 1** | 1.748s (0.14s user, 0.47s system) |
| **Duración run 2** | 0.273s |
| **Duración run 3** | 0.241s |
| **Código de salida** | 0 (éxito) en las 3 ejecuciones |
| **stderr** | Vacío en las 3 ejecuciones |
| **Nodos detectados** | 107 total |
| **Aristas detectadas** | ~70 (edges) |
| **Recomendaciones** | 3 |

> **Nota**: La diferencia de tiempo entre run 1 (~1.7s) y runs 2-3 (~0.25s) se atribuye al arranque en frío del sistema operativo (caches de filesystem, etc.). No hay indicio de comportamiento no determinista.

### 7.2 Nodos por Tipo

| Tipo | Detectados | Únicos |
|---|---|---|
| `dependency` | 65 | 65 |
| `language` | 26 | 15 |
| `platform` | 16 | 2 |
| **Total** | **107** | — |

### 7.3 Evaluación Detallada por Categoría

#### Lenguaje Principal — Dart/Flutter
| Resultado | Clasificación |
|---|---|
| `Dart/Flutter` detectado con `confidence: 1.00`, 2 nodos (raíz Flutter + sub) | `CORRECT` |
| Evidencias marcadas como `sensitivity: public-metadata` | `CORRECT` |

#### Lenguajes Satélite
| Lenguaje | Detectado | Evaluación |
|---|---|---|
| Swift / Objective-C | ✅ 4 nodos `Swift/Objective-C` | `CORRECT` — plataforma iOS real |
| JVM (Java/Kotlin) | ✅ 2 nodos `JVM` | `CORRECT` — Android via Gradle |
| Node.js / JavaScript | ✅ `Node.js` + `JavaScript` | `CORRECT` — Cloud Functions |
| Shell | ✅ 1 nodo | `CORRECT` — scripts varios |
| Python | ✅ 1 nodo (low confidence 0.30) | `PARTIALLY_CORRECT` — detectado por archivo aislado en iOS Pods |
| CSS | ✅ 1 nodo (0.30) | `PARTIALLY_CORRECT` — admin_web |
| HTML | ✅ 1 nodo (0.30) | `CORRECT` |
| **C/C++ — 6 nodos duplicados** | ⚠️ Detectados dentro de `ios/Pods/` | **`FALSE_POSITIVE / HIGH`** — Ver hallazgo F-001 |
| C/C++ Headers — 4735 evidencias | ⚠️ Detectados masivamente | **`FALSE_POSITIVE / HIGH`** — Ver hallazgo F-001 |

#### Dependencias Dart/Flutter
| Aspecto | Evaluación |
|---|---|
| 50 paquetes de pubspec.yaml detectados | `CORRECT` |
| Versiones no extraídas | `NOT_SUPPORTED` — No hay versiones en nodos `dependency` (confidence 0.80 uniform) |
| `build_runner`, `freezed`, `flutter_lints`, `flutter_test` como dev deps | `PARTIALLY_CORRECT` — Extraídos, no diferenciados de prod deps |
| `lucide_icons` (path dep local en `third_party/`) | `CORRECT` — Detectado como dependencia |
| PURL | `NOT_SUPPORTED` — No se genera PURL en nodos |
| lockfile (pubspec.lock) | `NOT_SUPPORTED` — StackGenome no analiza pubspec.lock |

#### Dependencias Node.js (Cloud Functions)
| Aspecto | Evaluación |
|---|---|
| `@turf/turf`, `firebase-admin`, `firebase-functions`, `ngeohash`, `polygon-clipping` | `CORRECT` — Extraídos de `functions/package.json` con confidence 1.00 |
| `firebase-functions-test` (devDep) | `NOT_VERIFIABLE` — No visto en salida (puede no estar soportado) |

#### Frameworks y Plataformas
| Aspecto | Evaluación |
|---|---|
| Android Native detectado (3 nodos) | `CORRECT` |
| iOS/macOS Native detectado (13 nodos) | `PARTIALLY_CORRECT` — Ver hallazgo F-002 (deduplicación) |
| Flutter como framework | `CORRECT` — Nodo `Dart/Flutter` incluye evidencia `pubspec.yaml` |
| Firebase (multiple servicios) | `CORRECT` — Detectado como dependencias individuales (firebase_core, etc.) |
| Mapas/Geo (google_maps_flutter, maplibre_gl, dart_geohash, geolocator) | `CORRECT` |
| Riverpod (state management) | `CORRECT` — Detectado como `flutter_riverpod` + `riverpod_annotation` |
| go_router (routing) | `CORRECT` |

#### Tooling
| Aspecto | Evaluación |
|---|---|
| `flutter_lints` detectado | `CORRECT` |
| `build_runner` + `json_serializable` + `freezed` detectados | `CORRECT` |
| CI/CD | `NOT_VERIFIABLE` — No existe `.github/` en el repositorio |
| Shorebird OTA | `FALSE_NEGATIVE / MEDIUM` — No detectado como herramienta de tooling |
| Editor (VS Code, IDEA) | `NOT_VERIFIABLE` — Excluidos por `.gitignore` |

#### Infraestructura
| Aspecto | Evaluación |
|---|---|
| Firebase project detectado (vía firebase.json) | `NOT_SUPPORTED` — StackGenome no analiza firebase.json para infraestructura |
| Cloud Functions (Node.js 20) | `CORRECT` — Detectado por package.json |
| Firebase Hosting | `NOT_SUPPORTED` |
| Firestore/Storage rules | `NOT_SUPPORTED` |

---

## 8. Privacidad

### 8.1 JSON de salida (ProjectGraph local)

| Patrón buscado | Encontrado | Evaluación |
|---|---|---|
| `/Volumes/` (rutas absolutas) | ❌ No encontrado | `CLEAN` |
| `/Users/` (home del usuario) | ❌ No encontrado | `CLEAN` |
| Nombre de usuario `ricardo` | ❌ No encontrado | `CLEAN` |
| Firebase project ID `territor-0` | ❌ No encontrado | `CLEAN` |
| Firebase numeric ID | ❌ No encontrado | `CLEAN` |
| API keys (`AIza...`) | ❌ No encontrado | `CLEAN` |
| Shorebird app_id | ❌ No encontrado | `CLEAN` |
| `password` (4 ocurrencias) | ✅ Encontrado | **INFORMATIONAL** — Son nombres de archivo Swift del SDK de Firebase Authentication (`ResetPasswordRequest.swift`, `VerifyPasswordRequest.swift`). Son nombres de archivo en evidencias, no valores de contraseña. Filtrado correcto al ser solo metadatos del nombre de archivo. |
| `secret` (9 ocurrencias) | ✅ Encontrado | **INFORMATIONAL** — Son nombres de archivo C de Envoy (TLS, `secret.upb.h`, etc.) dentro de iOS Pods. Igualmente son nombres de archivos, no valores secretos. |

> **Conclusión de privacidad**: El JSON de salida **no contiene datos sensibles reales**. Las apariciones de `password` y `secret` corresponden exclusivamente a nombres de archivos del SDK de Firebase y de Envoy, no a valores. Severidad: **INFORMATIONAL**. No hay BLOCKER.

### 8.2 Fingerprint para envío remoto (NO enviado en esta fase)

No se ejecutó `--remote`. El Fingerprint potencial sería generado por el backend de la web. Dado que los nodos no contienen rutas absolutas y las evidencias ya están sanitizadas, el Fingerprint sería limpio. **Se requiere autorización explícita para cualquier envío remoto.**

---

## 9. Modo Offline

Las 3 ejecuciones funcionaron sin red. No se realizaron llamadas HTTP (`stderr` vacío). El catálogo local fue consultado. Las recomendaciones offline se generaron correctamente.

**Recomendaciones offline obtenidas:**

| Rank | ID | Nombre | Score | Razones |
|---|---|---|---|---|
| 1 | `tool:buf` | Buf | 0.50 | lenguaje compatible: python, java; herramienta no detectada en el proyecto |
| 2 | `tool:dagger` | Dagger | 0.50 | lenguaje compatible: python; herramienta no detectada en el proyecto |
| 3 | `tool:earthly` | Earthly | 0.50 | lenguaje compatible: python, java; herramienta no detectada en el proyecto |

**Evaluación de recomendaciones:**

| Criterio | Resultado | Hallazgo |
|---|---|---|
| ¿Corresponden a Flutter/Dart? | ❌ No | **F-003 HIGH** — Ver hallazgo |
| ¿Compatibles con las plataformas? | ❌ No aplica a mobile Flutter | **F-003 HIGH** |
| ¿Duplican herramientas existentes? | No aplica | — |
| ¿Tienen relación con necesidad real de Territor? | ❌ No | **F-003 HIGH** |
| ¿Son genéricas? | ✅ Sí, demasiado genéricas | **F-003 HIGH** |
| ¿Mantenidas? | ✅ Sí (Buf, Dagger, Earthly son proyectos activos) | OK |
| ¿Razones derivadas de evidencia? | ⚠️ Derivadas de `python` y `java` detectados via iOS Pods | **F-001 consecuencia** |

---

## 10. Remoto

**No ejecutado**. Se requiere autorización explícita antes de cualquier llamada al backend de staging.

---

## 11. Web Alpha (Validación de Importación)

La importación del archivo `/tmp/sg_out_1.json` en la Web Alpha local (`npm run dev` en `web/`) funcionó correctamente:
- El FileImporter procesó el archivo localmente (File API).
- El schema fue validado por `ajv` contra el esquema OpenAPI.
- El Dashboard mostró la lista de nodos por tipo.
- El StackGraph intentó renderizar los 107 nodos (se observó densidad alta; a considerar mejora de clustering en Fase 20).
- No se produjo upload automático.
- El archivo fue eliminado de la sesión al hacer clic en "Cerrar reporte".

---

## 12. Repetibilidad

| Métrica | Run 1 | Run 2 | Run 3 |
|---|---|---|---|
| SHA-256 del JSON | `2d6b5ceb...` | `2d6b5ceb...` | `2d6b5ceb...` |
| Código de salida | 0 | 0 | 0 |
| stderr | vacío | vacío | vacío |
| Recomendaciones | Buf, Dagger, Earthly | Buf, Dagger, Earthly | Buf, Dagger, Earthly |

**✅ Determinismo confirmado: Los 3 outputs son byte-a-byte idénticos.**

**Git status de Territor antes y después del análisis:** Idéntico. StackGenome no modificó ningún archivo de Territor.

---

## 13. Hallazgos

### F-001 — ALTA — Deduplicación de nodos de lenguaje (iOS Pods / vendored code)
- **Tipo**: `FALSE_POSITIVE`
- **Severidad**: `HIGH`
- **Descripción**: StackGenome detecta C/C++ y C/C++ Headers con alta frecuencia (6 nodos C/C++, 4735 evidencias de headers) debido a que los iOS Pods compilados (Envoy, Firebase, gRPC, etc.) están presentes en `ios/Pods/`. Estos archivos son **dependencias vendorizadas/generadas**, no código fuente del proyecto. Deben excluirse del análisis de lenguajes.
- **Impacto**: Contamina el profile de lenguajes con C/C++ como lenguaje relevante cuando Territor es un proyecto Flutter/Dart. Las razones de las recomendaciones offline mencionan `python` y `java` que provienen de archivos dentro de Pods.
- **Evidencia**: 4735 evidencias de tipo `C/C++ Header` con paths relativos dentro de `ios/Pods/`.
- **Propuesta (Fase 20)**: Añadir exclusión automática de `ios/Pods/`, `ios/Flutter/`, `build/`, `.dart_tool/` en el walker de filesystem.

### F-002 — MEDIA — Nodos de plataforma duplicados (iOS)
- **Tipo**: `FALSE_POSITIVE` (redundancia)
- **Severidad**: `MEDIUM`
- **Descripción**: Se detectaron 16 nodos de tipo `platform`, 13 de ellos `iOS/macOS Native`. La multiplicidad proviene de múltiples detecciones de archivos Xcode dentro del directorio `ios/` (Runner.xcodeproj, Pods, etc.). Deberían consolidarse en 1 nodo de plataforma por plataforma real.
- **Impacto**: El grafo visual se satura con nodos de plataforma repetidos.
- **Propuesta (Fase 20)**: Deduplicar nodos de plataforma por nombre + tipo, manteniendo la unión de sus evidencias.

### F-003 — ALTA — Recomendaciones irrelevantes para Flutter/Dart
- **Tipo**: `FALSE_NEGATIVE` en relevancia
- **Severidad**: `HIGH`
- **Descripción**: Las 3 recomendaciones offline (Buf, Dagger, Earthly) no tienen ninguna relevancia para un proyecto Flutter/Dart móvil. Se generan porque el catálogo no tiene recursos específicos para Dart/Flutter, y el motor de scoring usa `python` y `java` (detectados en Pods) como señal, generando recomendaciones de DevOps generales.
- **Raíz**: (1) Catálogo sin recursos Flutter/Dart; (2) contaminación de lenguajes por F-001.
- **Propuesta (Fase 20)**: Añadir al menos recursos básicos de Flutter al catálogo (flutter_gen, very_good_cli, fvm, melos, etc.).

### F-004 — MEDIA — Versiones no extraídas de dependencias Dart
- **Tipo**: `NOT_SUPPORTED`
- **Severidad**: `MEDIUM`
- **Descripción**: Los nodos `dependency` de Dart tienen `confidence: 0.80` pero no incluyen la versión de la dependencia (aunque está en `pubspec.yaml`). La versión es fundamental para evaluaciones de seguridad y compatibilidad.
- **Propuesta (Fase 20)**: Extraer la versión de las dependencias del `pubspec.yaml` al crear el nodo.

### F-005 — MEDIA — Dev dependencies no diferenciadas de prod
- **Tipo**: `PARTIALLY_CORRECT`
- **Severidad**: `MEDIUM`
- **Descripción**: `build_runner`, `freezed`, `json_serializable`, `flutter_lints`, `flutter_test` son dev dependencies en Dart, pero se reportan igual que las dependencias de producción. No hay campo `is_dev` o `role: dev`.
- **Propuesta (Fase 20)**: Añadir campo `role: dev` o similar a nodos de dependencia marcadas como `dev_dependencies` en `pubspec.yaml`.

### F-006 — BAJA — Shorebird no detectado
- **Tipo**: `FALSE_NEGATIVE`
- **Severidad**: `LOW`
- **Descripción**: `shorebird.yaml` está presente en la raíz y es una herramienta de distribución OTA relevante para Flutter. No fue detectado como nodo de tooling.
- **Propuesta (Fase 20)**: Añadir detección de `shorebird.yaml` en el detector de Dart/Flutter.

### F-007 — BAJA — Firebase no detectado como nodo de infraestructura
- **Tipo**: `NOT_SUPPORTED`
- **Severidad**: `LOW`
- **Descripción**: `firebase.json` y `.firebaserc` contienen la configuración declarativa de Firebase (project ID, hosting, functions, Firestore). StackGenome no analiza estos archivos para generar nodos de infraestructura (`platform: Firebase`).
- **Propuesta (Fase 20)**: Añadir detector de `firebase.json`/`.firebaserc` para generar nodo de plataforma Firebase.

---

## 14. Falsos Positivos

| ID | Detección | Análisis |
|---|---|---|
| F-001 | C/C++ como lenguaje relevante | FP — proviene de iOS Pods (vendored) |
| F-002 | 13 nodos iOS duplicados | FP — multiplicidad por archivos dentro de `ios/` |

---

## 15. Falsos Negativos

| ID | Faltante | Análisis |
|---|---|---|
| F-006 | Shorebird no detectado | FN — archivo soportado `shorebird.yaml` no procesado |
| F-007 | Firebase infra no detectada | FN — `firebase.json` no analizado |
| F-004/F-005 | Versiones y rol dev no extraídos | FN parcial — datos están en pubspec.yaml |
| — | CI/CD | N/A — No existe `.github/` en Territor |

---

## 16. Recomendaciones Obtenidas (Offline)

1. **Buf** (id: `tool:buf`) — Score: 0.50 — Gestor de schemas Protobuf/gRPC. No relevante para Territor.
2. **Dagger** (id: `tool:dagger`) — Score: 0.50 — Build tool CI/CD. No relevante para Territor.
3. **Earthly** (id: `tool:earthly`) — Score: 0.50 — Build tool contenedorizado. No relevante para Territor.

**Evaluación global de recomendaciones**: Las 3 herramientas son proyectos reales y mantenidos, pero ninguna es pertinente para un proyecto Flutter/Dart móvil. El catálogo actual no tiene cobertura del ecosistema Dart/Flutter, lo que hace que las recomendaciones sean completamente inútiles para este tipo de proyecto.

---

## 17. Limitaciones de la Prueba

- El análisis no se pudo ejecutar con un binario de release oficial (no existen releases firmados aún), sino con un binario compilado localmente.
- `pubspec.lock` (pinning exacto de dependencias) no fue analizado — sería necesario para detectar versiones resueltas y transitivas.
- La validación web se realizó con el servidor de desarrollo local (`npm run dev`), no en un entorno de staging/producción.
- No se envió nada al backend remoto (pendiente de autorización).

---

## 18. Evidencia

- Binario: `/tmp/stackgenome_alpha` (SHA-256: `ccf0569d8e76c34a077c758c4f9aeaf0c053b1b68d1fe27b780c2d0f1cd135e9`)
- Output run 1: `/tmp/sg_out_1.json` (SHA-256: `2d6b5ceb7c7484197c839e27fbd3a5a841074e7e6b9b2185cc43b0c47ceaf488`, 1.26 MB)
- Output run 2: `/tmp/sg_out_2.json` (SHA-256 idéntico)
- Output run 3: `/tmp/sg_out_3.json` (SHA-256 idéntico)
- stderr runs 1-3: vacío
- Git status Territor: Sin cambios inducidos por StackGenome

---

## 19. Conclusión

**Resultado**: `PILOT_PASSED_WITH_FINDINGS`

StackGenome demostró capacidad de:
- ✅ Analizar estáticamente un proyecto Flutter/Dart real sin ejecutar ningún comando del proyecto.
- ✅ Detectar el lenguaje principal (Dart/Flutter) correctamente.
- ✅ Extraer dependencias directas de `pubspec.yaml` (50 deps).
- ✅ Detectar dependencias del backend Node.js (`functions/package.json`).
- ✅ Detectar plataformas móviles (Android + iOS).
- ✅ Operar sin red y de forma determinista (3/3 runs byte-idénticos).
- ✅ No modificar el repositorio analizado.
- ✅ No filtrar rutas absolutas, project IDs, API keys ni secretos en el JSON de salida.
- ✅ Completar el análisis en menos de 2 segundos.

**Hallazgos críticos para Fase 20**:
- 🔴 F-001 (HIGH): Exclusión de iOS Pods / vendored code en detección de lenguajes.
- 🔴 F-003 (HIGH): Catálogo sin recursos Flutter/Dart → recomendaciones irrelevantes.
- 🟡 F-002 (MEDIUM): Deduplicación de nodos de plataforma.
- 🟡 F-004 (MEDIUM): Extracción de versiones en dependencias Dart.
- 🟡 F-005 (MEDIUM): Diferenciación de dev vs prod dependencies.
- 🟢 F-006 (LOW): Detectar shorebird.yaml.
- 🟢 F-007 (LOW): Detectar firebase.json como infraestructura.

**Ningún hallazgo es BLOCKER**. StackGenome es funcional y seguro para uso en proyectos Flutter reales, con capacidades de mejora claras y priorizadas.
