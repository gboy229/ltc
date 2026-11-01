-- hg_admin_notice
CREATE TABLE IF NOT EXISTS `hg_admin_notice` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '公告ID',
  `title` varchar(64) NOT NULL COMMENT '公告标题',
  `type` bigint(20) NOT NULL COMMENT '公告类型',
  `tag` int(11) DEFAULT NULL COMMENT '标签',
  `content` longtext NOT NULL COMMENT '公告内容',
  `receiver` json DEFAULT NULL COMMENT '接收者',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '公告状态',
  `created_by` bigint(20) DEFAULT NULL COMMENT '发送人',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '修改人',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_通知公告';
