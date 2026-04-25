import apiService from './api'
import type { Memory } from '../types'

export interface MemoriesResponse {
  memories: Memory[]
  count: number
}

export interface AddMemoryRequest {
  type: string
  content: string
  metadata?: string
  [key: string]: unknown
}

export const memoryService = {
  fetchMemories: (type?: string) => {
    const params = type ? { type } : undefined
    return apiService.get<MemoriesResponse>('/memory', params as Record<string, unknown>)
  },
  addMemory: (data: AddMemoryRequest) => {
    return apiService.post('/memory', data as Record<string, unknown>)
  },
}

export default memoryService