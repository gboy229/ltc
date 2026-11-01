-- hg_sys_gen_codes
CREATE TABLE IF NOT EXISTS `hg_sys_gen_codes` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '生成ID',
  `gen_type` int(10) unsigned NOT NULL COMMENT '生成类型',
  `gen_template` int(11) DEFAULT '0' COMMENT '生成模板',
  `var_name` varchar(255) NOT NULL COMMENT '实体命名',
  `options` json DEFAULT NULL COMMENT '配置选项',
  `db_name` varchar(128) DEFAULT NULL COMMENT '数据库名称',
  `table_name` varchar(255) NOT NULL COMMENT '主表名称',
  `table_comment` varchar(255) DEFAULT NULL COMMENT '主表注释',
  `dao_name` varchar(255) DEFAULT NULL COMMENT '主表dao模型',
  `master_columns` json DEFAULT NULL COMMENT '主表字段',
  `addon_name` varchar(128) DEFAULT NULL COMMENT '插件名称',
  `status` tinyint(1) DEFAULT '1' COMMENT '生成状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COMMENT='系统_代码生成记录';
