-- hg_sys_attachment
CREATE TABLE IF NOT EXISTS `hg_sys_attachment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `app_id` varchar(64) NOT NULL COMMENT '应用ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '管理员ID',
  `cate_id` bigint(20) unsigned DEFAULT '0' COMMENT '上传分类',
  `drive` varchar(64) DEFAULT NULL COMMENT '上传驱动',
  `name` varchar(1000) DEFAULT NULL COMMENT '文件原始名',
  `kind` varchar(16) DEFAULT NULL COMMENT '上传类型',
  `mime_type` varchar(128) NOT NULL DEFAULT '' COMMENT '扩展类型',
  `naive_type` varchar(32) DEFAULT NULL COMMENT 'NaiveUI类型',
  `path` varchar(1000) DEFAULT NULL COMMENT '本地路径',
  `file_url` varchar(1000) DEFAULT NULL COMMENT 'url',
  `size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `ext` varchar(50) DEFAULT NULL COMMENT '扩展名',
  `md5` varchar(32) DEFAULT NULL COMMENT 'md5校验码',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `md5` (`md5`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='系统_附件管理';
