import { http } from '@/utils/http/axios';

// 资产账户列表
export function AccountList(params) {
  return http.request({
    url: '/users/asset/accountList',
    method: 'get',
    params,
  });
}

// 资金流水列表
export function LogList(params) {
  return http.request({
    url: '/users/asset/logList',
    method: 'get',
    params,
  });
}

// 资产统计
export function Stat(params) {
  return http.request({
    url: '/users/asset/stat',
    method: 'get',
    params,
  });
}
