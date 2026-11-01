-- 用户资产账户表：一个用户一条记录，当前仅余额
CREATE TABLE IF NOT EXISTS `hg_sys_asset_account` (
  `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '账户ID',
  `user_id`    bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID（hg_addon_users_user.id）',
  `account`    varchar(30)  NOT NULL DEFAULT 'USDT' COMMENT '资产账户标识，如 USDT',
  `available`  decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '可用余额',
  `frozen`     decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '冻结金额',
  `version`    int(11)      NOT NULL DEFAULT '0' COMMENT '乐观锁版本号',
  `status`     tinyint(1)   NOT NULL DEFAULT '1' COMMENT '状态：1=正常，2=禁用',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_account` (`user_id`,`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户资产账户表';
