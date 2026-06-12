import request from '@/utils/request'

export const logApi = {
  operationList: (params) => request.get('/logs/operation', { params }),
  clearOperationLog: () => request.delete('/logs/operation'),
  loginList: (params) => request.get('/logs/login', { params }),
  clearLoginLog: () => request.delete('/logs/login'),
}
