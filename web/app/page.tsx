import type { Metadata } from 'next'
import Link from 'next/link'
import styles from './page.module.css'

export const metadata: Metadata = {
  title: 'StackGenome — Analiza el ADN de tu stack tecnológico',
}

const ECOSYSTEMS = [
  { icon: '🐹', name: 'Go', desc: 'Módulos, workspaces, frameworks' },
  { icon: '🐍', name: 'Python', desc: 'pip, Poetry, pyproject.toml' },
  { icon: '🟨', name: 'JavaScript / TypeScript', desc: 'npm, pnpm, yarn, workspaces' },
  { icon: '🦀', name: 'Rust', desc: 'Cargo, workspaces' },
  { icon: '🎯', name: 'Dart / Flutter', desc: 'pubspec.yaml, Flutter SDK' },
  { icon: '☕', name: 'JVM', desc: 'Maven, Gradle (detección)' },
  { icon: '🔷', name: '.NET', desc: '.csproj, .sln (detección)' },
  { icon: '🐘', name: 'PHP / Ruby', desc: 'Composer, Bundler (detección)' },
  { icon: '⚙️', name: 'C / C++', desc: 'CMake, Conan (detección)' },
  { icon: '🍎', name: 'Swift', desc: 'SPM, CocoaPods (detección)' },
]

const FEATURES = [
  {
    icon: '🔒',
    title: 'Local-first, siempre',
    desc: 'El análisis ocurre en tu máquina. Tu código nunca sale de tu entorno sin tu consentimiento explícito.',
  },
  {
    icon: '🧬',
    title: 'Stack Graph completo',
    desc: 'Visualiza lenguajes, frameworks, package managers, dependencias e infraestructura como un grafo conectado.',
  },
  {
    icon: '📦',
    title: 'Un solo binario',
    desc: 'Sin Node, sin Docker, sin runtime. Descarga y ejecuta en macOS, Linux o Windows.',
  },
  {
    icon: '🎯',
    title: 'Recomendaciones objetivas',
    desc: 'Herramientas sugeridas basadas en tu stack real. Sin LLMs, sin guesses. Deterministas y trazables.',
  },
  {
    icon: '🔍',
    title: 'Más de 10 ecosistemas',
    desc: 'Detección profunda para Go, Python, JS/TS y Rust. Detección de presencia para 6 ecosistemas más.',
  },
  {
    icon: '🛡️',
    title: 'Sin secretos expuestos',
    desc: 'El sanitizer filtra .env, rutas absolutas y evidencias sensibles antes de cualquier serialización.',
  },
]

