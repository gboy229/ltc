-- hg_sys_cron
CREATE TABLE IF NOT EXISTS `hg_sys_cron` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `group_id` bigint(20) NOT NULL COMMENT '分组ID',
  `title` varchar(128) NOT NULL COMMENT '任务标题',
  `name` varchar(100) DEFAULT NULL COMMENT '任务方法',
  `params` varchar(255) DEFAULT NULL COMMENT '函数参数',
  `pattern` varchar(64) NOT NULL COMMENT '表达式',
  `policy` bigint(20) NOT NULL DEFAULT '1' COMMENT '策略',
  `count` bigint(20) NOT NULL DEFAULT '0' COMMENT '执行次数',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '任务状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COMMENT='系统_定时任务';
