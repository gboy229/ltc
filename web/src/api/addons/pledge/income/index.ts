import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({ url: '/pledge/income/list', method: 'get', params });
}
