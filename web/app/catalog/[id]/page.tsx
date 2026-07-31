import type { Metadata } from 'next'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { getCatalog } from '@/lib/api'

export async function generateMetadata({ params }: { params: { id: string } }): Promise<Metadata> {
  const res = await getCatalog()
  if (!res.ok) return { title: 'Herramienta no encontrada' }
  const r = res.data.find(r => r.id === params.id)
  if (!r) return { title: 'No encontrada' }
  return {
    title: r.canonical_name,
    description: r.summary,
  }
}

// Generate static pages for all tools in the catalog (SSG)
export async function generateStaticParams() {
  const res = await getCatalog()
  if (!res.ok) return []
  return res.data.map(r => ({ id: r.id }))
}

export default async function ResourcePage({ params }: { params: { id: string } }) {
  const res = await getCatalog()
  if (!res.ok) return <div className="container section">Error cargando el catálogo.</div>

  const resource = res.data.find(r => r.id === params.id)
  if (!resource) notFound()

  return (
    <div className="section">
      <div className="container" style={{ maxWidth: 800 }}>
        <div style={{ marginBottom: 'var(--space-6)' }}>
          <Link href="/catalog" style={{ color: 'var(--text-muted)', fontSize: 'var(--text-sm)' }}>
            ← Volver al catálogo
          </Link>
        </div>

        <div className="card">
          <div style={{ display: 'flex', gap: 'var(--space-3)', alignItems: 'center', marginBottom: 'var(--space-4)', flexWrap: 'wrap' }}>
            <h1 style={{ margin: 0 }}>{resource.canonical_name}</h1>
            <span className="badge badge--gray" style={{ fontFamily: 'var(--font-mono)' }}>{resource.id}</span>
            {resource.status !== 'active' && (
              <span className="badge badge--orange">{resource.status}</span>
            )}
          </div>

          <p style={{ fontSize: 'var(--text-lg)', marginBottom: 'var(--space-8)' }}>
            {resource.summary}
          </p>

          <div style={{ display: 'grid', gap: 'var(--space-6)', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', marginBottom: 'var(--space-8)' }}>
            <div>
              <h3 style={{ fontSize: 'var(--text-sm)', color: 'var(--text-muted)', marginBottom: 'var(--space-2)' }}>Tipo</h3>
              <div style={{ textTransform: 'capitalize' }}>{resource.type}</div>
            </div>
            <div>
              <h3 style={{ fontSize: 'var(--text-sm)', color: 'var(--text-muted)', marginBottom: 'var(--space-2)' }}>Licencia</h3>
              <div>{resource.license || 'No especificada'}</div>
            </div>
            <div>
              <h3 style={{ fontSize: 'var(--text-sm)', color: 'var(--text-muted)', marginBottom: 'var(--space-2)' }}>Ecosistemas</h3>
              <div style={{ display: 'flex', gap: 'var(--space-2)', flexWrap: 'wrap' }}>
                {resource.ecosystem.map(e => <span key={e} className="badge badge--gray">{e}</span>)}
              </div>
            </div>
            {resource.infra_targets.length > 0 && (
              <div>
                <h3 style={{ fontSize: 'var(--text-sm)', color: 'var(--text-muted)', marginBottom: 'var(--space-2)' }}>Infraestructura</h3>
                <div style={{ display: 'flex', gap: 'var(--space-2)', flexWrap: 'wrap' }}>
                  {resource.infra_targets.map(t => <span key={t} className="badge badge--gray">{t}</span>)}
                </div>
              </div>
            )}
          </div>

          {resource.canonical_url && (
            <div style={{ marginTop: 'var(--space-8)' }}>
              <Link href={resource.canonical_url} target="_blank" rel="noopener noreferrer" className="btn btn--secondary">
                Visitar sitio oficial ↗
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
