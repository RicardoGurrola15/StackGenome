import type { CatalogResource } from '@/types/stackgenome'
import styles from './ResourceModal.module.css'
import { useEffect } from 'react'

interface Props {
  resource: CatalogResource
  onClose: () => void
}

export default function ResourceModal({ resource, onClose }: Props) {
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handleEsc)
    return () => window.removeEventListener('keydown', handleEsc)
  }, [onClose])

  return (
    <div className={styles.overlay} onClick={onClose} aria-modal="true" role="dialog">
      <div className={styles.modal} onClick={e => e.stopPropagation()}>
        <button className={styles.closeBtn} onClick={onClose} aria-label="Cerrar modal">
          ✕
        </button>
        
        <div style={{ display: 'flex', gap: 'var(--space-3)', alignItems: 'center', marginBottom: 'var(--space-4)', flexWrap: 'wrap' }}>
          <h2 style={{ margin: 0, fontSize: 'var(--text-3xl)' }}>{resource.canonical_name}</h2>
          <span className="badge badge--gray" style={{ fontFamily: 'var(--font-mono)' }}>{resource.id}</span>
          {resource.status !== 'active' && (
            <span className="badge badge--orange">{resource.status}</span>
          )}
        </div>

        <p style={{ fontSize: 'var(--text-lg)', marginBottom: 'var(--space-8)', color: 'var(--text-secondary)' }}>
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
          {resource.infra_targets && resource.infra_targets.length > 0 && (
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
            <a href={resource.canonical_url} target="_blank" rel="noopener noreferrer" className="btn btn--primary btn--lg">
              Visitar sitio oficial ↗
            </a>
          </div>
        )}
      </div>
    </div>
  )
}
