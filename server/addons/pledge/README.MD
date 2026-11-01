## 质押管理

### 简介

质押投资插件：管理投资挡位（V1~V5）、投资套餐（3/6/12 个月）、直推收益、月化收益、每日交易手数、每手金额和提前退本违约金，并提供前台只读的套餐查询与收益试算接口。

### 数据表

- `hg_addon_pledge_level`：投资挡位（挡位名称、投资金额、排序、状态）。
- `hg_addon_pledge_plan`：投资套餐（挡位、周期、月化收益展示区间、固定月化率、直推率、每日手数、每手金额、违约金比例、状态）。

安装插件时自动执行 `migrations/` 建表与 `seeders/` 初始化数据，重复安装不会覆盖后台已修改的配置。

### 后台功能

安装插件后，后台「质押管理」下包含：

- 挡位管理：`/admin/pledge/level/*`
- 套餐管理：`/admin/pledge/plan/*`
- 配置管理：`/admin/pledge/config/*`

配置项存放在 `hg_sys_addons_config`，当前包含：

- `depositAddress`：入金地址。
- `orderEnabled`：是否开始下单。

新增参数时，在 `seeders/hg_sys_addons_config.sql` 增加一行（`addon_name='pledge'`、`group='basic'`），如需在代码中读取，再同步补充 `model.BasicConfig` 字段即可，前端配置页会自动渲染。

### 前台接口

无需登录：

- `GET /api/pledge/plan/list`：全部启用套餐，仅返回 `levelName`、`amount`、`periodMonths`、`monthlyRewardInfo`、`monthlyRewardRate`、`totalReturn`。
- 首页挡位由套餐列表按 `levelName` 反推，不再单独提供 C 端挡位列表接口；后台仍使用 `/admin/pledge/level/option`。
- `GET /api/pledge/plan/quote?levelName=V3&periodMonths=3`：单个套餐收益试算，挡位使用公开名称而非数据库ID。

需要登录：

- `POST /api/pledge/order/buy`：提交报单，参数 `levelName`、`periodMonths`。
- `GET /api/pledge/order/list?page=&pageSize=`：当前用户正式报单列表。
- `POST /api/pledge/order/refund`：提前退本，参数 `orderNo`，退款入账 USDT 资产。
- `GET /api/pledge/order/status?orderNo=`：查询报单支付处理状态。

下单流程：

1. `order/buy` 写入 `hg_addon_pledge_add_order`，`status=0`，同一用户同时只允许一条待处理报单；提交新报单时旧待处理报单自动置为 `status=2`（已取消/被替代）。
2. 用户向链上平台地址支付套餐原价。
3. chain 扫块写入 `chain_recharge`，通过 Redis `users:wallet:uid` 解析付款地址对应的用户ID。
4. pledge 按 `user_id + amount + status=0` 匹配报单，用 `plan_id` 查套餐，写入 `hg_addon_pledge_order`，并把报单 `status` 改为 `1`。
5. 前端每 3 秒轮询 `order/status`，直到 `status=1`。

试算规则：单利、每月按 30 天、金额保留 2 位小数；提前退本到账 = 投资金额 - 违约金（当前无订单，累计已释放按 0）。

### 迁移或安装

1、安装 HotGo (2.19.0及以上)

项目介绍：https://github.com/bufanyun/hotgo

2、将当前插件项目拷贝进 HotGo 根目录的 server/addons 目录下

3、在 HotGo 根目录的 server/addons/modules 目录下创建go文件:pledge.go，内容如下：
```go
package modules

import _ "hotgo/addons/pledge"
```

4、HotGo 后台进入 开发工具->插件管理->找到 质押管理 (pledge) 进行安装

5、重启服务即可生效

### 常用命令行

```shell
# 接口维护-gen service
gf gen service -s=addons/pledge/logic -d=addons/pledge/service

# 生成插件表模型
gf gen dao -p addons/pledge -d model/dao -o model/do -e model/entity \
  -t hg_addon_pledge_level,hg_addon_pledge_plan -r hg_addon_ -c -j CamelLower

# 单元测试
go test ./addons/pledge/...
```
