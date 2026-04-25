import apiService from './api'

interface ModelConfig {
  provider: string
  api_key: string
  model_name: string
  base_url: string
  temperature: number
  max_tokens: number
}

interface Config {
  agent: {
    name: string
    version: string
    data_dir: string
    log_level: string
    skill_nudge_interval: number
  }
  model: ModelConfig
  storage: {
    db_path: string
    skill_dir: string
    memory_dir: string
    session_dir: string
  }
  server: {
    host: string
    port: number
    mode: string
  }
  security: {
    allowed_paths: string[]
    blocked_cmds: string[]
  }
}

interface ModelList {
  provider: string
  models: string[]
}

export type { Config, ModelConfig, ModelList }

export const configService = {
  getConfig: async () => {
    const response = await apiService.get<Config>('/config')
    return response
  },
  updateConfig: async (config: { model: Partial<ModelConfig> }) => {
    const response = await apiService.put<{ message: string; config: Config }>('/config', config)
    return response
  },
  getModels: async () => {
    const response = await apiService.get<ModelList[]>('/models')
    return response
  },
}

export default configService