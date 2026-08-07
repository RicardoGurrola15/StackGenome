'use client'
import { useState, useEffect } from 'react'
import type { CatalogResource } from '@/types/stackgenome'
import { getCatalog } from '@/lib/api'
import ResourceCard from '@/components/catalog/ResourceCard'
import CatalogFilters from '@/components/catalog/CatalogFilters'
import ResourceModal from '@/components/catalog/ResourceModal'
import styles from './Catalog.module.css'

export default function CatalogClient() {
  const [resources, setResources] = useState<CatalogResource[]>([])
  const [filtered, setFiltered] = useState<CatalogResource[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [query, setQuery] = useState('')
  const [ecosystem, setEcosystem] = useState('')
  const [type, setType] = useState('')
  const [selectedResource, setSelectedResource] = useState<CatalogResource | null>(null)

  useEffect(() => {
    getCatalog().then(res => {
      if (res.ok) {
        setResources(res.data)
        setFiltered(res.data)
      } else {
        setError(res.error)
      }
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    let out = resources
    if (query) {
      const q = query.toLowerCase()
      out = out.filter(r =>
        r.canonical_name.toLowerCase().includes(q) ||
        r.summary.toLowerCase().includes(q)
      )
    }
    if (ecosystem) {
      out = out.filter(r => r.ecosystem.includes(ecosystem))
    }
    if (type) {
      out = out.filter(r => r.type === type)
    }
    setFiltered(out)
  }, [query, ecosystem, type, resources])

  const ecosystems = [...new Set(resources.flatMap(r => r.ecosystem))].sort()
  const types = [...new Set(resources.map(r => r.type))].sort()

  if (loading) {
    return (
      <div>
        <div className={styles.searchBar}>
          <div className="skeleton" style={{ height: 40, borderRadius: 'var(--radius)' }} />
        </div>
        <div className="grid-auto">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="skeleton" style={{ height: 140, borderRadius: 'var(--radius-lg)' }} />
          ))}
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="alert alert--error" role="alert">
        <span>❌</span>
        <div>
          <strong>No se pudo cargar el catálogo.</strong> {error}
          <br />
          <span style={{ fontSize: 'var(--text-xs)', color: 'var(--text-muted)' }}>
            El catálogo requiere conexión al backend de staging.
          </span>
        </div>
      </div>
    )
  }

  return (
    <div>
      {/* Search + Filters */}
      <div className={styles.toolbar} role="search">
        <label htmlFor="catalog-search" className="sr-only">Buscar en el catálogo</label>
        <input
          id="catalog-search"
          className="input"
          type="search"
          placeholder="Buscar herramienta, librería, servicio..."
          value={query}
          onChange={e => setQuery(e.target.value)}
          aria-label="Buscar en el catálogo"
        />
        <CatalogFilters
          ecosystems={ecosystems}
          types={types}
          selectedEcosystem={ecosystem}
          selectedType={type}
          onEcosystemChange={setEcosystem}
          onTypeChange={setType}
        />
      </div>

      {/* Results count */}
      <p style={{ fontSize: 'var(--text-sm)', color: 'var(--text-muted)', marginBottom: 'var(--space-5)' }}
         aria-live="polite" aria-atomic="true">
        {filtered.length} {filtered.length === 1 ? 'resultado' : 'resultados'}
        {(query || ecosystem || type) && ' para la búsqueda actual'}
      </p>

      {/* Grid */}
      {filtered.length === 0 ? (
        <div className={styles.emptyState} role="status">
          <span style={{ fontSize: '2rem' }}>🔍</span>
          <p>No se encontraron herramientas con esos filtros.</p>
          <button className="btn btn--ghost" onClick={() => { setQuery(''); setEcosystem(''); setType('') }}>
            Limpiar filtros
          </button>
        </div>
      ) : (
        <div className="grid-auto">
          {filtered.map(r => <ResourceCard key={r.id} resource={r} onClick={() => setSelectedResource(r)} />)}
        </div>
      )}

      {/* Modal */}
      {selectedResource && (
        <ResourceModal resource={selectedResource} onClose={() => setSelectedResource(null)} />
      )}
    </div>
  )
}
