-- 登录凭证表（一行一种登录方式）
CREATE TABLE IF NOT EXISTS `hg_addon_users_user_account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '凭证ID',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `type` varchar(20) NOT NULL DEFAULT '' COMMENT '凭证类型:username,mobile,email,wallet,wechat,google...',
  `account` varchar(128) NOT NULL DEFAULT '' COMMENT '凭证值(手机号/邮箱/用户名/钱包地址/openid)',
  `credential` varchar(255) NOT NULL DEFAULT '' COMMENT '凭证附加信息(可选,如钱包签名参数等)',
  `verified` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已验证(手机/邮箱/钱包通常需验证后才可登录)',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_type_account` (`type`,`account`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录凭证表';
