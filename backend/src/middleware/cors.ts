const ALLOWED_ORIGINS = [
  'http://localhost:3000',
  'http://localhost:5173',
  'https://staging.stackgenome.com',
  'https://stackgenome.com',
];

function isAllowedOrigin(origin: string | null): boolean {
  if (!origin) return false;
  if (ALLOWED_ORIGINS.includes(origin)) return true;
  // Permitir dominios de preview (ej: Cloudflare Pages)
  if (origin.endsWith('.stackgenome.pages.dev')) return true;
  return false;
}

export function corsHeaders(origin: string | null): Headers {
  const headers = new Headers();
  
  if (isAllowedOrigin(origin)) {
    headers.set('Access-Control-Allow-Origin', origin!);
  } else {
    // Fallback estricto
    headers.set('Access-Control-Allow-Origin', 'https://stackgenome.com');
  }

  headers.set('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  headers.set('Access-Control-Allow-Headers', 'Content-Type');
  headers.set('Access-Control-Max-Age', '86400');
  return headers;
}

export function handleOptions(request: Request): Response {
  const origin = request.headers.get('Origin');
  return new Response(null, {
    status: 204,
    headers: corsHeaders(origin),
  });
}

export function applyCors(req: Request, res: Response): Response {
  const origin = req.headers.get('Origin');
  const headers = new Headers(res.headers);
  const cors = corsHeaders(origin);
  for (const [key, value] of cors.entries()) {
    headers.set(key, value);
  }
  return new Response(res.body, {
    status: res.status,
    statusText: res.statusText,
    headers,
  });
}
