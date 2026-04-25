import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

export interface ApiError {
  message: string
  code?: string
  status?: number
}

export interface ApiResponse<T = unknown> {
  data: T
  message?: string
}

class ApiService {
  private instance: AxiosInstance

  constructor() {
    this.instance = axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors(): void {
    this.instance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        const sessionId = localStorage.getItem('sessionId')
        if (sessionId && config.headers) {
          config.headers['Authorization'] = sessionId
        }
        return config
      },
      (error: AxiosError) => {
        return Promise.reject(this.handleError(error))
      }
    )

    this.instance.interceptors.response.use(
      (response: AxiosResponse) => {
        return response
      },
      (error: AxiosError) => {
        // Handle 401 Unauthorized - redirect to login
        if (error.response?.status === 401) {
          // Clear session data
          localStorage.removeItem('sessionId')
          // Redirect to login page
          window.location.href = '/login'
        }
        return Promise.reject(this.handleError(error))
      }
    )
  }

  private handleError(error: AxiosError): ApiError {
    if (error.response) {
      const { status, data } = error.response as AxiosResponse & { data: { message?: string } }
      switch (status) {
        case 400:
          return { message: data?.message || '请求参数错误', status, code: 'BAD_REQUEST' }
        case 401:
          return { message: '未授权，请重新登录', status, code: 'UNAUTHORIZED' }
        case 403:
          return { message: '拒绝访问', status, code: 'FORBIDDEN' }
        case 404:
          return { message: '请求的资源不存在', status, code: 'NOT_FOUND' }
        case 500:
          return { message: '服务器内部错误', status, code: 'INTERNAL_ERROR' }
        default:
          return { message: data?.message || '请求失败', status, code: 'UNKNOWN' }
      }
    } else if (error.request) {
      return { message: '网络连接失败，请检查网络', code: 'NETWORK_ERROR' }
    } else {
      return { message: error.message || '请求配置错误', code: 'CONFIG_ERROR' }
    }
  }

  public get<T = unknown>(url: string, params?: Record<string, unknown>): Promise<AxiosResponse<T>> {
    return this.instance.get<T>(url, { params })
  }

  public post<T = unknown>(url: string, data?: Record<string, unknown>): Promise<AxiosResponse<T>> {
    return this.instance.post<T>(url, data)
  }

  public put<T = unknown>(url: string, data?: Record<string, unknown>): Promise<AxiosResponse<T>> {
    return this.instance.put<T>(url, data)
  }

  public delete<T = unknown>(url: string): Promise<AxiosResponse<T>> {
    return this.instance.delete<T>(url)
  }
}

export const apiService = new ApiService()
export default apiService