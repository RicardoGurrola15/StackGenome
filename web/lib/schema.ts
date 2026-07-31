import Ajv, { ValidateFunction } from 'ajv'
import type { ProjectGraphDTO } from '@/types/stackgenome'

const SUPPORTED_VERSIONS = ['1.0.0']
const MAX_FILE_SIZE_BYTES = 5 * 1024 * 1024 // 5MB

const ajv = new Ajv({ allErrors: true })

// Inlined schema derived from docs/api/openapi.yaml — ProjectGraphDTO
const projectGraphSchema = {
  type: 'object',
  required: ['schema_version', 'nodes', 'edges'],
  additionalProperties: true,
  properties: {
    schema_version: { type: 'string' },
    nodes: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'type', 'name', 'confidence'],
        properties: {
          id: { type: 'string' },
          type: { type: 'string' },
          name: { type: 'string' },
          confidence: { type: 'number' },
        },
      },
    },
    edges: {
      type: 'array',
      items: {
        type: 'object',
        required: ['source_id', 'target_id', 'relation'],
        properties: {
          source_id: { type: 'string' },
          target_id: { type: 'string' },
          relation: { type: 'string' },
        },
      },
    },
  },
}

let validate: ValidateFunction | null = null

function getValidator(): ValidateFunction {
  if (!validate) {
    validate = ajv.compile(projectGraphSchema)
  }
  return validate
}

export type ParseResult =
  | { ok: true; data: ProjectGraphDTO }
  | { ok: false; error: string }

export function parseGraphFile(content: string): ParseResult {
  let parsed: unknown
  try {
    parsed = JSON.parse(content)
  } catch {
    return { ok: false, error: 'El archivo no es JSON válido.' }
  }

  const v = getValidator()
  const valid = v(parsed)

  if (!valid) {
    const errs = v.errors?.map(e => `${e.instancePath} ${e.message}`).join('; ') ?? 'Schema inválido'
    return { ok: false, error: `El archivo no cumple el schema de StackGenome: ${errs}` }
  }

  const dto = parsed as ProjectGraphDTO

  if (!SUPPORTED_VERSIONS.includes(dto.schema_version)) {
    return {
      ok: false,
      error: `Versión de schema no compatible: "${dto.schema_version}". Versiones soportadas: ${SUPPORTED_VERSIONS.join(', ')}. Regenera el reporte con la versión actual del CLI.`,
    }
  }

  return { ok: true, data: dto }
}

export function checkFileSize(file: File): boolean {
  return file.size <= MAX_FILE_SIZE_BYTES
}

export const MAX_FILE_SIZE_MB = MAX_FILE_SIZE_BYTES / (1024 * 1024)
