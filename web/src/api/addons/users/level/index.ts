import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/users/level/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/users/level/view',
    method: 'get',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/users/level/edit',
    method: 'post',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/users/level/delete',
    method: 'post',
    params,
  });
}

export function Status(params) {
  return http.request({
    url: '/users/level/status',
    method: 'post',
    params,
  });
}

// 获取最大排序
export function MaxSort() {
  return http.request({
    url: '/users/level/maxSort',
    method: 'get',
  });
}

// 获取等级选项
export function Option() {
  return http.request({
    url: '/users/level/option',
    method: 'get',
  });
}
