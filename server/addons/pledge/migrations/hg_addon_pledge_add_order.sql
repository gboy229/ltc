-- 用户提交报单（支付前的待匹配记录）
CREATE TABLE IF NOT EXISTS `hg_addon_pledge_add_order` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_no` varchar(32) NOT NULL DEFAULT '' COMMENT '公开订单号',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `plan_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '套餐ID',
  `amount` decimal(20,2) NOT NULL DEFAULT '0.00' COMMENT '报单金额',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0=待处理，1=已处理，2=已取消/被新报单替代',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_status` (`user_id`,`status`),
  KEY `idx_amount_status` (`amount`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户提交报单表';
