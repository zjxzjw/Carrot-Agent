import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import { memoryService } from '../services'
import type { Memory } from '../types'

export interface MemoryState {
  memories: Memory[]
  loading: boolean
  error: string | null
}

const initialState: MemoryState = {
  memories: [],
  loading: false,
  error: null,
}

export const fetchMemories = createAsyncThunk(
  'memory/fetchMemories',
  async (type: string = '', { rejectWithValue }) => {
    try {
      const response = await memoryService.fetchMemories(type)
      return response.data.memories
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '获取记忆列表失败')
    }
  }
)

export const addMemory = createAsyncThunk(
  'memory/addMemory',
  async (memory: { type: string; content: string; metadata?: string }, { rejectWithValue }) => {
    try {
      await memoryService.addMemory(memory)
      return memory
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '添加记忆失败')
    }
  }
)

const memorySlice = createSlice({
  name: 'memory',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchMemories.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchMemories.fulfilled, (state, action) => {
        state.loading = false
        state.memories = action.payload
      })
      .addCase(fetchMemories.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '获取记忆列表失败'
      })
      .addCase(addMemory.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(addMemory.fulfilled, (state) => {
        state.loading = false
      })
      .addCase(addMemory.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '添加记忆失败'
      })
  },
})

export default memorySlice.reducer