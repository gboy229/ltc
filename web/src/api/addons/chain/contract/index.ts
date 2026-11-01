import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/chain/contract/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/chain/contract/view',
    method: 'get',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/chain/contract/edit',
    method: 'post',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/chain/contract/delete',
    method: 'post',
    params,
  });
}
