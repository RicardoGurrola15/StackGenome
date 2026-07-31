# Protocolo de ejecución para agentes

## 1. Objetivo

Permitir que Codex, Antigravity u otro agente continúen el proyecto sin perder contexto ni ejecutar etapas no autorizadas.

## 2. Ciclo obligatorio

### A. Orientar

- leer reglas;
- leer fase;
- revisar Git;
- buscar memoria;
- identificar restricciones;
- no modificar.

### B. Planear

Presentar:

- objetivo;
- tareas;
- archivos;
- comandos;
- dependencias;
- riesgos;
- criterios de aceptación.

Si una instalación toca configuración global, detenerse para autorización.

### C. Implementar

- cambios pequeños;
- pruebas junto al código;
- no refactors ajenos;
- no fase siguiente.

### D. Verificar

- format;
- vet/lint;
- tests;
- race cuando corresponda;
- fixtures;
- inspección de diff;
- comprobación de privacidad.

### E. Documentar

Actualizar:

- estado;
- fase;
- handoff;
- decisiones;
- ADR;
- documentación técnica.

### F. Recordar

Si Engram está disponible:

- guardar decisión;
- guardar descubrimientos;
- guardar error/solución;
- guardar siguiente paso.

No sustituir archivos por memoria.

### G. Entregar y parar

Usar el formato definido en `AGENTS.md` y pedir autorización.

## 3. Clasificación de cambios

### Pequeño

- documentación;
- test aislado;
- bug local;
- sin cambio de contrato.

### Medio

- nueva regla;
- adaptador;
- comando;
- migración compatible.

### Grande

- cambio de schema;
- dependencia central;
- seguridad;
- backend;
- plugin model.

Cambios grandes requieren ADR antes de implementación.

## 4. Manejo de incertidumbre

El agente debe:

1. revisar documentación y código;
2. buscar evidencia;
3. formular la mínima suposición;
4. marcarla;
5. no inventar API ni comportamiento externo;
6. detenerse si afecta arquitectura, seguridad o datos.

## 5. Política de dependencias

Antes de añadir:

- problema;
- alternativa estándar;
- licencia;
- mantenimiento;
- tamaño;
- riesgo;
- versión;
- plan de eliminación.

Registrar en decisión o ADR.

## 6. Prohibiciones

- no continuar automáticamente;
- no afirmar que pasó una prueba no ejecutada;
- no ocultar warnings;
- no cambiar rutas del usuario;
- no instalar globalmente por comodidad;
- no almacenar secretos en Engram;
- no crear un LLM escondido en el pipeline;
- no usar repositorios externos como fixtures sin licencia y minimización.

## 7. Handoff

Debe permitir que otro agente responda:

- qué existe;
- qué funciona;
- qué falló;
- qué se decidió;
- qué sigue;
- qué no debe tocar.
