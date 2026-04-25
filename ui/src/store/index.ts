import { configureStore } from '@reduxjs/toolkit'
import chatReducer from './chatSlice'
import skillsReducer from './skillsSlice'
import memoryReducer from './memorySlice'
import sessionsReducer from './sessionsSlice'
import statsReducer from './statsSlice'
import authReducer from './authSlice'

export const store = configureStore({
  reducer: {
    chat: chatReducer,
    skills: skillsReducer,
    memory: memoryReducer,
    sessions: sessionsReducer,
    stats: statsReducer,
    auth: authReducer,
  },
})  

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch