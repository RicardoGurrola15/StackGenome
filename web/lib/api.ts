import type { RecommendationRequest, RecommendationResponse, CatalogResource, ErrorResponse } from '@/types/stackgenome'

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'https://stackgenome-api-staging.stackgenome.workers.dev'
const TIMEOUT_MS = 10_000

async function fetchWithTimeout(url: string, options: RequestInit): Promise<Response> {
  const controller = new AbortController()
  const id = setTimeout(() => controller.abort(), TIMEOUT_MS)
  try {
    return await fetch(url, { ...options, signal: controller.signal })
  } finally {
    clearTimeout(id)
  }
}

export type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; error: string; status?: number }

export async function getRecommendations(req: RecommendationRequest): Promise<ApiResult<RecommendationResponse>> {
  try {
    const res = await fetchWithTimeout(`${API_BASE}/api/v1/recommendations`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req),
    })

    if (!res.ok) {
      const err: Partial<ErrorResponse> = await res.json().catch(() => ({}))
      return { ok: false, error: err.message ?? `Error del servidor (${res.status})`, status: res.status }
    }

    const data: RecommendationResponse = await res.json()
    return { ok: true, data }
  } catch (e) {
    if (e instanceof DOMException && e.name === 'AbortError') {
      return { ok: false, error: 'Tiempo de espera agotado. El backend no respondió.' }
    }
    return { ok: false, error: 'No se pudo conectar al servidor. Verifica tu conexión.' }
  }
}

export async function getCatalog(): Promise<ApiResult<CatalogResource[]>> {
  try {
    const res = await fetchWithTimeout(`${API_BASE}/api/v1/resources`, { method: 'GET' })
    if (!res.ok) {
      return { ok: false, error: `Error al cargar el catálogo (${res.status})`, status: res.status }
    }
    const data = await res.json()
    return { ok: true, data: data.resources ?? data }
  } catch {
    return { ok: false, error: 'No se pudo cargar el catálogo. Verifica tu conexión.' }
  }
}
