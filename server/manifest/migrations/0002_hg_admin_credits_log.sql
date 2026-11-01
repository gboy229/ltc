-- hg_admin_credits_log
CREATE TABLE IF NOT EXISTS `hg_admin_credits_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '变动ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '管理员ID',
  `app_id` varchar(64) DEFAULT NULL COMMENT '应用id',
  `addons_name` varchar(100) NOT NULL DEFAULT '' COMMENT '插件名称',
  `credit_type` varchar(32) NOT NULL DEFAULT '' COMMENT '变动类型',
  `credit_group` varchar(32) DEFAULT NULL COMMENT '变动组别',
  `before_num` decimal(10,2) DEFAULT '0.00' COMMENT '变动前',
  `num` decimal(10,2) DEFAULT '0.00' COMMENT '变动数据',
  `after_num` decimal(10,2) DEFAULT '0.00' COMMENT '变动后',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `ip` varchar(20) DEFAULT NULL COMMENT '操作人IP',
  `map_id` bigint(20) DEFAULT '0' COMMENT '关联ID',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `member_id` (`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_资产变动表';
