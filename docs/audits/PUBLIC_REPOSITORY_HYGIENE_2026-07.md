# Informe de Higiene del Repositorio Público
**Documento**: `docs/audits/PUBLIC_REPOSITORY_HYGIENE_2026-07.md`  
**Versión**: 2.0 (actualizado Fase 20)  
**Fecha**: 2026-08-07  
**Auditor**: Agente StackGenome (Fase 20A)

---

## 1. Inventario Git

### 1.1 Remotes
```
origin  https://github.com/RicardoGurrola15/StackGenome.git (fetch/push)
```

### 1.2 Commits rastreados (historial completo, 6 commits)
```
cf18679 validation: Phase 19 Territor pilot — PILOT_PASSED_WITH_FINDINGS
005fb8e feat: complete Phase 18B Web Alpha
373b766 chore: exclude backend/node_modules and .wrangler from git tracking
54ea6e3 feat(18a): web readiness gate — CORS hardening, OpenAPI contract, catalog validator
59723ba ci: automate release and deployment pipelines (Phase 17)
abc3f34 chore: initial commit (Alpha Release)
```

### 1.3 Tags
Ninguno (sin releases publicados).

### 1.4 Ramas
Solo `main` (local y remoto).

---

## 2. Secretos y datos sensibles

### 2.1 Resultado del escaneo manual del historial

| Patrón | Historial | Working Tree | Evaluación |
|--------|-----------|--------------|------------|
| Claves API (`AIza`, `sk-`, `Bearer`) | ❌ No encontrado | ❌ No encontrado | CLEAN |
| `.env` archivos | ❌ No encontrado | ❌ No encontrado | CLEAN |
| `wrangler.toml` con secrets | ❌ No encontrado | ❌ No encontrado | CLEAN |
| Firebase project IDs | ❌ No encontrado | ❌ No encontrado | CLEAN |
| Rutas absolutas `/Users/ricardo` | ❌ No encontrado | ❌ No encontrado | CLEAN |
| Shorebird `app_id` | ❌ No encontrado | ❌ No encontrado | CLEAN |
| Tokens CI/CD | ❌ No encontrado | ❌ No encontrado | CLEAN |
| Cloudflare Account ID | ❌ No encontrado en código | Variables de entorno en CI | INFORMATIONAL |
| Credenciales BD | ❌ No encontrado | ❌ No encontrado | CLEAN |

> **Conclusión**: No se encontraron secretos reales en el historial ni en el working tree. No hay BLOCKER. El `CLOUDFLARE_ACCOUNT_ID` referenciado en CI se lee de GitHub Secrets (no está en código).

---

## 3. Archivos que deben permanecer públicos

| Archivo/Directorio | Motivo |
|--------------------|--------|
| `cmd/`, `internal/`, `pkg/` | Código fuente del CLI |
| `backend/src/`, `backend/migrations/` | Código del Worker y migraciones D1 |
| `web/app/`, `web/components/`, `web/lib/` | Frontend Next.js |
| `docs/` | Documentación técnica y de producto |
| `.project/` | Estado y continuidad del proyecto |
| `go.mod`, `go.sum` | Dependencias Go |
| `web/package.json`, `web/package-lock.json` | Dependencias Node |
| `.github/workflows/` | CI/CD |
| `.gitignore` | Exclusiones |
| `README.md` | Presentación pública |
| `internal/catalog/catalog.json` | Catálogo local de recursos |

---

## 4. Archivos que deben ser locales (excluidos)

| Archivo/Directorio | Razón | Estado .gitignore |
|--------------------|-------|--------------------|
| `.engram/` | Memoria local auxiliar con posibles observaciones privadas | ✅ Cubierto (`.engram/`) |
| `stackgenome` | Binario Go compilado | ✅ Cubierto |
| `stackgenome_cli` | Binario compilado localmente (Fase 19) | ✅ Añadido en Fase 20A |
| `fix_catalog.py` | Script temporal de trabajo | ✅ Añadido en Fase 20A |
| `territor-*.json` | Outputs de análisis de proyectos externos | ✅ Cubierto |
| `sg_out_*.json` | Outputs del CLI | ✅ Cubierto |
| `.wrangler/` | Estado local de Wrangler/Cloudflare | ✅ Cubierto |
| `.dev.vars` | Variables de entorno locales Cloudflare | ✅ Cubierto |
| `backend/node_modules/` | Dependencias Node del backend | ✅ Cubierto |
| `web/node_modules/` | Dependencias Node del frontend | ✅ Cubierto |
| `web/.next/` | Build de desarrollo Next.js | ✅ Cubierto |
| `web/out/` | Build de producción Next.js | ✅ Cubierto |
| `.DS_Store` | Metadatos macOS | ✅ Cubierto |
| `coverage.out`, `*.test`, `*.prof` | Artefactos de tests | ✅ Cubierto |

---

## 5. Cambios aplicados al .gitignore en Fase 20A

```diff
+ stackgenome_cli        # Binario compilado localmente (no release oficial)
+ fix_catalog.py        # Script temporal de trabajo
```

El resto del `.gitignore` ya estaba en buen estado desde Fase 20 previa.

---

## 6. .engram — Auditoría y decisión

### Contenido encontrado
- `.engram/README.md` (304 bytes) — Solo contiene texto explicativo genérico sobre Engram. No contiene observaciones privadas, rutas de usuario, datos de proyectos externos ni secretos.

### Acción tomada
- `.engram/README.md` fue staged para eliminación del índice Git (`D  .engram/README.md`) en una sesión anterior
- `.engram/` completo está ahora cubierto por `.gitignore`
- El contenido local se preserva (no borrado del filesystem)
- Ninguna información técnica relevante fue migrada ya que el README era solo explicativo

### Estado posterior
- `.engram/` no se publica en el repositorio público
- La documentación viva del proyecto continúa en `docs/`, ADR y `.project/`

---

## 7. Riesgos históricos

| Riesgo | Nivel | Detalle |
|--------|-------|---------|
| Historial con `.engram/README.md` en commits previos | LOW | El archivo solo contenía texto genérico sin datos privados. No requiere reescritura de historial. |
| Cloudflare Account ID visible en CI si los logs son públicos | INFORMATIONAL | Referenciado solo como `${{ secrets.CLOUDFLARE_ACCOUNT_ID }}`. No en código fuente. |

---

## 8. Acciones no ejecutadas (requieren autorización adicional)

Ninguna. No se encontró ningún BLOCKER ni acción que requiera autorización adicional.

---

## 9. Resultado

**HYGIENE_PASSED** — El repositorio es seguro para publicación pública. No hay secretos, rutas absolutas de usuario, IDs de proyectos externos ni datos privados en el historial o working tree relevante.
