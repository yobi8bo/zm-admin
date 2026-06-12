import request from '@/utils/request'

export const userApi = {
  getMe: () => request.get('/user/me'),
  updateMe: (data) => request.put('/user/me', data),
  updateMyPassword: (data) => request.put('/user/me/password', data),
  getMyMenus: () => request.get('/user/me/menus'),
  getMyPermissions: () => request.get('/user/me/permissions'),

  list: (params) => request.get('/users', { params }),
  get: (id) => request.get(`/users/${id}`),
  create: (data) => request.post('/users', data),
  update: (id, data) => request.put(`/users/${id}`, data),
  delete: (id) => request.delete(`/users/${id}`),
  updateStatus: (id, data) => request.put(`/users/${id}/status`, data),
  resetPassword: (id, data) => request.put(`/users/${id}/password`, data),
  assignRoles: (id, data) => request.put(`/users/${id}/roles`, data),
}
