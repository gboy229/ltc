import { http } from '@/utils/http/axios';

export function getConfig(params) {
  return http.request({
    url: '/pledge/config/get',
    method: 'get',
    params,
  });
}

export function updateConfig(params) {
  return http.request({
    url: '/pledge/config/update',
    method: 'post',
    params,
  });
}

export function getConfigItems() {
  return http.request({
    url: '/pledge/config/items',
    method: 'get',
  });
}
