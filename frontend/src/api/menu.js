import request from '@/utils/request'

export const menuApi = {
  list: () => request.get('/menus'),
  get: (id) => request.get(`/menus/${id}`),
  create: (data) => request.post('/menus', data),
  update: (id, data) => request.put(`/menus/${id}`, data),
  delete: (id) => request.delete(`/menus/${id}`),
}
