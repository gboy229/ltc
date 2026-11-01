import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/users/user/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/users/user/view',
    method: 'get',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/users/user/edit',
    method: 'post',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/users/user/delete',
    method: 'post',
    params,
  });
}

export function Status(params) {
  return http.request({
    url: '/users/user/status',
    method: 'post',
    params,
  });
}
