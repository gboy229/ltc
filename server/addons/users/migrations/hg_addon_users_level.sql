-- 用户等级表
CREATE TABLE IF NOT EXISTS `hg_addon_users_level` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '等级名称，例如 LV1、新手',
  `personal_amount` decimal(20,2) NOT NULL DEFAULT '0.00' COMMENT '个人业绩',
  `team_amount` decimal(20,2) NOT NULL DEFAULT '0.00' COMMENT '团队业绩',
  `diff_reward` decimal(20,2) NOT NULL DEFAULT '0.00' COMMENT '级差奖励',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注说明',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1=启用，2=禁用',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序值，越大越靠前',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户等级表';
