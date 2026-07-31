// Bindings available from wrangler.toml
export interface Env {
  DB: D1Database;
  ENVIRONMENT: string;
  RANKING_VERSION: string;
  MAX_PAYLOAD_BYTES: string;
  RATE_LIMIT_RPM: string;
}
