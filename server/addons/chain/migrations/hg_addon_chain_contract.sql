-- 链合约表
CREATE TABLE IF NOT EXISTS `hg_addon_chain_contract` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain` varchar(30) DEFAULT NULL COMMENT '链名',
  `address` varchar(100) DEFAULT NULL COMMENT '合约地址',
  `symbol` varchar(30) DEFAULT NULL COMMENT '币种',
  `decimal` tinyint NOT NULL DEFAULT 18 COMMENT '精度',
  `min_amount` decimal(20,8) NOT NULL DEFAULT 0.00000000,
  `max_amount` decimal(20,8) NOT NULL DEFAULT 0.00000000,
  `multiple_amount` int NOT NULL DEFAULT 0 COMMENT '倍数',
  `verify_amount` decimal(20,8) NOT NULL DEFAULT 0.00000000,
  `scale` tinyint NOT NULL DEFAULT 0 COMMENT '保留几位小数',
  `day_limit` int NOT NULL DEFAULT 0 COMMENT '日上限笔数',
  `fee_type` tinyint NOT NULL DEFAULT 0 COMMENT '手续费类型，0=百分比，1=个数',
  `rate` decimal(20,8) NOT NULL DEFAULT 0.00000000,
  `recharge_switch` tinyint NOT NULL DEFAULT 0 COMMENT '充币开关',
  `withdraw_switch` tinyint NOT NULL DEFAULT 0 COMMENT '提现开关',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `scan_block` bigint unsigned NOT NULL DEFAULT 0 COMMENT '已扫块高',
  `scan_block_hash` varchar(66) NOT NULL DEFAULT '' COMMENT '已扫块哈希',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `chain_contract_chain_address_unique` (`chain`,`address`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='链合约表';
