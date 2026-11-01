-- hg_admin_role_casbin
CREATE TABLE IF NOT EXISTS `hg_admin_role_casbin` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `p_type` varchar(64) DEFAULT NULL,
  `v0` varchar(256) DEFAULT NULL,
  `v1` varchar(256) DEFAULT NULL,
  `v2` varchar(256) DEFAULT NULL,
  `v3` varchar(256) DEFAULT NULL,
  `v4` varchar(256) DEFAULT NULL,
  `v5` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='管理员_casbin权限表';
