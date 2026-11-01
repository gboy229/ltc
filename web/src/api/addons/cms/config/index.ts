import { http } from '@/utils/http/axios';

export function getConfig(params) {
  return http.request({
    url: '/cms/config/get',
    method: 'get',
    params,
  });
}

export function updateConfig(params) {
  return http.request({
    url: '/cms/config/update',
    method: 'post',
    params,
  });
}
