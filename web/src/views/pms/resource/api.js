import axios from 'axios'
import { request } from '@/utils'

export default {
  getMenuTree: () => request.get('/permission/menu/tree'),
  getButtonAndApi: (parentId) => request.get(`/permission/button-and-api/${parentId}`),
  getComponents: () => axios.get(`${import.meta.env.VITE_PUBLIC_PATH}components.json`),
  addPermission: (data) => request.post('/permission', data),
  savePermission: (id, data) => request.patch(`/permission/${id}`, data),
  deletePermission: (id) => request.delete(`permission/${id}`),
}
