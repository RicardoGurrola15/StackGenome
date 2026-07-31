'use client'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import styles from './Nav.module.css'

const links = [
  { href: '/', label: 'Inicio' },
  { href: '/install', label: 'Instalar' },
  { href: '/catalog', label: 'Catálogo' },
  { href: '/report', label: 'Mi Reporte' },
]

export default function Nav() {
  const pathname = usePathname()
  return (
    <header className={styles.header} role="banner">
      <nav className={styles.nav} aria-label="Navegación principal">
        <Link href="/" className={styles.logo} aria-label="StackGenome — Inicio">
          <span className={styles.logoIcon}>⬡</span>
          <span className={styles.logoText}>StackGenome</span>
          <span className={styles.logoBadge}>alpha</span>
        </Link>
        <ul className={styles.links} role="list">
          {links.map(({ href, label }) => (
            <li key={href}>
              <Link
                href={href}
                className={`${styles.link} ${pathname === href ? styles.active : ''}`}
                aria-current={pathname === href ? 'page' : undefined}
              >
                {label}
              </Link>
            </li>
          ))}
        </ul>
        <Link href="https://github.com/RicardoGurrola15/StackGenome" target="_blank" rel="noopener noreferrer" className="btn btn--secondary btn--sm" aria-label="Ver en GitHub (abre en nueva pestaña)">
          GitHub ↗
        </Link>
      </nav>
    </header>
  )
}
