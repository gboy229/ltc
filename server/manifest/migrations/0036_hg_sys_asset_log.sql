-- 用户资产流水表：只追加，不修改；一笔业务的一种账务动作只允许一条
CREATE TABLE IF NOT EXISTS `hg_sys_asset_log` (
  `id`               bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID',
  `user_id`          bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `account`          varchar(30)  NOT NULL DEFAULT 'USDT' COMMENT '资产账户标识，如 USDT',
  `addon_name`       varchar(30)  NOT NULL DEFAULT 'system' COMMENT '来源插件',
  `change_type`      varchar(20)  NOT NULL DEFAULT 'incr' COMMENT '账务动作：incr=增加，decr=扣减，freeze=冻结，unfreeze=解冻，settle=扣减冻结',
  `biz_type`         varchar(50)  NOT NULL DEFAULT '' COMMENT '业务类型',
  `biz_id`           bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '业务单号',
  `biz_sn`           varchar(128) NOT NULL DEFAULT '' COMMENT '业务外部单号',
  `before_available` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '变动前可用余额',
  `change_available` decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '可用余额变动值（正负）',
  `after_available`  decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '变动后可用余额',
  `before_frozen`    decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '变动前冻结金额',
  `change_frozen`    decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '冻结金额变动值（正负）',
  `after_frozen`     decimal(20,8) NOT NULL DEFAULT '0.00000000' COMMENT '变动后冻结金额',
  `remark`           varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `operator_id`      bigint(20)   NOT NULL DEFAULT '0' COMMENT '操作人ID（后台管理员，0=系统）',
  `created_at`       datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_biz` (`account`,`addon_name`,`biz_type`,`biz_id`,`change_type`),
  KEY `idx_user` (`user_id`,`id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户资产流水表';
