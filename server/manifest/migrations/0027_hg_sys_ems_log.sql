-- hg_sys_ems_log
CREATE TABLE IF NOT EXISTS `hg_sys_ems_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `event` varchar(64) NOT NULL COMMENT '事件',
  `email` varchar(512) NOT NULL COMMENT '邮箱地址，多个用;隔开',
  `code` varchar(256) DEFAULT '' COMMENT '验证码',
  `times` bigint(20) NOT NULL COMMENT '验证次数',
  `content` longtext COMMENT '邮件内容',
  `ip` varchar(128) DEFAULT NULL COMMENT 'ip地址',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态(1未验证,2已验证)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='系统_邮件发送记录';
