import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({ url: '/pledge/periodConfig/list', method: 'get', params });
}

export function View(params) {
  return http.request({ url: '/pledge/periodConfig/view', method: 'get', params });
}

export function Edit(params) {
  return http.request({ url: '/pledge/periodConfig/edit', method: 'post', params });
}

export function Delete(params) {
  return http.request({ url: '/pledge/periodConfig/delete', method: 'post', params });
}

export function Status(params) {
  return http.request({ url: '/pledge/periodConfig/status', method: 'post', params });
}
