-- 内容文章表（当前仅公告）
CREATE TABLE IF NOT EXISTS `hg_addon_cms_article` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `article_no` varchar(32) NOT NULL DEFAULT '' COMMENT '前台公开ID',
  `title` varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '类型：1=公告',
  `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '封面图',
  `content` longtext COMMENT '内容',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序，越大越靠前',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1=正常，2=禁用',
  `is_recommend` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否推荐：1=推荐，0=不推荐',
  `created_by` bigint(20) DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint(20) DEFAULT NULL COMMENT '修改人',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_article_no` (`article_no`),
  KEY `idx_type_status` (`type`,`status`),
  KEY `idx_recommend` (`is_recommend`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容文章表';
