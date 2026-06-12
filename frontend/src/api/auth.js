import request from '@/utils/request'

export const authApi = {
  getCaptcha: () => request.get('/auth/captcha'),

  login: (data) => request.post('/auth/login', data),

  logout: () => request.post('/auth/logout'),

  refreshToken: (data) => request.post('/auth/refresh', data),
}
