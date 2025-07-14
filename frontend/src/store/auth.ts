import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface User {
  id: number
  email: string
  name: string
  role: 'admin' | 'general'
  created_at: string
  updated_at: string
}

export interface Attendance {
  id: number
  user_id: number
  date: string
  clock_in_time?: string
  clock_out_time?: string
  work_hours?: number
  created_at: string
  updated_at: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (token: string, user: User) => void
  logout: () => void
  setUser: (user: User) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (token: string, user: User) =>
        set({ token, user, isAuthenticated: true }),
      logout: () =>
        set({ token: null, user: null, isAuthenticated: false }),
      setUser: (user: User) => set({ user }),
    }),
    {
      name: 'auth-storage',
    }
  )
)

interface AttendanceState {
  todayAttendance: Attendance | null
  attendanceHistory: Attendance[]
  isLoading: boolean
  setTodayAttendance: (attendance: Attendance | null) => void
  setAttendanceHistory: (history: Attendance[]) => void
  setLoading: (loading: boolean) => void
}

export const useAttendanceStore = create<AttendanceState>((set) => ({
  todayAttendance: null,
  attendanceHistory: [],
  isLoading: false,
  setTodayAttendance: (attendance) => set({ todayAttendance: attendance }),
  setAttendanceHistory: (history) => set({ attendanceHistory: history }),
  setLoading: (loading) => set({ isLoading: loading }),
}))