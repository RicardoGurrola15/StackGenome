import type { Metadata } from 'next'
import './globals.css'
import Nav from '@/components/ui/Nav'
import Footer from '@/components/ui/Footer'

export const metadata: Metadata = {
  title: {
    default: 'StackGenome — Analiza el ADN de tu stack tecnológico',
    template: '%s | StackGenome',
  },
  description:
    'StackGenome analiza localmente tu repositorio y mapea lenguajes, frameworks, dependencias e infraestructura. Privacidad local-first: tu código nunca sale de tu máquina.',
  keywords: ['stack tecnológico', 'análisis de repositorio', 'CLI', 'developer tools', 'open source'],
  openGraph: {
    title: 'StackGenome — Analiza el ADN de tu stack',
    description: 'Análisis local-first de repositorios. Detecta lenguajes, frameworks, dependencias e infraestructura.',
    type: 'website',
    locale: 'es_MX',
    siteName: 'StackGenome',
  },
  twitter: { card: 'summary_large_image' },
  robots: { index: true, follow: true },
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="es">
      <body>
        <a href="#main-content" className="sr-only">Saltar al contenido principal</a>
        <Nav />
        <main id="main-content">
          {children}
        </main>
        <Footer />
      </body>
    </html>
  )
}
