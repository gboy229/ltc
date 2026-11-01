import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/cms/article/list',
    method: 'get',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/cms/article/view',
    method: 'get',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/cms/article/edit',
    method: 'post',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/cms/article/delete',
    method: 'post',
    params,
  });
}

export function Status(params) {
  return http.request({
    url: '/cms/article/status',
    method: 'post',
    params,
  });
}

// 获取最大排序
export function MaxSort() {
  return http.request({
    url: '/cms/article/maxSort',
    method: 'get',
  });
}
