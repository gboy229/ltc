-- hg_admin_member_role
CREATE TABLE IF NOT EXISTS `hg_admin_member_role` (
  `member_id` bigint(20) NOT NULL COMMENT '管理员ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`member_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员_用户角色关联';
