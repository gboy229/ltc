-- hg_sys_login_log
CREATE TABLE IF NOT EXISTS `hg_sys_login_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `req_id` varchar(50) DEFAULT NULL COMMENT '请求ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '用户ID',
  `username` varchar(64) DEFAULT NULL COMMENT '用户名',
  `response` json DEFAULT NULL COMMENT '响应数据',
  `login_at` datetime DEFAULT NULL COMMENT '登录时间',
  `login_ip` varchar(128) DEFAULT NULL COMMENT '登录IP',
  `province_id` bigint(20) DEFAULT NULL COMMENT '省编码',
  `city_id` bigint(20) DEFAULT NULL COMMENT '市编码',
  `user_agent` varchar(512) DEFAULT NULL COMMENT 'UA信息',
  `err_msg` varchar(1000) DEFAULT NULL COMMENT '错误提示',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `member_id` (`member_id`),
  KEY `req_id` (`req_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_登录日志';
