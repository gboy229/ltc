-- 质押插件配置项初始化（幂等，INSERT IGNORE）
INSERT IGNORE INTO `hg_sys_addons_config`
(`addon_name`, `group`, `name`, `type`, `key`, `value`, `default_value`, `sort`, `tip`, `is_default`, `status`) VALUES
('pledge', 'basic', '是否开启', 'bool', 'orderEnabled', 'false', 'false', 20, '关闭后用户不能发起新的质押下单', 1, 1)
,('pledge', 'basic', '入金地址', 'string', 'depositAddress', '', '', 10, '用户前台入金时展示的钱包地址', 1, 1);
