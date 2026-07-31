export default function CatalogFilters({
  ecosystems,
  types,
  selectedEcosystem,
  selectedType,
  onEcosystemChange,
  onTypeChange,
}: {
  ecosystems: string[]
  types: string[]
  selectedEcosystem: string
  selectedType: string
  onEcosystemChange: (val: string) => void
  onTypeChange: (val: string) => void
}) {
  return (
    <div style={{ display: 'flex', gap: 'var(--space-3)', flexWrap: 'wrap' }}>
      <select
        className="input"
        style={{ width: 'auto' }}
        value={selectedEcosystem}
        onChange={e => onEcosystemChange(e.target.value)}
        aria-label="Filtrar por ecosistema"
      >
        <option value="">Todos los ecosistemas</option>
        {ecosystems.map(e => (
          <option key={e} value={e}>{e}</option>
        ))}
      </select>

      <select
        className="input"
        style={{ width: 'auto' }}
        value={selectedType}
        onChange={e => onTypeChange(e.target.value)}
        aria-label="Filtrar por tipo"
      >
        <option value="">Todos los tipos</option>
        {types.map(t => (
          <option key={t} value={t}>{t}</option>
        ))}
      </select>
    </div>
  )
}
