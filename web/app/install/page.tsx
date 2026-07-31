import type { Metadata } from 'next'
import Link from 'next/link'
import styles from './page.module.css'

export const metadata: Metadata = {
  title: 'Instalar StackGenome CLI',
  description: 'Descarga e instala el CLI de StackGenome para macOS, Linux o Windows. Sin dependencias externas.',
}

// These reference actual GitHub Releases — update on each real release
const REPO = 'https://github.com/RicardoGurrola15/StackGenome'
const RELEASES = `${REPO}/releases`
const VERSION = 'dev' // Updated by release workflow

const PLATFORMS = [
  {
    id: 'macos-arm',
    os: 'macOS',
    arch: 'Apple Silicon (arm64)',
    icon: '🍎',
    steps: [
      {
        label: 'Descarga el binario',
        code: `curl -L ${RELEASES}/latest/download/stackgenome-darwin-arm64 -o stackgenome`,
      },
      { label: 'Hazlo ejecutable', code: 'chmod +x stackgenome' },
      { label: 'Muévelo a tu PATH', code: 'sudo mv stackgenome /usr/local/bin/' },
      { label: 'Verifica la instalación', code: 'stackgenome version' },
    ],
  },
  {
    id: 'macos-amd',
    os: 'macOS',
    arch: 'Intel (amd64)',
    icon: '🍎',
    steps: [
      {
        label: 'Descarga el binario',
        code: `curl -L ${RELEASES}/latest/download/stackgenome-darwin-amd64 -o stackgenome`,
      },
      { label: 'Hazlo ejecutable', code: 'chmod +x stackgenome' },
      { label: 'Muévelo a tu PATH', code: 'sudo mv stackgenome /usr/local/bin/' },
      { label: 'Verifica la instalación', code: 'stackgenome version' },
    ],
  },
  {
    id: 'linux-amd',
    os: 'Linux',
    arch: 'amd64',
    icon: '🐧',
    steps: [
      {
        label: 'Descarga el binario',
        code: `curl -L ${RELEASES}/latest/download/stackgenome-linux-amd64 -o stackgenome`,
      },
      { label: 'Hazlo ejecutable', code: 'chmod +x stackgenome' },
      { label: 'Muévelo a tu PATH', code: 'sudo mv stackgenome /usr/local/bin/' },
      { label: 'Verifica la instalación', code: 'stackgenome version' },
    ],
  },
  {
    id: 'linux-arm',
    os: 'Linux',
    arch: 'arm64',
    icon: '🐧',
    steps: [
      {
        label: 'Descarga el binario',
        code: `curl -L ${RELEASES}/latest/download/stackgenome-linux-arm64 -o stackgenome`,
      },
      { label: 'Hazlo ejecutable', code: 'chmod +x stackgenome' },
      { label: 'Muévelo a tu PATH', code: 'sudo mv stackgenome /usr/local/bin/' },
      { label: 'Verifica la instalación', code: 'stackgenome version' },
    ],
  },
  {
    id: 'windows',
    os: 'Windows',
    arch: 'amd64',
    icon: '🪟',
    steps: [
      {
        label: 'Descarga el ejecutable',
        code: `# PowerShell\nInvoke-WebRequest -Uri "${RELEASES}/latest/download/stackgenome-windows-amd64.exe" -OutFile stackgenome.exe`,
      },
      { label: 'Agrégalo a tu PATH', code: '# Mueve stackgenome.exe a una carpeta en tu PATH' },
      { label: 'Verifica la instalación', code: 'stackgenome.exe version' },
    ],
  },
]

