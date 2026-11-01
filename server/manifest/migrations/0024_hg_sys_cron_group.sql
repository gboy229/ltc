-- hg_sys_cron_group
CREATE TABLE IF NOT EXISTS `hg_sys_cron_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '任务分组ID',
  `pid` bigint(20) NOT NULL COMMENT '父类任务分组ID',
  `name` varchar(100) DEFAULT '' COMMENT '分组名称',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '分组状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='系统_定时任务分组';
