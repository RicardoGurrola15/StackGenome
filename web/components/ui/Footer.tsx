import Link from 'next/link'

export default function Footer() {
  return (
    <footer style={{
      borderTop: '1px solid var(--border)',
      background: 'var(--bg-secondary)',
      padding: 'var(--space-8) 0',
      marginTop: 'auto',
    }}>
      <div className="container" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 'var(--space-4)' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-2)' }}>
          <span style={{ color: 'var(--accent-green)', fontSize: '1.25rem' }}>⬡</span>
          <span style={{ fontWeight: 600, fontSize: 'var(--text-sm)', color: 'var(--text-secondary)' }}>
            StackGenome Alpha — Local-first. Sin cuentas. Sin upload automático.
          </span>
        </div>
        <div style={{ display: 'flex', gap: 'var(--space-6)', fontSize: 'var(--text-sm)', color: 'var(--text-muted)' }}>
          <Link href="https://github.com/RicardoGurrola15/StackGenome" target="_blank" rel="noopener noreferrer">GitHub ↗</Link>
          <Link href="/install">Instalar CLI</Link>
          <Link href="/catalog">Catálogo</Link>
          <Link href="/report">Mi Reporte</Link>
        </div>
      </div>
    </footer>
  )
}
