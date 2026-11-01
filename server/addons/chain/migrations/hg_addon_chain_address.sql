-- 链地址表
CREATE TABLE IF NOT EXISTS `hg_addon_chain_address` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain` varchar(30) NOT NULL DEFAULT '0' COMMENT '链名',
  `uid` int NOT NULL DEFAULT 0 COMMENT '用户ID',
  `address` varchar(100) NOT NULL COMMENT '地址',
  `keys` varchar(255) NOT NULL COMMENT '私钥',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `chain_address_address_unique` (`address`) USING BTREE,
  KEY `chain_address_uid_index` (`uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='链地址表';
