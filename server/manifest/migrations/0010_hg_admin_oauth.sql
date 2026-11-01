-- hg_admin_oauth
CREATE TABLE IF NOT EXISTS `hg_admin_oauth` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `member_id` bigint(20) DEFAULT '0' COMMENT '用户ID',
  `unionid` varchar(64) DEFAULT '' COMMENT '唯一ID',
  `oauth_client` varchar(32) DEFAULT NULL COMMENT '授权组别',
  `oauth_openid` varchar(128) DEFAULT NULL COMMENT '授权开放ID',
  `sex` tinyint(1) DEFAULT '1' COMMENT '性别',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '昵称',
  `head_portrait` varchar(512) DEFAULT NULL COMMENT '头像',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `country` varchar(100) DEFAULT '' COMMENT '国家',
  `province` varchar(100) DEFAULT '' COMMENT '省',
  `city` varchar(100) DEFAULT '' COMMENT '市',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(1) DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `oauth_client` (`oauth_client`,`oauth_openid`),
  KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员_第三方登录';
