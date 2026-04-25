import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'
import { chatService } from '../services'

export interface ChatMessage {
  id: string
  content: string
  role: 'user' | 'assistant'
  timestamp: number
}

export interface ChatState {
  messages: ChatMessage[]
  loading: boolean
  error: string | null
}

const initialState: ChatState = {
  messages: [],
  loading: false,
  error: null,
}

export const sendMessage = createAsyncThunk(
  'chat/sendMessage',
  async (message: string, { rejectWithValue }) => {
    try {
      const response = await chatService.sendMessage({ message })
      return {
        message: response.data.message,
        usage: response.data.usage,
      }
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '发送消息失败')
    }
  }
)

const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    setSessionId: (_state, _action: PayloadAction<string>) => {
    },
    clearMessages: (state) => {
      state.messages = []
      state.error = null
    },
    clearError: (state) => {
      state.error = null
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(sendMessage.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(sendMessage.fulfilled, (state, action) => {
        state.loading = false
        const userMessage: ChatMessage = {
          id: `user-${Date.now()}`,
          content: action.meta.arg,
          role: 'user',
          timestamp: Date.now(),
        }
        const assistantMessage: ChatMessage = {
          id: `assistant-${Date.now()}`,
          content: action.payload.message,
          role: 'assistant',
          timestamp: Date.now(),
        }
        state.messages.push(userMessage, assistantMessage)
      })
      .addCase(sendMessage.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '发送消息失败'
      })
  },
})

export const { setSessionId, clearMessages, clearError } = chatSlice.actions
export default chatSlice.reducer