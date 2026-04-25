import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import { skillsService } from '../services'
import type { Skill } from '../types'

export interface SkillsState {
  skills: Skill[]
  loading: boolean
  error: string | null
}

const initialState: SkillsState = {
  skills: [],
  loading: false,
  error: null,
}

export const fetchSkills = createAsyncThunk(
  'skills/fetchSkills',
  async (_, { rejectWithValue }) => {
    try {
      const response = await skillsService.fetchSkills()
      return response.data.skills
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '获取技能列表失败')
    }
  }
)

export const createSkill = createAsyncThunk(
  'skills/createSkill',
  async (skill: { name: string; description: string; content: string }, { rejectWithValue }) => {
    try {
      await skillsService.createSkill(skill)
      return skill
    } catch (error: unknown) {
      const err = error as { message?: string }
      return rejectWithValue(err.message || '创建技能失败')
    }
  }
)

const skillsSlice = createSlice({
  name: 'skills',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchSkills.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchSkills.fulfilled, (state, action) => {
        state.loading = false
        state.skills = action.payload
      })
      .addCase(fetchSkills.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '获取技能列表失败'
      })
      .addCase(createSkill.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(createSkill.fulfilled, (state) => {
        state.loading = false
      })
      .addCase(createSkill.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string || '创建技能失败'
      })
  },
})

export default skillsSlice.reducer