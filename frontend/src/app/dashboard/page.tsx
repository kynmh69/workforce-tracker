'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useAuthStore, useAttendanceStore } from '@/store/auth'
import { api } from '@/lib/api'
import { Clock, LogOut, Users, History } from 'lucide-react'

export default function DashboardPage() {
  const [currentTime, setCurrentTime] = useState(new Date())
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState('')
  const router = useRouter()
  
  const { user, isAuthenticated, logout } = useAuthStore()
  const { todayAttendance, setTodayAttendance } = useAttendanceStore()

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }

    // Update time every second
    const timer = setInterval(() => setCurrentTime(new Date()), 1000)
    
    // Load today's attendance
    const loadTodayAttendance = async () => {
      try {
        const response = await api.get('/attendance/today')
        setTodayAttendance(response.data)
      } catch (error) {
        console.error('Failed to load today attendance:', error)
      }
    }

    loadTodayAttendance()

    return () => clearInterval(timer)
  }, [isAuthenticated, router, setTodayAttendance])

  const refreshTodayAttendance = async () => {
    try {
      const response = await api.get('/attendance/today')
      setTodayAttendance(response.data)
    } catch (error) {
      console.error('Failed to load today attendance:', error)
    }
  }

  const handleClockIn = async () => {
    setIsLoading(true)
    setMessage('')

    try {
      await api.post('/attendance/clock-in')
      setMessage('出勤打刻が完了しました')
      refreshTodayAttendance()
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string } } }
      setMessage(error.response?.data?.message || '出勤打刻に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  const handleClockOut = async () => {
    setIsLoading(true)
    setMessage('')

    try {
      await api.post('/attendance/clock-out')
      setMessage('退勤打刻が完了しました')
      refreshTodayAttendance()
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string } } }
      setMessage(error.response?.data?.message || '退勤打刻に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  const formatTime = (date: Date) => {
    return date.toLocaleString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  const formatWorkHours = (hours?: number) => {
    if (!hours) return '0時間0分'
    const h = Math.floor(hours)
    const m = Math.floor((hours - h) * 60)
    return `${h}時間${m}分`
  }

  if (!isAuthenticated || !user) {
    return <div>Loading...</div>
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <h1 className="text-xl font-semibold text-gray-900">勤怠管理システム</h1>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-700">
                {user.name}さん ({user.role === 'admin' ? '管理者' : '一般ユーザー'})
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={handleLogout}
                className="flex items-center space-x-2"
              >
                <LogOut className="h-4 w-4" />
                <span>ログアウト</span>
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Clock and Time Display */}
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Clock className="h-5 w-5" />
                  <span>打刻</span>
                </CardTitle>
                <CardDescription>
                  出勤・退勤の打刻を行ってください
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-center space-y-6">
                  <div className="text-3xl font-mono font-bold text-gray-900">
                    {formatTime(currentTime)}
                  </div>
                  
                  <div className="flex justify-center space-x-4">
                    <Button
                      onClick={handleClockIn}
                      disabled={isLoading || !!todayAttendance?.clock_in_time}
                      className="w-32"
                    >
                      出勤
                    </Button>
                    <Button
                      onClick={handleClockOut}
                      disabled={isLoading || !todayAttendance?.clock_in_time || !!todayAttendance?.clock_out_time}
                      variant="secondary"
                      className="w-32"
                    >
                      退勤
                    </Button>
                  </div>

                  {message && (
                    <div className={`p-3 rounded ${
                      message.includes('失敗') 
                        ? 'bg-red-50 text-red-600 border border-red-200'
                        : 'bg-green-50 text-green-600 border border-green-200'
                    }`}>
                      {message}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Today's Attendance Status */}
          <div>
            <Card>
              <CardHeader>
                <CardTitle>本日の勤怠状況</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm text-gray-600">出勤時刻:</span>
                    <span className="text-sm font-medium">
                      {todayAttendance?.clock_in_time 
                        ? new Date(todayAttendance.clock_in_time).toLocaleTimeString('ja-JP', {
                            hour: '2-digit',
                            minute: '2-digit',
                          })
                        : '未打刻'
                      }
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-gray-600">退勤時刻:</span>
                    <span className="text-sm font-medium">
                      {todayAttendance?.clock_out_time 
                        ? new Date(todayAttendance.clock_out_time).toLocaleTimeString('ja-JP', {
                            hour: '2-digit',
                            minute: '2-digit',
                          })
                        : '未打刻'
                      }
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-gray-600">勤務時間:</span>
                    <span className="text-sm font-medium">
                      {formatWorkHours(todayAttendance?.work_hours)}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Navigation Menu */}
        <div className="mt-8">
          <Card>
            <CardHeader>
              <CardTitle>メニュー</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Button
                  variant="outline"
                  className="h-20 flex flex-col items-center justify-center space-y-2"
                  onClick={() => router.push('/attendance/history')}
                >
                  <History className="h-6 w-6" />
                  <span>勤怠履歴</span>
                </Button>
                
                {user.role === 'admin' && (
                  <>
                    <Button
                      variant="outline"
                      className="h-20 flex flex-col items-center justify-center space-y-2"
                      onClick={() => router.push('/users')}
                    >
                      <Users className="h-6 w-6" />
                      <span>ユーザー管理</span>
                    </Button>
                    
                    <Button
                      variant="outline"
                      className="h-20 flex flex-col items-center justify-center space-y-2"
                      onClick={() => router.push('/admin/attendance')}
                    >
                      <History className="h-6 w-6" />
                      <span>全ユーザー勤怠</span>
                    </Button>
                  </>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}