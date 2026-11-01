-- 质押投资套餐表
CREATE TABLE IF NOT EXISTS `hg_addon_pledge_plan` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '套餐名称，如 V1',
  `amount` decimal(20,2) NOT NULL DEFAULT '0.00' COMMENT '入金金额',
  `period_months` int NOT NULL COMMENT '投资周期（月）',
  `daily_lots` int NOT NULL DEFAULT '0' COMMENT '每日交易手数',
  `reward_per_lot` decimal(20,2) NOT NULL DEFAULT '0.00' COMMENT '每手金额（USDT）',
  `early_exit_penalty_rate` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '提前退本违约金比例（%）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1=启用，2=禁用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name_period` (`name`,`period_months`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='质押投资套餐表';
