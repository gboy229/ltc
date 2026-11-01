import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/pledge/plan/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/pledge/plan/view',
    method: 'get',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/pledge/plan/edit',
    method: 'post',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/pledge/plan/delete',
    method: 'post',
    params,
  });
}

export function Status(params) {
  return http.request({
    url: '/pledge/plan/status',
    method: 'post',
    params,
  });
}
