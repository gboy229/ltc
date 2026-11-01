-- 收益流水：静态、直推、本金等统一记录
CREATE TABLE IF NOT EXISTS `hg_addon_pledge_income_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '收益用户',
  `order_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '来源报单',
  `from_user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '来源用户（直推/团队）',
  `type` varchar(32) NOT NULL COMMENT 'static/direct/indirect/team/principal',
  `amount` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '收益金额',
  `base_amount` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '计算基数',
  `rate` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '收益率',
  `period_no` int NOT NULL DEFAULT '0' COMMENT '第几期',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '1=已入账 2=已撤销',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_income` (`user_id`,`type`,`order_id`,`period_no`,`from_user_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='质押收益流水表';
