import { http } from '@/utils/http/axios';

export function getConfig(params) {
  return http.request({
    url: '/users/config/get',
    method: 'get',
    params,
  });
}

export function updateConfig(params) {
  return http.request({
    url: '/users/config/update',
    method: 'post',
    params,
  });
}

export function getConfigItems() {
  return http.request({
    url: '/users/config/items',
    method: 'get',
  });
}
