-- hg_sys_addons_config
CREATE TABLE IF NOT EXISTS `hg_sys_addons_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `addon_name` varchar(128) NOT NULL COMMENT '插件名称',
  `group` varchar(128) NOT NULL COMMENT '分组',
  `name` varchar(100) DEFAULT '' COMMENT '参数名称',
  `type` varchar(32) NOT NULL COMMENT '键值类型:string,int,uint,bool,datetime,date',
  `key` varchar(100) DEFAULT '' COMMENT '参数键名',
  `value` varchar(500) DEFAULT '' COMMENT '参数键值',
  `default_value` varchar(500) NOT NULL COMMENT '默认值',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `tip` varchar(500) DEFAULT NULL COMMENT '变量描述',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否为系统默认',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `addon_name_2` (`addon_name`,`key`),
  KEY `addon_name` (`addon_name`),
  KEY `addon_name_3` (`addon_name`,`group`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='系统_插件配置';
