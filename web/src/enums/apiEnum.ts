export enum ApiEnum {
  // api前缀
  Prefix = '/api',

  // 基础
  SiteAccountLogin = '/site/accountLogin', // 账号登录
  SiteLoginConfig = '/site/loginConfig', // 登录配置
  SiteLogout = '/site/logout', // 注销
  SiteConfig = '/site/config', // 配置信息

  // 用户
  MemberInfo = '/member/info', // 登录用户信息

  // 角色
  RoleDynamic = '/role/dynamic', // 动态路由
}
