-- 插件配置项初始化（幂等，INSERT IGNORE）
INSERT IGNORE INTO `hg_sys_addons_config` (`addon_name`, `group`, `name`, `type`, `key`, `value`, `default_value`, `sort`, `tip`, `is_default`, `status`) VALUES
('users', 'login', '注册开关', 'bool', 'registerSwitch', 'true', 'true', 10, '是否开启注册', 1, 1),
('users', 'login', '强制邀请', 'bool', 'forceInvite', 'false', 'false', 20, '注册是否强制填写邀请码', 1, 1);
