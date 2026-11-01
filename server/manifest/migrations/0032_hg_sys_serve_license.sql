-- hg_sys_serve_license
CREATE TABLE IF NOT EXISTS `hg_sys_serve_license` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '许可ID',
  `group` varchar(50) NOT NULL COMMENT '分组',
  `name` varchar(128) NOT NULL COMMENT '许可名称',
  `appid` varchar(64) NOT NULL COMMENT '应用ID',
  `secret_key` varchar(255) DEFAULT NULL COMMENT '应用秘钥',
  `remote_addr` varchar(64) DEFAULT NULL COMMENT '最后连接地址',
  `online_limit` int(11) DEFAULT '1' COMMENT '在线限制',
  `login_times` bigint(20) DEFAULT NULL COMMENT '登录次数',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_active_at` datetime DEFAULT NULL COMMENT '最后心跳',
  `routes` json DEFAULT NULL COMMENT '路由表，空使用默认分组路由',
  `allowed_ips` varchar(512) DEFAULT NULL COMMENT 'IP白名单',
  `end_at` datetime NOT NULL COMMENT '授权有效期',
  `remark` varchar(512) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `appid` (`appid`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='系统_服务许可证';
