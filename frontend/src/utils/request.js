import axios from 'axios'
import { message } from 'ant-design-vue'
import { getToken, setToken, getRefreshToken, clearAuth } from './auth'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE,
  timeout: 15000,
})

// 是否正在刷新 token
let isRefreshing = false
// 等待刷新的请求队列
let requestQueue = []

const processQueue = (error, token = null) => {
  requestQueue.forEach((p) => (error ? p.reject(error) : p.resolve(token)))
  requestQueue = []
}

// 请求拦截器
request.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const { code, msg, data } = response.data

    if (code === 200) return data

    // token 过期，尝试刷新
    if (code === 10041) {
      return handleTokenExpired(response.config)
    }

    // token 无效或未登录，跳转登录
    if (code === 401 || code === 10042) {
      redirectToLogin()
      return Promise.reject(new Error(msg))
    }

    message.error(msg || '请求失败')
    return Promise.reject(new Error(msg))
  },
  (error) => {
    message.error(error.message || '网络错误，请稍后重试')
    return Promise.reject(error)
  }
)

async function handleTokenExpired(originalConfig) {
  if (isRefreshing) {
    return new Promise((resolve, reject) => {
      requestQueue.push({ resolve, reject })
    }).then((token) => {
      originalConfig.headers.Authorization = `Bearer ${token}`
      return request(originalConfig)
    })
  }

  isRefreshing = true
  const refreshToken = getRefreshToken()

  if (!refreshToken) {
    redirectToLogin()
    return Promise.reject(new Error('无效的刷新令牌'))
  }

  try {
    const res = await axios.post(`${import.meta.env.VITE_API_BASE}/auth/refresh`, {
      refresh_token: refreshToken,
    })
    const { access_token } = res.data.data
    setToken(access_token)
    processQueue(null, access_token)
    originalConfig.headers.Authorization = `Bearer ${access_token}`
    return request(originalConfig)
  } catch {
    processQueue(new Error('刷新token失败'))
    clearAuth()
    redirectToLogin()
    return Promise.reject(new Error('刷新token失败'))
  } finally {
    isRefreshing = false
  }
}

function redirectToLogin() {
  clearAuth()
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

export default request