export default function InstallPage() {
  return (
    <div className="section">
      <div className="container" style={{ maxWidth: 900 }}>
        {/* Header */}
        <div style={{ marginBottom: 'var(--space-12)' }}>
          <div className="badge badge--blue" style={{ marginBottom: 'var(--space-4)' }}>CLI</div>
          <h1 style={{ marginBottom: 'var(--space-4)' }}>Instalar StackGenome</h1>
          <p style={{ fontSize: 'var(--text-lg)', maxWidth: 600 }}>
            Un solo binario. Sin Node, sin Docker, sin dependencias de runtime.
            Elige tu plataforma y empieza a analizar en segundos.
          </p>
        </div>

        {/* Alpha notice */}
        <div className="alert alert--warning" style={{ marginBottom: 'var(--space-8)' }} role="note">
          <span>⚠️</span>
          <div>
            <strong>Alpha: </strong>
            Los releases automatizados aún no están publicados (requieren un tag de versión).
            Puedes{' '}
            <Link href={`${REPO}`} target="_blank" rel="noopener noreferrer">
              compilar desde el código fuente ↗
            </Link>{' '}
            o esperar la primera release oficial.
          </div>
        </div>

        {/* Platforms */}
        <div className={styles.platforms}>
          {PLATFORMS.map(p => (
            <div key={p.id} className="card" id={p.id}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-3)', marginBottom: 'var(--space-5)' }}>
                <span style={{ fontSize: '1.5rem' }}>{p.icon}</span>
                <div>
                  <div style={{ fontWeight: 700, fontSize: 'var(--text-base)' }}>{p.os}</div>
                  <div style={{ fontSize: 'var(--text-sm)', color: 'var(--text-muted)', fontFamily: 'var(--font-mono)' }}>{p.arch}</div>
                </div>
              </div>
              <ol className={styles.steps}>
                {p.steps.map((s, i) => (
                  <li key={i} className={styles.step}>
                    <div className={styles.stepLabel}>{i + 1}. {s.label}</div>
                    <pre className={styles.codeBlock}><code>{s.code}</code></pre>
                  </li>
                ))}
              </ol>
            </div>
          ))}
        </div>

        <hr className="divider" />

        {/* Usage */}
        <section aria-labelledby="usage-title">
          <h2 id="usage-title" style={{ marginBottom: 'var(--space-6)' }}>Primeros pasos</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
            {[
              {
                title: 'Analizar un repositorio',
                code: 'stackgenome analyze /ruta/a/tu/proyecto',
              },
              {
                title: 'Exportar como JSON (para importar en la web)',
                code: 'stackgenome analyze --json /ruta/a/tu/proyecto > mi-reporte.json',
              },
              {
                title: 'Ver la versión',
                code: 'stackgenome version',
              },
              {
                title: 'Ver ayuda',
                code: 'stackgenome --help',
              },
            ].map(cmd => (
              <div key={cmd.title} className="card">
                <div style={{ fontSize: 'var(--text-sm)', color: 'var(--text-secondary)', marginBottom: 'var(--space-2)' }}>{cmd.title}</div>
                <pre style={{ background: 'var(--bg-primary)', border: '1px solid var(--border)', borderRadius: 'var(--radius)', padding: 'var(--space-3) var(--space-4)', overflow: 'auto' }}>
                  <code style={{ color: 'var(--accent-green)', fontSize: 'var(--text-sm)' }}>{cmd.code}</code>
                </pre>
              </div>
            ))}
          </div>
        </section>

        <hr className="divider" />

        {/* Compile from source */}
        <section aria-labelledby="source-title">
          <h2 id="source-title" style={{ marginBottom: 'var(--space-4)' }}>Compilar desde el código fuente</h2>
          <p style={{ marginBottom: 'var(--space-5)' }}>Requiere Go 1.22+</p>
          <pre className={styles.codeBlock}>
            <code>{`git clone https://github.com/RicardoGurrola15/StackGenome.git\ncd StackGenome/StackGenome\ngo build -o stackgenome ./cmd/stackgenome`}</code>
          </pre>
        </section>

        <div style={{ marginTop: 'var(--space-10)', textAlign: 'center' }}>
          <Link href="/report" className="btn btn--primary btn--lg">
            Importar mi reporte →
          </Link>
        </div>
      </div>
    </div>
  )
}
