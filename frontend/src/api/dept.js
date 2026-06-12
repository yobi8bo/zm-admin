import request from '@/utils/request'

export const deptApi = {
  list: () => request.get('/depts'),
  get: (id) => request.get(`/depts/${id}`),
  create: (data) => request.post('/depts', data),
  update: (id, data) => request.put(`/depts/${id}`, data),
  delete: (id) => request.delete(`/depts/${id}`),
}
