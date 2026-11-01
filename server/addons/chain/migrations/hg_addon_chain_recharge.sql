-- 链充值表
CREATE TABLE IF NOT EXISTS `hg_addon_chain_recharge` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(30) DEFAULT NULL,
  `uid` int NOT NULL COMMENT '用户ID',
  `chain` varchar(30) NOT NULL COMMENT '链',
  `symbol` varchar(20) NOT NULL COMMENT '币种',
  `from_address` varchar(100) DEFAULT NULL COMMENT '转出地址',
  `to_address` varchar(100) DEFAULT NULL COMMENT '转入地址',
  `contract_address` varchar(100) DEFAULT NULL COMMENT '合约地址',
  `amount` decimal(20,4) NOT NULL DEFAULT 0.0000 COMMENT '充值数量',
  `block_height` bigint NOT NULL DEFAULT 0 COMMENT '块高',
  `hash` varchar(100) NOT NULL COMMENT '交易哈希',
  `log_index` int(11) NOT NULL DEFAULT '0' COMMENT '日志序号',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_hash_log_index` (`hash`,`log_index`) USING BTREE,
  KEY `chain_recharge_uid_index` (`uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='链充值表';
