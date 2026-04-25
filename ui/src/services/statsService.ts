import apiService from './api'

export interface StatsResponse {
  tool_call_count: number
  skill_count: number
  memory_stats: {
    snapshot: number
    session: number
    longterm: number
  }
  conversation_len: number
}

export const statsService = {
  fetchStats: () => {
    return apiService.get<StatsResponse>('/api/stats')
  },
}

export default statsService