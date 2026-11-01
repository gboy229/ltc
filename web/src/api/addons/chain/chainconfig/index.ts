import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/chain/chainconfig/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/chain/chainconfig/view',
    method: 'get',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/chain/chainconfig/edit',
    method: 'post',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/chain/chainconfig/delete',
    method: 'post',
    params,
  });
}

export function GenerateAddress(params) {
  return http.request({
    url: '/chain/chainconfig/generateAddress',
    method: 'post',
    params,
  });
}

export function Options() {
  return http.request({
    url: '/chain/chainconfig/options',
    method: 'get',
  });
}
