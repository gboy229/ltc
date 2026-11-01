-- 用户业绩汇总表
CREATE TABLE IF NOT EXISTS `hg_addon_users_performance` (
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `personal_amount` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '个人业绩',
  `direct_amount` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '直推业绩',
  `team_amount` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '团队业绩',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户业绩汇总';
