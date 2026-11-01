import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/chain/recharge/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/chain/recharge/view',
    method: 'get',
    params,
  });
}
