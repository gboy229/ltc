-- 用户等级变更记录
CREATE TABLE IF NOT EXISTS `hg_addon_users_level_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `from_level_id` int NOT NULL DEFAULT '0' COMMENT '变更前等级ID',
  `to_level_id` int NOT NULL DEFAULT '0' COMMENT '变更后等级ID',
  `personal_amount` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '变更时个人业绩',
  `team_amount` decimal(20,4) NOT NULL DEFAULT '0.0000' COMMENT '变更时团队业绩',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户等级变更记录';
