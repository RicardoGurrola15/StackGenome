'use client'
import { useState, useCallback } from 'react'
import { parseGraphFile, checkFileSize, MAX_FILE_SIZE_MB } from '@/lib/schema'
import type { ProjectGraphDTO } from '@/types/stackgenome'

export default function FileImporter({ onImport }: { onImport: (data: ProjectGraphDTO) => void }) {
  const [error, setError] = useState<string | null>(null)
  const [isDragging, setIsDragging] = useState(false)

  const handleFile = async (file: File) => {
    setError(null)
    if (!checkFileSize(file)) {
      setError(`El archivo es demasiado grande (máx ${MAX_FILE_SIZE_MB}MB)`)
      return
    }

    try {
      const text = await file.text()
      const result = parseGraphFile(text)
      if (result.ok) {
        onImport(result.data)
      } else {
        setError(result.error)
      }
    } catch {
      setError('No se pudo leer el archivo')
    }
  }

  const onDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    const file = e.dataTransfer.files?.[0]
    if (file) handleFile(file)
  }, [])

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) handleFile(file)
  }

  return (
    <div style={{ maxWidth: 600, margin: '0 auto' }}>
      <div
        style={{
          border: `2px dashed ${isDragging ? 'var(--accent-blue)' : 'var(--border)'}`,
          borderRadius: 'var(--radius-lg)',
          padding: 'var(--space-12) var(--space-6)',
          textAlign: 'center',
          background: isDragging ? 'rgba(88,166,255,0.05)' : 'var(--bg-secondary)',
          transition: 'all var(--transition)',
          cursor: 'pointer',
        }}
        onDragOver={e => { e.preventDefault(); setIsDragging(true) }}
        onDragLeave={() => setIsDragging(false)}
        onDrop={onDrop}
        onClick={() => document.getElementById('file-upload')?.click()}
      >
        <div style={{ fontSize: '3rem', marginBottom: 'var(--space-4)' }}>📁</div>
        <h3 style={{ marginBottom: 'var(--space-2)' }}>Arrastra tu archivo aquí</h3>
        <p style={{ color: 'var(--text-muted)', fontSize: 'var(--text-sm)', marginBottom: 'var(--space-6)' }}>
          o haz clic para seleccionar (solo .json)
        </p>
        <p style={{ color: 'var(--text-muted)', fontSize: 'var(--text-xs)', fontFamily: 'var(--font-mono)' }}>
          El archivo se procesa 100% en tu navegador.
        </p>
        <input
          id="file-upload"
          type="file"
          accept="application/json,.json"
          onChange={onChange}
          style={{ display: 'none' }}
        />
      </div>

      {error && (
        <div className="alert alert--error" style={{ marginTop: 'var(--space-4)' }} role="alert">
          <span>❌</span>
          <div>{error}</div>
        </div>
      )}
    </div>
  )
}
