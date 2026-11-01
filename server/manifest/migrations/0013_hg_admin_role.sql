-- hg_admin_role
CREATE TABLE IF NOT EXISTS `hg_admin_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(32) NOT NULL COMMENT '角色名称',
  `key` varchar(128) NOT NULL COMMENT '角色权限字符串',
  `data_scope` tinyint(1) DEFAULT '1' COMMENT '数据范围',
  `custom_dept` json DEFAULT NULL COMMENT '自定义部门权限',
  `pid` bigint(20) DEFAULT '0' COMMENT '上级角色ID',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '关系树等级',
  `tree` varchar(512) DEFAULT NULL COMMENT '关系树',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '角色状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=211 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_角色信息';
