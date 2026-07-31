'use client'
import { useState } from 'react'
import type { ProjectGraphDTO, RecommendationResult } from '@/types/stackgenome'
import { extractFingerprint, describeFingerprintPayload } from '@/lib/fingerprint'
import { getRecommendations } from '@/lib/api'

export default function RemoteRecommend({ 
  graph, 
  onSuccess 
}: { 
  graph: ProjectGraphDTO
  onSuccess: (recs: RecommendationResult[]) => void 
}) {
  const [isOpen, setIsOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fp = extractFingerprint(graph)
  const desc = describeFingerprintPayload(fp)

  const handleSend = async () => {
    setLoading(true)
    setError(null)
    
    const res = await getRecommendations({
      schema_version: graph.schema_version,
      fingerprint: fp,
    })

    if (res.ok) {
      onSuccess(res.data.recommendations)
      setIsOpen(false)
    } else {
      setError(res.error)
    }
    setLoading(false)
  }

  if (!isOpen) {
    return (
      <div className="card" style={{ background: 'var(--bg-tertiary)', textAlign: 'center' }}>
        <h3 style={{ marginBottom: 'var(--space-2)' }}>Recomendaciones</h3>
        <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--text-sm)', marginBottom: 'var(--space-4)' }}>
          Las recomendaciones adjuntas en tu JSON son estáticas. 
          Puedes consultar el catálogo actual para obtener recomendaciones frescas basadas en tu stack.
        </p>
        <button className="btn btn--primary" onClick={() => setIsOpen(true)}>
          Consultar catálogo remoto
        </button>
      </div>
    )
  }

  return (
    <div className="card" style={{ borderColor: 'var(--accent-blue-dim)' }}>
      <div style={{ display: 'flex', gap: 'var(--space-3)', alignItems: 'center', marginBottom: 'var(--space-4)' }}>
        <span style={{ fontSize: '1.5rem' }}>🔒</span>
        <h3 style={{ margin: 0 }}>Confirmación de Privacidad</h3>
      </div>
      
      <p style={{ fontSize: 'var(--text-sm)', marginBottom: 'var(--space-4)' }}>
        Estás a punto de consultar el backend de StackGenome. Se enviará un <strong>Fingerprint</strong> anónimo 
        de tu stack para calcular las recomendaciones.
      </p>
      
      <div style={{ background: 'var(--bg-primary)', padding: 'var(--space-3)', borderRadius: 'var(--radius)', border: '1px solid var(--border)', marginBottom: 'var(--space-4)', fontSize: 'var(--text-sm)', fontFamily: 'var(--font-mono)' }}>
        <strong>Qué se envía:</strong><br/>
        <span className="text-green">{desc}</span><br/><br/>
        <strong>Qué NO se envía:</strong><br/>
        <span className="text-muted">Nombres de archivos, contenidos (evidencias), tu IP, o secretos.</span>
      </div>

      {error && (
        <div className="alert alert--error" style={{ marginBottom: 'var(--space-4)' }}>
          {error}
        </div>
      )}

      <div style={{ display: 'flex', gap: 'var(--space-3)' }}>
        <button className="btn btn--primary" onClick={handleSend} disabled={loading}>
          {loading ? 'Consultando...' : 'Aceptar y Enviar'}
        </button>
        <button className="btn btn--secondary" onClick={() => setIsOpen(false)} disabled={loading}>
          Cancelar
        </button>
      </div>
    </div>
  )
}
