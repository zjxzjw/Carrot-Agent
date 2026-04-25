import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import { sessionsService } from '../services'
import type { Session } from '../types'

export interface SessionsState {
  sessions: Session[]
  loading: boolean
  error: string | null
}

const initialState: SessionsState = {
  sessions: [],
  loading: false,
  error: null,
}

export const fetchSessions = createAsyncThunk(
  'sessions/fetchSessions',
  async (_, { rejectWithValue }) => {
    try {
      const response = await sessionsService.fetchSessions()
      return response.data.sessions
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '获取会话列表失败')
    }
  }
)

export const deleteSession = createAsyncThunk(
  'sessions/deleteSession',
  async (sessionId: string, { rejectWithValue }) => {
    try {
      await sessionsService.deleteSession(sessionId)
      return sessionId
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '删除会话失败')
    }
  }
)

const sessionsSlice = createSlice({
  name: 'sessions',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchSessions.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchSessions.fulfilled, (state, action) => {
        state.loading = false
        state.sessions = action.payload
      })
      .addCase(fetchSessions.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '获取会话列表失败'
      })
      .addCase(deleteSession.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(deleteSession.fulfilled, (state, action) => {
        state.loading = false
        state.sessions = state.sessions.filter(session => session.id !== action.payload)
      })
      .addCase(deleteSession.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '删除会话失败'
      })
  },
})

export default sessionsSlice.reducer