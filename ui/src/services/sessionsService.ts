import apiService from './api'
import type { Session } from '../types'

export interface SessionsResponse {
  sessions: Session[]
  count: number
}

export interface DeleteSessionResponse {
  message: string
}

export const sessionsService = {
  fetchSessions: () => {
    return apiService.get<SessionsResponse>('/session')
  },
  deleteSession: (sessionId: string) => {
    return apiService.delete<DeleteSessionResponse>(`/session/${sessionId}`)
  },
}

export default sessionsService