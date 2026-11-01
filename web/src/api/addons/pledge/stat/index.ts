import { http } from '@/utils/http/axios';

export function Stat(params) {
  return http.request({
    url: '/pledge/stat',
    method: 'get',
    params,
  });
}
