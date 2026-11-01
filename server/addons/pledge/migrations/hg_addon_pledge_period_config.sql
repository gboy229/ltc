-- 周期公共收益配置：同一投资周期共用一份当前配置
CREATE TABLE IF NOT EXISTS `hg_addon_pledge_period_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `period_months` int NOT NULL COMMENT '投资周期（月）',
  `monthly_reward_info` varchar(64) NOT NULL DEFAULT '' COMMENT '月化收益展示区间',
  `monthly_reward_rate` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '当前静态月化收益率（%）',
  `direct_reward_rate` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '当前直推收益率（%）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1=启用，2=禁用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_period_months` (`period_months`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='质押周期收益配置表';

INSERT IGNORE INTO `hg_addon_pledge_period_config`
(`period_months`, `monthly_reward_info`, `monthly_reward_rate`, `direct_reward_rate`, `status`)
VALUES
(3,  '8%～10%',  9.0000,  6.0000,  1),
(6,  '10%～12%', 11.0000, 8.0000,  1),
(12, '12%～15%', 13.5000, 10.0000, 1);
