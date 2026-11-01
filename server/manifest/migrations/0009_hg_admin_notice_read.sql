-- hg_admin_notice_read
CREATE TABLE IF NOT EXISTS `hg_admin_notice_read` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `notice_id` bigint(20) NOT NULL COMMENT '公告ID',
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `clicks` int(11) DEFAULT '1' COMMENT '已读次数',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '阅读时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `notice_id` (`notice_id`,`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_公告已读记录';
