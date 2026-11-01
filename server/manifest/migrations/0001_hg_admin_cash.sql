-- hg_admin_cash
CREATE TABLE IF NOT EXISTS `hg_admin_cash` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `member_id` bigint(20) NOT NULL COMMENT '管理员ID',
  `money` decimal(10,2) NOT NULL COMMENT '提现金额',
  `fee` decimal(10,2) NOT NULL COMMENT '手续费',
  `last_money` decimal(10,2) NOT NULL COMMENT '最终到账金额',
  `ip` varchar(128) NOT NULL COMMENT '申请人IP',
  `status` bigint(20) NOT NULL COMMENT '状态码',
  `msg` varchar(128) DEFAULT NULL COMMENT '处理结果',
  `handle_at` datetime DEFAULT NULL COMMENT '处理时间',
  `created_at` datetime NOT NULL COMMENT '申请时间',
  PRIMARY KEY (`id`),
  KEY `admin_id` (`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_提现记录表';
