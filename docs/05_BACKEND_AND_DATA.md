# Backend, catálogo y datos

## 1. Momento de implementación

El backend no se construye hasta que:

- ProjectGraph y Fingerprint estén estables;
- el CLI produzca recomendaciones con catálogo local;
- existan fixtures y evaluación básica;
- la fase sea autorizada.

## 2. Stack previsto

- Cloudflare Workers;
- TypeScript;
- D1;
- R2 posterior;
- Pages/Workers para web posterior;
- API REST versionada.

El plan gratuito es apropiado para prototipo si el análisis pesado permanece local. Los límites se deben verificar nuevamente al iniciar la fase de backend.

## 3. Responsabilidades

El backend:

- valida fingerprint;
- consulta catálogo;
- aplica reglas;
- calcula ranking;
- devuelve recursos y razones;
- recibe submissions en una fase posterior.

No:

- recibe repositorios;
- compila código;
- ejecuta archivos del usuario;
- almacena código;
- usa un LLM para ranking;
- promete auditoría de seguridad.

## 4. API v1 conceptual

```text
GET  /v1/health
GET  /v1/schema/fingerprint
POST /v1/recommendations
GET  /v1/resources/{id}
GET  /v1/search
POST /v1/submissions          posterior
POST /v1/feedback             posterior
```

### `POST /v1/recommendations`

Request:

- schema version;
- fingerprint;
- goal;
- policy;
- limit.

Response:

- catalog snapshot id;
- ranking version;
- recommendations;
- reasons;
- exclusions resumidas;
- warnings;
- expiration.

## 5. Modelo de datos

### resources

- id;
- type;
- canonical_name;
- summary;
- canonical_url;
- source;
- source_id;
- ecosystem;
- license;
- status;
- verified_at;
- updated_at.

### resource_versions

- resource_id;
- version;
- published_at;
- yanked/deprecated;
- compatibility metadata.

### technologies

- id;
- kind;
- name;
- aliases.

### resource_technologies

- resource_id;
- technology_id;
- relation;
- version_constraint.

### source_metrics

- stars;
- forks;
- downloads;
- likes;
- popularity;
- last_release;
- last_commit;
- fetched_at.

### security_signals

- advisory source;
- severity;
- affected range;
- scorecard signals;
- fetched_at.

### compatibility_rules

- subject;
- predicate;
- object;
- constraint;
- confidence;
- source;
- valid_from/to.

### ranking_snapshots

- algorithm version;
- weights;
- catalog version.

## 6. Identificadores

- PURL para paquetes;
- URL canónica para repositorios;
- ids internos opacos;
- aliases normalizados;
- hashes solo para integridad, no identidad personal.

## 7. Ingesta

Prioridad:

1. APIs y dumps oficiales;
2. ecosyste.ms/deps.dev u otras fuentes abiertas;
3. registries;
4. GitHub/GitLab bajo límites;
5. submissions verificadas.

No hacer crawling indiscriminado.

## 8. Actualización

- recursos populares: frecuente;
- long tail: bajo demanda;
- seguridad: prioridad alta;
- cache condicional;
- provenance por campo;
- no sobrescribir un dato de mayor autoridad con uno inferior.

## 9. Búsqueda

Primera versión:

- D1 + FTS5;
- coincidencia exacta y aliases;
- filtros SQL;
- ranking técnico separado del ranking textual.

Vectores solo si una necesidad medida lo justifica.

## 10. Retención

- fingerprints no se persisten por defecto;
- logs sin payload completo;
- feedback separado;
- política pública antes de beta;
- herramientas para eliminación cuando haya cuentas.

## 11. Free tier

Diseñar para:

- pocas consultas D1 por request;
- statements preparados;
- caché de recursos;
- payloads pequeños;
- sin procesamiento de repositorios en Worker.

Los límites actuales son una condición operativa, no una garantía contractual.
