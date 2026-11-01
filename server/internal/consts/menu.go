// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import "hotgo/internal/library/menu"

// init 注册后台基础菜单（目录 + 菜单，按钮由 Ops 自动生成）
// 排序按注册顺序、目录默认重定向到第一个子菜单、目录默认始终显示，均无需手动指定
func init() {
	menu.Register(
		// Dashboard
		&menu.Item{Name: "Dashboard", Title: "数据看板", Path: "/dashboard", Icon: "DashboardOutlined", Dir: true, Sort: 1,
			Children: []*menu.Item{
				{Name: "dashboard_console", Title: "主控台", Path: "console", Component: "/dashboard/console/console", Perm: "/console/stat"},
			},
		},

		// 组织权限
		&menu.Item{Name: "Org", Title: "组织权限", Path: "/org", Icon: "SafetyCertificateOutlined", Dir: true, Sort: 2,
			Children: []*menu.Item{
				{Name: "user", Title: "后台用户", Path: "user", Component: "/org/user/user", Api: "member", Perm: "/dept/list,/post/list,/role/list,/member/list,/dept/option", Ops: "edit,delete,status"},
				{Name: "org_dept", Title: "部门管理", Path: "dept", Component: "/org/dept/dept", Ops: "edit,delete,status"},
				//{Name: "org_post", Title: "岗位管理", Path: "post", Component: "/org/post/post", Ops: "edit,delete,status"},
				//{Name: "permission_menu", Title: "菜单权限", Path: "menu", Component: "/permission/menu/menu", Ops: "edit,delete,status"},
				{Name: "permission_role", Title: "角色权限", Path: "role", Component: "/permission/role/role", Perm: "/role/list,/role/dataScope/select,/role/getPermissions", Ops: "edit,delete,status"},
			},
		},

		/*
			// 运营中心
			&menu.Item{Name: "Operations", Title: "运营中心", Path: "/operations", Icon: "ClusterOutlined", Dir: true, Hidden: true,
				Children: []*menu.Item{
					{Name: "apply_notice", Title: "通知公告", Path: "notice", Component: "/apply/notice/index", Ops: "edit,delete,status,switch"},
					// 资金管理
					{Name: "asset", Title: "资金管理", Path: "asset", Dir: true,
						Children: []*menu.Item{
							{Name: "creditsLogIndex", Title: "资金变动", Path: "creditsLogIndex", Component: "/asset/creditsLog/index", Api: "creditsLog", Perm: "/creditsLog/list,/creditsLog/option", Ops: "export"},
							{Name: "asset_recharge", Title: "在线充值", Path: "recharge", Component: "/asset/recharge/recharge", Perm: "/order/create,/order/option",
								Buttons: []*menu.Button{
									{Name: "asset_recharge_log", Title: "充值记录", Permissions: "/order/list"},
									{Name: "asset_recharge_accept_refund", Title: "受理退款", Permissions: "/order/acceptRefund"},
									{Name: "asset_recharge_apply_refund", Title: "申请退款", Permissions: "/order/applyRefund"},
									{Name: "asset_recharge_order_delete", Title: "删除关闭订单", Permissions: "/order/delete"},
								}},
							{Name: "asset_cash", Title: "提现管理", Path: "cash", Component: "/asset/cash/index", Perm: "/cash/list,/cash/apply,/config/getCash",
								Buttons: []*menu.Button{
									{Name: "cash_apply", Title: "申请提现", Permissions: "/cash/apply"},
									{Name: "cash_payment", Title: "提现打款处理", Permissions: "/cash/payment,/cash/view"},
								}},
						},
					},
					{Name: "apply_provinces", Title: "地区编码", Path: "provinces", Component: "/apply/provinces/index"},
				},
			},

		*/
		/*
			// 系统管理
			&menu.Item{Name: "System", Title: "系统管理", Path: "/system", Icon: "SettingOutlined", Dir: true, Sort: 3,
				Children: []*menu.Item{
					{Name: "system_config", Title: "配置管理", Path: "config", Component: "/system/config/system", Perm: "/config/get",
						Buttons: []*menu.Button{
							{Name: "config_get", Title: "查看配置", Permissions: "/config/get"},
							{Name: "config_update", Title: "更新配置", Permissions: "/config/update"},
						}},
					{Name: "system_dict", Title: "字典管理", Path: "dict", Component: "/system/dict/index", Perm: "/dictType/tree,/dictData/list", Ops: "edit,delete"},
					{Name: "system_cron", Title: "定时任务", Path: "cron", Component: "/system/cron/index", Perm: "/cron/list,/cronGroup/list", Ops: "edit,delete,status"},
					{Name: "system_blacklist", Title: "黑名单", Path: "blacklist", Component: "/system/blacklist/index", Ops: "edit,delete,status"},
					// 日志管理
					{Name: "Logs", Title: "日志管理", Path: "log", Dir: true,
						Children: []*menu.Item{
							{Name: "log", Title: "全局日志", Path: "log", Component: "/log/log/index", Api: "log", Ops: "delete"},
							{Name: "login_log", Title: "登录日志", Path: "login-log", Component: "/log/login-log/index", Api: "loginLog", Ops: "delete,export"},
							{Name: "sms_log", Title: "短信记录", Path: "sms", Component: "/log/sms-log/index", Api: "smsLog", Ops: "delete"},
							{Name: "monitor_serve_log", Title: "服务日志", Path: "serve_log", Component: "/monitor/serve-log/index", Api: "serveLog", Ops: "delete,export"},
						},
					},
					{Name: "monitor_serve_monitor", Title: "服务监控", Path: "serve_monitor", Component: "/monitor/serve-monitor/index"},
					{Name: "monitor_online", Title: "在线用户", Path: "online", Component: "/monitor/online/index", Perm: "/monitor/userOnlineList",
						Buttons: []*menu.Button{
							{Name: "monitorUserOffline", Title: "强制退出用户", Permissions: "/monitor/userOffline"},
						}},
					{Name: "monitor_netconn", Title: "在线服务", Path: "netconn", Component: "/monitor/netconn/index"},
					{Name: "apply_attachment", Title: "附件管理", Path: "attachment", Component: "/apply/attachment/index"},
				},
			},
		*/
		// 开发工具
		&menu.Item{Name: "Develops", Title: "开发工具", Path: "/develop", Icon: "CodeOutlined", Dir: true, Sort: 10000,
			Children: []*menu.Item{
				//{Name: "develop_code", Title: "代码生成", Path: "code", Component: "/develop/code/index", Api: "genCodes"},
				{Name: "develop_code_deploy", Title: "生成配置", Path: "deploy/:id?", Component: "/develop/code/deploy", ActiveMenu: "develop_code", Hidden: true},
				{Name: "develop_addons", Title: "插件管理", Path: "addons", Component: "/develop/addons/index"},
			},
		},

		/*
			// 个人中心（隐藏，通过右上角头像进入，保留以生成路由）
			&menu.Item{Name: "home", Title: "个人中心", Path: "/home", Icon: "UserOutlined", Dir: true, Hidden: true,
				Children: []*menu.Item{
					{Name: "home_account", Title: "个人设置", Path: "account", Component: "/home/account/account",
						Buttons: []*menu.Button{
							{Name: "member_update_pwd", Title: "修改登录密码", Permissions: "/member/updatePwd"},
							{Name: "member_update_mobile", Title: "换绑手机号", Permissions: "/member/updateMobile"},
							{Name: "member_update_email", Title: "换绑邮箱", Permissions: "/member/updateEmail"},
							{Name: "member_update_profile", Title: "更新基本设置", Permissions: "/member/updateProfile"},
							{Name: "member_update_cash", Title: "更新提现设置", Permissions: "/member/updateCash"},
						}},
					{Name: "home_message", Title: "我的消息", Path: "message", Component: "/home/message/message", Perm: "/notice/messageList"},
				},
			},

		*/
	)
}
