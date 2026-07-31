import type { Metadata } from 'next'
import CatalogClient from '@/components/catalog/CatalogClient'

export const metadata: Metadata = {
  title: 'Catálogo de Herramientas',
  description: 'Explora herramientas, frameworks y servicios. El catálogo de StackGenome base para las recomendaciones.',
}

export default function CatalogPage() {
  return (
    <div className="section">
      <div className="container" style={{ maxWidth: 1000 }}>
        <div style={{ marginBottom: 'var(--space-8)' }}>
          <h1 style={{ marginBottom: 'var(--space-2)' }}>Catálogo de Herramientas</h1>
          <p style={{ fontSize: 'var(--text-lg)' }}>
            Base de conocimiento utilizada para emitir recomendaciones deterministas.
          </p>
        </div>
        <CatalogClient />
      </div>
    </div>
  )
}
