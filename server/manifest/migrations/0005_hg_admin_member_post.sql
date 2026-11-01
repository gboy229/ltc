-- hg_admin_member_post
CREATE TABLE IF NOT EXISTS `hg_admin_member_post` (
  `member_id` bigint(20) NOT NULL COMMENT '管理员ID',
  `post_id` bigint(20) NOT NULL COMMENT '岗位ID',
  PRIMARY KEY (`member_id`,`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员_用户岗位关联';
