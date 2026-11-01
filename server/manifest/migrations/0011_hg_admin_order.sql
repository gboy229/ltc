-- hg_admin_order
CREATE TABLE IF NOT EXISTS `hg_admin_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `member_id` bigint(20) DEFAULT '0' COMMENT '管理员id',
  `order_type` varchar(32) NOT NULL COMMENT '订单类型',
  `product_id` bigint(20) DEFAULT NULL COMMENT '产品id',
  `order_sn` varchar(64) DEFAULT '' COMMENT '关联订单号',
  `money` decimal(10,2) NOT NULL COMMENT '充值金额',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `refund_reason` varchar(255) DEFAULT NULL COMMENT '退款原因',
  `reject_refund_reason` varchar(255) DEFAULT NULL COMMENT '拒绝退款原因',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `order_sn` (`order_sn`),
  KEY `member_id` (`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_充值订单';
