-- hg_sys_dict_data
CREATE TABLE IF NOT EXISTS `hg_sys_dict_data` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `label` varchar(100) DEFAULT NULL COMMENT '字典标签',
  `value` varchar(100) DEFAULT NULL COMMENT '字典键值',
  `value_type` varchar(255) NOT NULL DEFAULT 'string' COMMENT '键值数据类型：string,int,uint,bool,datetime,date',
  `type` varchar(100) DEFAULT NULL COMMENT '字典类型',
  `list_class` varchar(100) DEFAULT NULL COMMENT '表格回显样式',
  `is_default` tinyint(1) DEFAULT '2' COMMENT '是否为系统默认',
  `sort` int(11) DEFAULT '0' COMMENT '字典排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `dict_data_idx` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=171 DEFAULT CHARSET=utf8mb4 COMMENT='系统_字典数据';
