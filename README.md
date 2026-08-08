# StackGenome

**StackGenome** es un analizador universal y local-first de proyectos de software. Construye un grafo técnico normalizado del proyecto y lo utiliza para recomendar herramientas, paquetes, repositorios, plugins, frameworks e implementaciones compatibles, mantenidos y seguros.

## Estado

- Nombre aprobado: **StackGenome**
- Etapa actual: **Public Alpha / Pre-Release (Fase 21)**
- Producto inicial: **CLI multiplataforma**
- Lenguaje del CLI: **Go**
- Backend: **Cloudflare Workers + D1**
- Privacidad predeterminada: **análisis local y envío exclusivo de metadatos autorizados (--remote opt-in)**

## 🚀 Uso Rápido (Alpha)

Puedes descargar los binarios pre-compilados (para macOS, Linux o Windows) desde los Releases.

```bash
# Analiza el directorio actual
stackgenome analyze .

# Analiza y obtén recomendaciones de herramientas usando el catálogo offline
stackgenome analyze --recommend .

# Analiza y consulta la nube (Cloudflare) para recomendaciones actualizadas (con anonimización forzada)
stackgenome analyze --remote .
```

## 🔒 Privacidad Extrema

StackGenome se construyó con la premisa de "Zero-Telemetry".
- El análisis local y el catálogo embebido funcionan sin conexión. Las recomendaciones remotas son opcionales y requieren conexión explícita.
- No extrae código fuente; solo lee manifiestos (ej: `package.json`, `go.mod`, etc.).
- Cuando usas la bandera `--remote`, el sistema envía a Cloudflare únicamente la topología básica de tu stack. El Fingerprint utiliza una política *metadata-only* y las pruebas realizadas no detectaron exposición de los campos sensibles prohibidos evaluados.

## ⚠️ Limitaciones Conocidas (Alpha)

Al ser un release Alpha, existen ciertas limitaciones temporales:
1. **Catálogo remoto limitado**: El backend en Cloudflare (D1) contiene actualmente una lista curada de herramientas. Se ampliará.
2. **Rate Limits Remotos**: Si usas `--remote`, estás sujeto a los límites del plan Free de Cloudflare Workers (aprox 100K requests por día globales).
3. **Lenguajes Soportados**: StackGenome soporta en profundidad Go, Node.js (JS/TS), Python, Rust, Java/JVM, .NET, Swift, PHP, Ruby y C/C++. Otros ecosistemas podrían ser identificados genéricamente, pero sin extracción profunda de paquetes.
4. **Symlink Depth**: El `SafeWalker` está configurado por seguridad para detenerse tras ciertos niveles de recursión o si detecta bucles.

## Comunidad

- [Guía de Contribución](CONTRIBUTING.md)
- [Política de Seguridad](SECURITY.md)
- Licencia: [MIT](LICENSE)

## Ruta local prevista

```text
/Volumes/intento1/Repos/StackGenome
```

La partición debe estar montada antes de instalar herramientas o trabajar en el proyecto. Los scripts se niegan a continuar si `/Volumes/intento1` no está disponible, evitando que cachés o SDK se escriban accidentalmente en la partición principal.

## Inicio rápido para agentes

1. Leer `AGENTS.md`.
2. Leer `docs/00_INDEX.md`.
3. Leer `.project/CURRENT_PHASE.md`.
4. No implementar fases posteriores.
5. Ejecutar únicamente el trabajo autorizado.
6. Probar, documentar y detenerse para solicitar autorización.

El primer prompt se encuentra en [`START_PROMPT.md`](START_PROMPT.md).

## Inicio rápido para el entorno

```bash
cd /Volumes/intento1/Repos/StackGenome
source scripts/activate-env.sh
./scripts/bootstrap-macos.sh
./scripts/verify-environment.sh
```

El script de bootstrap no modifica automáticamente `~/.zshrc`, no instala Homebrew y no mueve configuraciones globales.

## Documentación

Consulta [`docs/00_INDEX.md`](docs/00_INDEX.md) para el orden de lectura y la autoridad de cada documento.
