-- hg_admin_dept
CREATE TABLE IF NOT EXISTS `hg_admin_dept` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '部门ID',
  `pid` bigint(20) DEFAULT '0' COMMENT '父部门ID',
  `name` varchar(32) DEFAULT NULL COMMENT '部门名称',
  `code` varchar(255) DEFAULT NULL COMMENT '部门编码',
  `type` varchar(10) DEFAULT NULL COMMENT '部门类型',
  `leader` varchar(32) DEFAULT NULL COMMENT '负责人',
  `phone` varchar(11) DEFAULT NULL COMMENT '联系电话',
  `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
  `level` int(11) NOT NULL COMMENT '关系树等级',
  `tree` varchar(512) DEFAULT NULL COMMENT '关系树',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '部门状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `pid` (`pid`)
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_部门';
