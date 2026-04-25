export interface Message {
  id: number
  content: string
  role: 'user' | 'ai'
  timestamp: Date
}

export interface Skill {
  id: string
  name: string
  description: string
  version: string
  platforms: string
  content: string
  created_at: string
  updated_at: string
}

export interface Memory {
  id: string
  type: string
  content: string
  metadata: string
  created_at: string
}

export interface Stats {
  tool_call_count: number
  skill_count: number
  memory_stats: Record<string, number>
  conversation_len: number
}

export interface Conversation {
  id: string
  title: string
  updatedAt: Date
}

export interface AgentState {
  messages: Message[]
  skills: Skill[]
  memories: Memory[]
  stats: Stats | null
  conversations: Conversation[]
  currentConversationId: string
  loading: boolean
  error: string | null
  connected: boolean
}

export type AgentAction =
  | { type: 'SET_MESSAGES'; payload: Message[] }
  | { type: 'ADD_MESSAGE'; payload: Message }
  | { type: 'SET_SKILLS'; payload: Skill[] }
  | { type: 'SET_MEMORIES'; payload: Memory[] }
  | { type: 'SET_STATS'; payload: Stats | null }
  | { type: 'SET_CONVERSATIONS'; payload: Conversation[] }
  | { type: 'SELECT_CONVERSATION'; payload: string }
  | { type: 'ADD_CONVERSATION'; payload: Conversation }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_CONNECTED'; payload: boolean }
  | { type: 'CLEAR_MESSAGES' }
  | { type: 'CLEAR_ERROR' }

const API_BASE = 'http://localhost:8080/api'

export const initialState: AgentState = {
  messages: [],
  skills: [],
  memories: [],
  stats: null,
  conversations: [
    { id: 'default', title: '对话 1', updatedAt: new Date() },
    { id: '2', title: '对话 2', updatedAt: new Date(Date.now() - 86400000) },
    { id: '3', title: '对话 3', updatedAt: new Date(Date.now() - 172800000) },
  ],
  currentConversationId: 'default',
  loading: false,
  error: null,
  connected: false
}

export function agentReducer(state: AgentState, action: AgentAction): AgentState {
  switch (action.type) {
    case 'SET_MESSAGES':
      return { ...state, messages: action.payload }
    case 'ADD_MESSAGE':
      return { ...state, messages: [...state.messages, action.payload] }
    case 'SET_SKILLS':
      return { ...state, skills: action.payload }
    case 'SET_MEMORIES':
      return { ...state, memories: action.payload }
    case 'SET_STATS':
      return { ...state, stats: action.payload }
    case 'SET_CONVERSATIONS':
      return { ...state, conversations: action.payload }
    case 'SELECT_CONVERSATION':
      return { ...state, currentConversationId: action.payload, messages: [] }
    case 'ADD_CONVERSATION':
      return { ...state, conversations: [action.payload, ...state.conversations] }
    case 'SET_LOADING':
      return { ...state, loading: action.payload }
    case 'SET_ERROR':
      return { ...state, error: action.payload }
    case 'SET_CONNECTED':
      return { ...state, connected: action.payload }
    case 'CLEAR_MESSAGES':
      return { ...state, messages: [] }
    case 'CLEAR_ERROR':
      return { ...state, error: null }
    default:
      return state
  }
}

export async function checkConnection(dispatch: React.Dispatch<AgentAction>) {
  try {
    const response = await fetch(`${API_BASE}/stats`)
    dispatch({ type: 'SET_CONNECTED', payload: response.ok })
  } catch {
    dispatch({ type: 'SET_CONNECTED', payload: false })
  }
}

