import request from '@/utils/request'

export const roleApi = {
  list: (params) => request.get('/roles', { params }),
  all: () => request.get('/roles/all'),
  get: (id) => request.get(`/roles/${id}`),
  create: (data) => request.post('/roles', data),
  update: (id, data) => request.put(`/roles/${id}`, data),
  delete: (id) => request.delete(`/roles/${id}`),
  updateStatus: (id, data) => request.put(`/roles/${id}/status`, data),
  assignMenus: (id, data) => request.put(`/roles/${id}/menus`, data),
  getMenuIDs: (id) => request.get(`/roles/${id}/menus`),
}
