-- 链配置表
CREATE TABLE IF NOT EXISTS `hg_addon_chain_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `driver` varchar(20) NOT NULL COMMENT '驱动',
  `chain` varchar(20) NOT NULL COMMENT '网络',
  `rpc_url` varchar(100) NOT NULL COMMENT 'rpc地址',
  `chain_id` int NOT NULL DEFAULT 0 COMMENT '链id',
  `browser_url` varchar(200) NOT NULL COMMENT '浏览器地址',
  `transfer_address` varchar(100) DEFAULT NULL COMMENT '转账地址',
  `transfer_private` varchar(200) DEFAULT NULL COMMENT '转账私钥',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `confirmations` int NOT NULL DEFAULT 15 COMMENT '确认数',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `chain_config_chain_unique` (`chain`) USING BTREE,
  UNIQUE KEY `chain_config_driver_chain_unique` (`driver`,`chain`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='链配置表';
