import apiService from './api'

export interface ChatMessage {
  message: string
  session_id?: string
  [key: string]: unknown
}

export interface ChatResponse {
  message: string
  usage: {
    prompt_tokens: number
    completion_tokens: number
    total_tokens: number
  }
}

export const chatService = {
  sendMessage: (data: ChatMessage) => {
    return apiService.post<ChatResponse>('/api/chat', data as Record<string, unknown>)
  },
}

export default chatService