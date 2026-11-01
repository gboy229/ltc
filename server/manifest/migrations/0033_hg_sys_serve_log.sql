-- hg_sys_serve_log
CREATE TABLE IF NOT EXISTS `hg_sys_serve_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `trace_id` varchar(50) DEFAULT NULL COMMENT '链路ID',
  `level_format` varchar(32) DEFAULT NULL COMMENT '日志级别',
  `content` text COMMENT '日志内容',
  `stack` json DEFAULT NULL COMMENT '打印堆栈',
  `line` varchar(255) NOT NULL COMMENT '调用行',
  `trigger_ns` bigint(20) DEFAULT NULL COMMENT '触发时间(ns)',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `member_id` (`level_format`),
  KEY `traceid` (`trace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_服务日志';
