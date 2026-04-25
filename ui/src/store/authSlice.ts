import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface AuthState {
  isAuthenticated: boolean
  sessionId: string | null
}

const initialState: AuthState = {
  isAuthenticated: false,
  sessionId: localStorage.getItem('sessionId') || null,
}

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    login: (state, action: PayloadAction<string>) => {
      state.isAuthenticated = true
      state.sessionId = action.payload
      localStorage.setItem('sessionId', action.payload)
    },
    logout: (state) => {
      state.isAuthenticated = false
      state.sessionId = null
      localStorage.removeItem('sessionId')
    },
  },
})

export const { login, logout } = authSlice.actions

export default authSlice.reducer