import { request } from '@/utils'

export default {
  toggleRole: (data) => request.post('/auth/role/toggle', data),
  login: (data) => request.post('/auth/login', data, { noNeedToken: true }),
  getUser: () => request.get('/user/detail'),
}
