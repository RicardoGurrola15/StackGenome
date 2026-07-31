# ADR-0007: Elección de Next.js 14 para la Web Alpha

- **Estado**: aceptado
- **Fecha**: 2026-07-31
- **Fase**: Fase 18B

## Contexto
El CLI de StackGenome genera reportes locales en JSON y requiere una Web Alpha que sirva como landing page, catálogo público y visor de estos reportes. 
Las opciones evaluadas fueron Astro (excelente SSG), SvelteKit y Next.js.

## Decisión
Se eligió **Next.js 14 (App Router)** por los siguientes motivos:
1. Combinación híbrida nativa: SSG estático para el catálogo y la landing, CSR (React del lado del cliente) para el visualizador del grafo (importador de JSON).
2. Misma tecnología base (TypeScript + Node) que se usa y se usará en el backend de Cloudflare Workers.
3. Tipado sin costuras: Se reutilizan las mismas interfaces DTO (`ProjectGraphDTO`) del backend y del CLI.
4. Compatibilidad con el ecosistema de Cloudflare Pages (vía `@cloudflare/next-on-pages` a futuro).

## Consecuencias
- La aplicación web depende de React y del framework Next.js.
- Se debe tener especial cuidado con el tamaño del bundle, evitando librerías pesadas donde no sean indispensables.
- La web reside en el directorio `/web/` dentro del repositorio principal (monorepo).
