'use client'
import { useState } from 'react'
import type { ProjectGraphDTO, RecommendationResult } from '@/types/stackgenome'
import FileImporter from '@/components/report/FileImporter'
import Dashboard from '@/components/report/Dashboard'

export default function ReportPageClient() {
  const [graph, setGraph] = useState<ProjectGraphDTO | null>(null)
  
  // Handlers for remote recommendations to update the state of the loaded graph
  const handleRecommendations = (recs: RecommendationResult[]) => {
    setGraph(prev => prev ? { ...prev, recommendations: recs } : null)
  }

  return (
    <div className="section">
      <div className="container">
        {!graph ? (
          <>
            <div style={{ textAlign: 'center', marginBottom: 'var(--space-8)' }}>
              <h1 style={{ marginBottom: 'var(--space-2)' }}>Mi Reporte</h1>
              <p style={{ fontSize: 'var(--text-lg)', color: 'var(--text-secondary)' }}>
                Importa el JSON generado por el CLI.
              </p>
            </div>
            <FileImporter onImport={setGraph} />
          </>
        ) : (
          <Dashboard 
            graph={graph} 
            onClear={() => setGraph(null)} 
            onRecommendationsUpdate={handleRecommendations}
          />
        )}
      </div>
    </div>
  )
}
