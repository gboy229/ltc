-- 用户主体表
CREATE TABLE IF NOT EXISTS `hg_addon_users_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `uid` varchar(32) NOT NULL DEFAULT '' COMMENT '对外唯一用户号',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '上级用户ID(邀请关系)',
  `level_id` int(11) NOT NULL DEFAULT '0' COMMENT '会员等级ID',
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `password` varchar(60) NOT NULL DEFAULT '' COMMENT '登录密码(bcrypt哈希,自带盐)',
  `trade_password` varchar(60) NOT NULL DEFAULT '' COMMENT '交易密码(bcrypt哈希,自带盐)',
  `authenticator_key` varchar(100) NOT NULL DEFAULT '' COMMENT 'TOTP验证器密钥',
  `authenticator_enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否开启验证器',
  `auth_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '实名认证状态',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `extend` json DEFAULT NULL COMMENT '扩展资料(零散资料兜底,如性别/生日/地区等)',
  `last_login_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
