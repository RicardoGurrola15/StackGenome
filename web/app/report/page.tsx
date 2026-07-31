import type { Metadata } from 'next'
import ReportClient from './ReportClient'

export const metadata: Metadata = {
  title: 'Mi Reporte',
  description: 'Visualiza y explora tu StackGraph generado localmente.',
}

export default function ReportPage() {
  return <ReportClient />
}
