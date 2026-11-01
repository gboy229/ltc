-- hg_sys_blacklist
CREATE TABLE IF NOT EXISTS `hg_sys_blacklist` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '黑名单ID',
  `ip` varchar(100) DEFAULT '' COMMENT 'IP地址',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='系统_访问黑名单';
