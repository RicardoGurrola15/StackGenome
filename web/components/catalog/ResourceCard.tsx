import type { CatalogResource } from '@/types/stackgenome'
import Link from 'next/link'
import styles from './ResourceCard.module.css'

const TYPE_COLORS: Record<string, string> = {
  tool: 'badge--blue',
  library: 'badge--green',
  service: 'badge--purple',
  platform: 'badge--orange',
}

export default function ResourceCard({ resource: r }: { resource: CatalogResource }) {
  return (
    <Link href={`/catalog/${encodeURIComponent(r.id)}`} className={`card card--interactive ${styles.card}`} aria-label={`Ver detalle de ${r.canonical_name}`}>
      <div className={styles.header}>
        <span className={`badge ${TYPE_COLORS[r.type] ?? 'badge--gray'}`}>{r.type}</span>
        {r.status !== 'active' && (
          <span className="badge badge--orange">{r.status}</span>
        )}
      </div>
      <h3 className={styles.name}>{r.canonical_name}</h3>
      <p className={styles.summary}>{r.summary}</p>
      <div className={styles.footer}>
        {r.ecosystem.slice(0, 3).map(e => (
          <span key={e} className="badge badge--gray">{e}</span>
        ))}
        {r.ecosystem.length > 3 && (
          <span className="badge badge--gray">+{r.ecosystem.length - 3}</span>
        )}
      </div>
    </Link>
  )
}
