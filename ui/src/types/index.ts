export interface BaseEntity {
  id: string
  created_at: string
  updated_at?: string
}

export interface PaginationParams {
  page?: number
  pageSize?: number
}

export interface ListResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export type LoadingState = 'idle' | 'loading' | 'succeeded' | 'failed'

export interface AsyncState<T> {
  data: T | null
  loading: boolean
  error: string | null
}

export type MemoryType = 'snapshot' | 'session' | 'longterm'

export interface SelectOption {
  value: string
  label: string
}

export const MEMORY_TYPE_OPTIONS: SelectOption[] = [
  { value: '', label: '全部类型' },
  { value: 'snapshot', label: '快照记忆' },
  { value: 'session', label: '会话记忆' },
  { value: 'longterm', label: '长期记忆' },
]

export interface Skill extends BaseEntity {
  name: string
  description: string
  version: string
  platforms: string
  content: string
}

export interface Memory extends BaseEntity {
  type: MemoryType
  content: string
  metadata: string
}

export interface Session extends BaseEntity {
}

export interface Stats {
  tool_call_count: number
  skill_count: number
  memory_stats: {
    snapshot: number
    session: number
    longterm: number
  }
  conversation_len: number
}