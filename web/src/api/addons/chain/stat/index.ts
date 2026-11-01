import { http } from '@/utils/http/axios';

export function Stat(params) {
  return http.request({
    url: '/chain/stat',
    method: 'get',
    params,
  });
}
