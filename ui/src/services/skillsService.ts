import apiService from './api'
import type { Skill } from '../types'

export interface SkillsResponse {
  skills: Skill[]
  count: number
}

export interface CreateSkillRequest {
  name: string
  description: string
  content: string
  [key: string]: unknown
}

export const skillsService = {
  fetchSkills: () => {
    return apiService.get<SkillsResponse>('/skills')
  },
  createSkill: (data: CreateSkillRequest) => {
    return apiService.post('/skills', data as Record<string, unknown>)
  },
}

export default skillsService