export default function HomePage() {
  return (
    <>
      {/* ── Hero ── */}
      <section className={styles.hero} aria-label="Introducción">
        <div className="container">
          <div className={styles.heroBadge}>
            <span className="badge badge--orange">Alpha</span>
            <span style={{ color: 'var(--text-muted)', fontSize: 'var(--text-sm)' }}>
              Herramienta de análisis local-first
            </span>
          </div>
          <h1 className={styles.heroTitle}>
            El ADN de tu{' '}
            <span className={styles.heroAccent}>stack tecnológico</span>
          </h1>
          <p className={styles.heroDesc}>
            StackGenome analiza tu repositorio localmente y genera un mapa completo de lenguajes,
            frameworks, dependencias e infraestructura. Sin subir tu código. Sin cuentas.
            Sin sorpresas.
          </p>
          <div className={styles.heroCta}>
            <Link href="/install" className="btn btn--primary btn--lg">
              Descargar CLI ↓
            </Link>
            <Link href="/report" className="btn btn--secondary btn--lg">
              Ver demo de reporte →
            </Link>
          </div>
          <div className={styles.heroCode} aria-label="Ejemplo de uso del CLI">
            <pre>
              <code>
                <span style={{ color: 'var(--text-muted)' }}>$</span>{' '}
                <span style={{ color: 'var(--accent-green)' }}>stackgenome</span>{' '}
                <span style={{ color: 'var(--accent-blue)' }}>analyze</span>{' '}
                <span style={{ color: 'var(--accent-orange)' }}>--json</span>{' '}{' '}
                <span style={{ color: 'var(--text-muted)' }}>/mi-proyecto</span>
                {'\n'}
                <span style={{ color: 'var(--text-muted)', fontSize: 'var(--text-xs)' }}>
                  {'{'}&quot;schema_version&quot;: &quot;1.0.0&quot;, &quot;nodes&quot;: [...], ...{'}'}
                </span>
              </code>
            </pre>
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section className="section" aria-label="Características">
        <div className="container">
          <h2 style={{ textAlign: 'center', marginBottom: 'var(--space-3)' }}>¿Por qué StackGenome?</h2>
          <p style={{ textAlign: 'center', marginBottom: 'var(--space-12)', maxWidth: 600, margin: '0 auto var(--space-12)' }}>
            Diseñado para desarrolladores que quieren entender su stack sin comprometer su privacidad.
          </p>
          <div className="grid-auto">
            {FEATURES.map(f => (
              <div key={f.title} className="card">
                <div style={{ fontSize: '1.75rem', marginBottom: 'var(--space-3)' }}>{f.icon}</div>
                <h3 style={{ marginBottom: 'var(--space-2)', fontSize: 'var(--text-base)' }}>{f.title}</h3>
                <p style={{ fontSize: 'var(--text-sm)' }}>{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Ecosystems ── */}
      <section className="section" style={{ background: 'var(--bg-secondary)' }} aria-label="Ecosistemas soportados">
        <div className="container">
          <h2 style={{ textAlign: 'center', marginBottom: 'var(--space-3)' }}>Ecosistemas soportados</h2>
          <p style={{ textAlign: 'center', marginBottom: 'var(--space-12)', maxWidth: 500, margin: '0 auto var(--space-12)' }}>
            Detección profunda para los más usados. Detección de presencia para el resto.
          </p>
          <div className="grid-auto">
            {ECOSYSTEMS.map(e => (
              <div key={e.name} className="card" style={{ display: 'flex', gap: 'var(--space-3)', alignItems: 'flex-start' }}>
                <span style={{ fontSize: '1.5rem', flexShrink: 0 }}>{e.icon}</span>
                <div>
                  <div style={{ fontWeight: 600, fontSize: 'var(--text-sm)', marginBottom: 'var(--space-1)' }}>{e.name}</div>
                  <div style={{ fontSize: 'var(--text-xs)', color: 'var(--text-muted)', fontFamily: 'var(--font-mono)' }}>{e.desc}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Privacy callout ── */}
      <section className="section" aria-label="Política de privacidad">
        <div className="container" style={{ maxWidth: 800 }}>
          <div className="card" style={{ borderColor: 'var(--accent-green-dim)', textAlign: 'center' }}>
            <div style={{ fontSize: '2.5rem', marginBottom: 'var(--space-4)' }}>🔒</div>
            <h2 style={{ marginBottom: 'var(--space-4)' }}>Tu código nunca sale de tu máquina</h2>
            <p style={{ marginBottom: 'var(--space-4)' }}>
              StackGenome opera localmente. El análisis ocurre en tu proceso, en tu disco.
              Si usas la consulta remota opcional, solo se envían metadatos anónimos —tipos de nodos y relaciones —
              previa revisión y confirmación tuya. Nunca se envía código fuente, rutas absolutas ni secretos.
            </p>
            <Link href="/report" className="btn btn--primary">
              Ver cómo funciona la importación local →
            </Link>
          </div>
        </div>
      </section>

      {/* ── Limitations ── */}
      <section className="section" style={{ background: 'var(--bg-secondary)' }} aria-label="Limitaciones conocidas">
        <div className="container" style={{ maxWidth: 700 }}>
          <h2 style={{ marginBottom: 'var(--space-6)' }}>Limitaciones conocidas del Alpha</h2>
          <ul style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-3)', listStyle: 'none' }}>
            {[
              'Los lockfiles (yarn.lock, Cargo.lock, etc.) no son analizados aún.',
              'JVM, .NET, Swift, PHP y Ruby tienen detección básica — sin extracción de dependencias.',
              'El catálogo tiene ~30 herramientas. La expansión masiva es posterior.',
              'Sin modo offline para el catálogo web — requiere conexión al backend de staging.',
              'Sin autenticación, sin cuentas, sin historial remoto (por diseño en Alpha).',
            ].map(l => (
              <li key={l} style={{ display: 'flex', gap: 'var(--space-3)', alignItems: 'flex-start', fontSize: 'var(--text-sm)', color: 'var(--text-secondary)' }}>
                <span style={{ color: 'var(--accent-orange)', flexShrink: 0 }}>⚠</span>
                {l}
              </li>
            ))}
          </ul>
        </div>
      </section>

      {/* ── CTA final ── */}
      <section className="section" aria-label="Comenzar">
        <div className="container" style={{ textAlign: 'center' }}>
          <h2 style={{ marginBottom: 'var(--space-4)' }}>Empieza ahora</h2>
          <p style={{ marginBottom: 'var(--space-8)' }}>Descarga el CLI y analiza tu primer repositorio en menos de un minuto.</p>
          <div style={{ display: 'flex', gap: 'var(--space-4)', justifyContent: 'center', flexWrap: 'wrap' }}>
            <Link href="/install" className="btn btn--primary btn--lg">Descargar CLI</Link>
            <Link href="/catalog" className="btn btn--secondary btn--lg">Explorar catálogo</Link>
          </div>
        </div>
      </section>
    </>
  )
}
