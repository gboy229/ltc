-- hg_sys_provinces
CREATE TABLE IF NOT EXISTS `hg_sys_provinces` (
  `id` bigint(20) NOT NULL COMMENT '省市区ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '栏目名称',
  `pinyin` varchar(100) DEFAULT '' COMMENT '拼音',
  `lng` varchar(20) DEFAULT '' COMMENT '经度',
  `lat` varchar(20) DEFAULT '' COMMENT '纬度',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '父栏目',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '关系树等级',
  `tree` varchar(200) DEFAULT NULL COMMENT '关系树',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `pid` (`pid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='系统_省市区编码';
