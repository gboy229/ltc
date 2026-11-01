import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/chain/withdraw/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/chain/withdraw/view',
    method: 'get',
    params,
  });
}

export function Status(params) {
  return http.request({
    url: '/chain/withdraw/status',
    method: 'post',
    params,
  });
}
