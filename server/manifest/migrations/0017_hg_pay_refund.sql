-- hg_pay_refund
CREATE TABLE IF NOT EXISTS `hg_pay_refund` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '会员ID',
  `app_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '应用ID',
  `order_sn` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '业务订单号',
  `refund_trade_no` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '退款交易号',
  `refund_money` decimal(10,2) DEFAULT NULL COMMENT '退款金额',
  `refund_way` tinyint(4) DEFAULT '1' COMMENT '退款方式',
  `ip` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '申请者IP',
  `reason` varchar(255) DEFAULT NULL COMMENT '申请退款原因',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '退款备注',
  `status` tinyint(4) DEFAULT '1' COMMENT '退款状态',
  `created_at` datetime DEFAULT NULL COMMENT '申请时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `order_sn` (`order_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付_退款记录';