export async function sendMessage(dispatch: React.Dispatch<AgentAction>, content: string, currentConversationId: string) {
  if (!content.trim()) return

  const userMessage: Message = {
    id: Date.now(),
    content: content.trim(),
    role: 'user',
    timestamp: new Date()
  }
  dispatch({ type: 'ADD_MESSAGE', payload: userMessage })

  dispatch({ type: 'SET_LOADING', payload: true })
  dispatch({ type: 'SET_ERROR', payload: null })

  try {
    const response = await fetch(`${API_BASE}/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: content,
        session_id: currentConversationId
      })
    })

    if (!response.ok) {
      throw new Error(`请求失败: ${response.status}`)
    }

    const data = await response.json()

    const aiMessage: Message = {
      id: Date.now() + 1,
      content: data.message || '无响应',
      role: 'ai',
      timestamp: new Date()
    }
    dispatch({ type: 'ADD_MESSAGE', payload: aiMessage })
  } catch (e) {
    dispatch({ type: 'SET_ERROR', payload: e instanceof Error ? e.message : '发送消息失败' })
  } finally {
    dispatch({ type: 'SET_LOADING', payload: false })
  }
}

export async function fetchSkills(dispatch: React.Dispatch<AgentAction>) {
  dispatch({ type: 'SET_LOADING', payload: true })
  dispatch({ type: 'SET_ERROR', payload: null })

  try {
    const response = await fetch(`${API_BASE}/skills`)
    if (!response.ok) throw new Error('获取技能列表失败')

    const data = await response.json()
    dispatch({ type: 'SET_SKILLS', payload: data.skills || [] })
  } catch (e) {
    dispatch({ type: 'SET_ERROR', payload: e instanceof Error ? e.message : '获取技能列表失败' })
  } finally {
    dispatch({ type: 'SET_LOADING', payload: false })
  }
}

export async function fetchMemories(dispatch: React.Dispatch<AgentAction>, memType?: string) {
  dispatch({ type: 'SET_LOADING', payload: true })
  dispatch({ type: 'SET_ERROR', payload: null })

  try {
    const url = memType
      ? `${API_BASE}/memory?type=${memType}`
      : `${API_BASE}/memory`

    const response = await fetch(url)
    if (!response.ok) throw new Error('获取记忆列表失败')

    const data = await response.json()
    dispatch({ type: 'SET_MEMORIES', payload: data.memories || [] })
  } catch (e) {
    dispatch({ type: 'SET_ERROR', payload: e instanceof Error ? e.message : '获取记忆列表失败' })
  } finally {
    dispatch({ type: 'SET_LOADING', payload: false })
  }
}

export async function fetchStats(dispatch: React.Dispatch<AgentAction>) {
  dispatch({ type: 'SET_LOADING', payload: true })
  dispatch({ type: 'SET_ERROR', payload: null })

  try {
    const response = await fetch(`${API_BASE}/stats`)
    if (!response.ok) throw new Error('获取统计数据失败')

    const data = await response.json()
    dispatch({ type: 'SET_STATS', payload: data })
  } catch (e) {
    dispatch({ type: 'SET_ERROR', payload: e instanceof Error ? e.message : '获取统计数据失败' })
  } finally {
    dispatch({ type: 'SET_LOADING', payload: false })
  }
}

export async function createSkill(dispatch: React.Dispatch<AgentAction>, name: string, description: string, content: string) {
  dispatch({ type: 'SET_LOADING', payload: true })
  dispatch({ type: 'SET_ERROR', payload: null })

  try {
    const response = await fetch(`${API_BASE}/skills`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, description, content })
    })

    if (!response.ok) throw new Error('创建技能失败')

    await fetchSkills(dispatch)
  } catch (e) {
    dispatch({ type: 'SET_ERROR', payload: e instanceof Error ? e.message : '创建技能失败' })
    throw e
  } finally {
    dispatch({ type: 'SET_LOADING', payload: false })
  }
}

export async function addMemory(dispatch: React.Dispatch<AgentAction>, type_: string, content: string, metadata?: string) {
  dispatch({ type: 'SET_LOADING', payload: true })
  dispatch({ type: 'SET_ERROR', payload: null })

  try {
    const response = await fetch(`${API_BASE}/memory`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: type_, content, metadata: metadata || '{}' })
    })

    if (!response.ok) throw new Error('添加记忆失败')

    await fetchMemories(dispatch)
  } catch (e) {
    dispatch({ type: 'SET_ERROR', payload: e instanceof Error ? e.message : '添加记忆失败' })
    throw e
  } finally {
    dispatch({ type: 'SET_LOADING', payload: false })
  }
}

export function selectConversation(dispatch: React.Dispatch<AgentAction>, id: string) {
  dispatch({ type: 'SELECT_CONVERSATION', payload: id })
}

export function addConversation(dispatch: React.Dispatch<AgentAction>) {
  const newId = Date.now().toString()
  const newConversation: Conversation = {
    id: newId,
    title: `新对话 ${initialState.conversations.length + 1}`,
    updatedAt: new Date()
  }
  dispatch({ type: 'ADD_CONVERSATION', payload: newConversation })
  dispatch({ type: 'SELECT_CONVERSATION', payload: newId })
}

export function clearMessages(dispatch: React.Dispatch<AgentAction>) {
  dispatch({ type: 'CLEAR_MESSAGES' })
}

export function clearError(dispatch: React.Dispatch<AgentAction>) {
  dispatch({ type: 'CLEAR_ERROR' })
}

export function getCurrentConversation(state: AgentState) {
  return state.conversations.find(c => c.id === state.currentConversationId)
}

export function getSortedMemories(state: AgentState) {
  return [...state.memories].sort((a, b) =>
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )
}

export function getSortedSkills(state: AgentState) {
  return [...state.skills].sort((a, b) =>
    new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
  )
}