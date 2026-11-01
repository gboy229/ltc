-- hg_sys_dict_type
CREATE TABLE IF NOT EXISTS `hg_sys_dict_type` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '父类字典类型ID',
  `name` varchar(100) DEFAULT '' COMMENT '字典类型名称',
  `type` varchar(100) DEFAULT '' COMMENT '字典类型',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '字典类型状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `dict_type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8mb4 COMMENT='系统_字典类型';
