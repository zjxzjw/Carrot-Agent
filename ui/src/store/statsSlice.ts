import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import { statsService } from '../services'

export interface MemoryStats {
  snapshot: number
  session: number
  longterm: number
}

export interface StatsState {
  tool_call_count: number
  skill_count: number
  memory_stats: MemoryStats
  conversation_len: number
  loading: boolean
  error: string | null
}

const initialState: StatsState = {
  tool_call_count: 0,
  skill_count: 0,
  memory_stats: {
    snapshot: 0,
    session: 0,
    longterm: 0,
  },
  conversation_len: 0,
  loading: false,
  error: null,
}

export const fetchStats = createAsyncThunk(
  'stats/fetchStats',
  async (_, { rejectWithValue }) => {
    try {
      const response = await statsService.fetchStats()
      return response.data
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '获取统计信息失败')
    }
  }
)

const statsSlice = createSlice({
  name: 'stats',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchStats.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchStats.fulfilled, (state, action) => {
        state.loading = false
        state.tool_call_count = action.payload.tool_call_count
        state.skill_count = action.payload.skill_count
        state.memory_stats = action.payload.memory_stats
        state.conversation_len = action.payload.conversation_len
      })
      .addCase(fetchStats.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '获取统计信息失败'
      })
  },
})

export default statsSlice.reducer