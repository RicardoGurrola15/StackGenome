import type { ProjectGraphDTO, RecommendationResult } from '@/types/stackgenome'
import StackGraph from './StackGraph'
import RemoteRecommend from './RemoteRecommend'

export default function Dashboard({ 
  graph, 
  onClear,
  onRecommendationsUpdate
}: { 
  graph: ProjectGraphDTO
  onClear: () => void
  onRecommendationsUpdate: (recs: RecommendationResult[]) => void
}) {
  const nodesByType = graph.nodes.reduce((acc, node) => {
    if (!acc[node.type]) acc[node.type] = []
    acc[node.type].push(node)
    return acc
  }, {} as Record<string, typeof graph.nodes>)

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 'var(--space-6)', flexWrap: 'wrap', gap: 'var(--space-4)' }}>
        <div>
          <h2 style={{ marginBottom: 'var(--space-1)' }}>Dashboard del Proyecto</h2>
          <div className="badge badge--gray">Esquema v{graph.schema_version}</div>
        </div>
        <button className="btn btn--secondary btn--sm" onClick={onClear}>
          Cerrar reporte
        </button>
      </div>

      <div style={{ display: 'grid', gap: 'var(--space-6)', gridTemplateColumns: '1fr 300px', alignItems: 'start' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-6)' }}>
          {/* Stack Graph */}
          <section className="card" style={{ padding: 0, overflow: 'hidden' }}>
            <div style={{ padding: 'var(--space-4)', borderBottom: '1px solid var(--border)' }}>
              <h3 style={{ fontSize: 'var(--text-base)' }}>Visualización del Grafo</h3>
            </div>
            <StackGraph nodes={graph.nodes} edges={graph.edges} />
          </section>

          {/* Nodos */}
          <section className="card">
            <h3 style={{ marginBottom: 'var(--space-4)' }}>Inventario Detectado</h3>
            <div className="grid-auto">
              {Object.entries(nodesByType).map(([type, nodes]) => (
                <div key={type}>
                  <h4 style={{ textTransform: 'capitalize', fontSize: 'var(--text-sm)', color: 'var(--text-muted)', marginBottom: 'var(--space-3)', borderBottom: '1px solid var(--border)', paddingBottom: 'var(--space-1)' }}>
                    {type} ({nodes.length})
                  </h4>
                  <ul style={{ listStyle: 'none', display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
                    {nodes.map(n => (
                      <li key={n.id} style={{ fontSize: 'var(--text-sm)', display: 'flex', justifyContent: 'space-between' }}>
                        <span>{n.name} {n.version && <span className="text-muted">v{n.version}</span>}</span>
                        <span className="badge badge--gray">{Math.round(n.confidence * 100)}%</span>
                      </li>
                    ))}
                  </ul>
                </div>
              ))}
            </div>
          </section>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-6)' }}>
          {/* Entorno */}
          {graph.environment && (
            <section className="card">
              <h3 style={{ fontSize: 'var(--text-base)', marginBottom: 'var(--space-4)' }}>Entorno de análisis</h3>
              {graph.environment.os && (
                <div style={{ fontSize: 'var(--text-sm)', marginBottom: 'var(--space-2)' }}>
                  <span className="text-muted">OS:</span> {graph.environment.os} ({graph.environment.arch})
                </div>
              )}
              {graph.environment.tools && graph.environment.tools.length > 0 && (
                <div style={{ marginTop: 'var(--space-3)' }}>
                  <div className="text-muted" style={{ fontSize: 'var(--text-sm)', marginBottom: 'var(--space-2)' }}>Toolchain detectado:</div>
                  <ul style={{ listStyle: 'none', display: 'flex', flexDirection: 'column', gap: 'var(--space-1)', fontSize: 'var(--text-sm)' }}>
                    {graph.environment.tools.map(t => (
                      <li key={t.name}>{t.name} <span className="text-muted">{t.version}</span></li>
                    ))}
                  </ul>
                </div>
              )}
            </section>
          )}

          {/* Recomendaciones remotas */}
          <RemoteRecommend graph={graph} onSuccess={onRecommendationsUpdate} />

          {/* Lista de recomendaciones actuales */}
          {graph.recommendations && graph.recommendations.length > 0 && (
            <section className="card" style={{ borderColor: 'var(--accent-green-dim)' }}>
              <h3 style={{ fontSize: 'var(--text-base)', marginBottom: 'var(--space-4)', color: 'var(--accent-green)' }}>
                Recomendaciones ({graph.recommendations.length})
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
                {graph.recommendations.map(r => (
                  <div key={r.id} style={{ paddingBottom: 'var(--space-3)', borderBottom: '1px solid var(--border-subtle)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 'var(--space-2)' }}>
                      <strong style={{ fontSize: 'var(--text-sm)' }}>{r.name}</strong>
                      <span className="badge badge--green">Score: {Math.round(r.score * 100)}</span>
                    </div>
                    <p style={{ fontSize: 'var(--text-xs)', color: 'var(--text-secondary)', marginBottom: 'var(--space-2)' }}>{r.description}</p>
                    <ul style={{ margin: '0 0 0 var(--space-4)', padding: 0, fontSize: 'var(--text-xs)', color: 'var(--text-muted)' }}>
                      {r.reasons.map((reason, i) => (
                        <li key={i}>{reason}</li>
                      ))}
                    </ul>
                  </div>
                ))}
              </div>
            </section>
          )}
        </div>
      </div>
    </div>
  )
}
