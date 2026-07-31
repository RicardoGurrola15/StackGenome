# ProjectGraph Specification v0

## 1. Objetivo

`ProjectGraph` es el modelo interno canónico. Representa la composición técnica del proyecto sin conservar código fuente.

`Fingerprint` es una proyección sanitizada y más pequeña del ProjectGraph, apta para vista previa y eventual envío.

## 2. Identidad

```json
{
  "schema": "dev.stackgenome.projectgraph",
  "schema_version": "0.1.0",
  "scanner_version": "0.0.0-dev",
  "scan_id": "local-random-id",
  "created_at": "RFC3339"
}
```

`scan_id` no debe ser estable entre escaneos salvo opción local explícita. No contiene identificadores de máquina.

## 3. Entidades

### Project

- tipo predominante;
- módulos;
- repositorios;
- capacidades;
- objetivos opcionales;
- warnings;
- cobertura de análisis.

### Repository

- raíz relativa;
- VCS;
- branch opcional local;
- remote excluido del fingerprint predeterminado;
- submodules declarados;
- estado de workspace opcional y local.

### Module

- id estable derivado de ruta relativa y tipo;
- nombre sanitizado;
- raíz relativa;
- module type;
- outputs;
- ecosystems;
- relationships.

### Language

- nombre canónico;
- bytes;
- files;
- role;
- confidence;
- evidence refs.

### Ecosystem / PackageManager

- ecosystem;
- manager;
- manifests;
- lockfiles;
- version declarada;
- confidence.

### Component

Dependencia o recurso técnico:

- PURL cuando aplique;
- name;
- version;
- scope;
- direct/transitive/unknown;
- development/runtime/build;
- source evidence.

### Framework

- name;
- version/range;
- role;
- module;
- evidence.

### Platform

- web;
- backend;
- android;
- ios;
- macos;
- windows;
- linux;
- embedded;
- wasm;
- cloud;
- other.

### Service

- database;
- queue;
- cache;
- object storage;
- analytics;
- auth;
- AI provider;
- external API.

No se deben extraer endpoints, tenant ids o nombres de proyectos por defecto.

### Tooling

- build;
- test;
- lint;
- format;
- CI/CD;
- editor declarations;
- container;
- infrastructure-as-code.

### EnvironmentSnapshot

Separado y opt-in:

- OS y arquitectura;
- versiones de runtimes;
- compiladores;
- editor;
- extensiones permitidas.

Nunca contiene usuario, hostname, IP, serial, rutas de home o variables completas.

### Evidence

- id;
- kind;
- relative path;
- selector o línea aproximada;
- valor sanitizado;
- sensitivity;
- detector;
- confidence contribution.

### Relationship

Ejemplos:

- `contains`;
- `depends_on`;
- `builds`;
- `targets`;
- `uses`;
- `generated_by`;
- `configured_by`;
- `compatible_with`;
- `conflicts_with`.

## 4. Rutas

- Siempre relativas a la raíz autorizada.
- Separador normalizado `/`.
- No incluir `/Users/<name>`.
- No seguir symlinks fuera de raíz.
- Los nombres de carpetas potencialmente sensibles pueden pseudonimizarse en Fingerprint.

## 5. Fingerprint

Incluye únicamente allowlist:

- tipos de módulos;
- lenguajes y proporciones;
- ecosystems;
- nombres y versiones de dependencias públicas;
- frameworks;
- plataformas;
- herramientas declaradas;
- categorías de servicios;
- capacidades;
- cobertura y warnings no sensibles;
- objetivos de recomendación elegidos.

Excluye:

- código;
- README completo;
- remote URL;
- branch privado;
- rutas absolutas;
- nombre del proyecto por defecto;
- nombres internos de módulos cuando no sean necesarios;
- variables;
- endpoints;
- secretos;
- comentarios;
- historial Git.

## 6. Versionado

- SemVer para schema.
- Cambios aditivos compatibles: minor.
- Eliminaciones o cambio semántico: major.
- Cada fingerprint declara versión.
- Backend rechaza versiones incompatibles con error accionable.

## 7. Exportaciones

- JSON canónico de StackGenome.
- CycloneDX opcional para componentes y relaciones.
- Resumen terminal.
- No se promete equivalencia completa entre ProjectGraph y SBOM.

## 8. Validación

- JSON Schema generado o mantenido junto al modelo.
- invariantes de grafo;
- ids únicos;
- referencias válidas;
- rutas relativas;
- ausencia de campos prohibidos;
- orden estable para serialización y snapshots.
