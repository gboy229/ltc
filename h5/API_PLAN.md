# H5 C 端 API 对接与后台完善规划

> 更新日期：2026-10-28
> 范围：`h5/`（C 端 H5） ↔ `server/`（HotGo API）
> 当前进度：首页 + 钱包登录已对接，公告公开接口已补齐。

---

## 一、已完成的首页链路

### 1. 钱包自动登录

```
进入 H5
  └─ App.vue onMounted → userStore.bootstrap()
       ├─ EIP-6963 + 注入式全局变量发现钱包
       │    MetaMask / TokenPocket / OKX / Bitget / Trust / imToken / 其它 EVM 钱包
       ├─ 已有 token + 地址：eth_accounts 静默校验 → 恢复会话（不弹签名）
       ├─ 无会话：eth_requestAccounts → GET nonce → personal_sign → POST login
       └─ 无钱包环境：弹出钱包选择/下载引导
```

前端文件：

| 文件 | 作用 |
|------|------|
| `h5/src/web3/wallet.ts` | 多钱包发现、连接、签名、账户/链变化监听、deep link |
| `h5/src/stores/user.ts` | 自动登录、会话恢复、退出、登录态 |
| `h5/src/api/users/index.ts` | nonce / login / me / logout |
| `h5/src/components/WalletSheet.vue` | 未检测到钱包时的选择/下载面板 |
| `h5/src/components/AppHeader.vue` | 显示真实钱包地址 / 连接入口 |

已对接接口：

| 接口 | 说明 |
|------|------|
| `GET /api/users/auth/nonce?address=` | 获取签名挑战值 |
| `POST /api/users/auth/login` | `{type:"wallet", account, signature, nonce, uid?}` |
| `GET /api/users/center/info` | 当前用户信息（含钱包地址、等级） |
| `POST /api/users/auth/logout` | 退出登录 |

### 2. 首页数据

| 页面区域 | 接口 | 状态 |
|----------|------|------|
| 公告上下滚动 | `GET /api/cms/article/recommend` | 已对接（全部推荐公告，不分页） |
| 入金挡位 | 由 `GET /api/pledge/plan/list` 的 `levelName/amount` 反推 | 已对接 |
| 投资周期 / 套餐 | `GET /api/pledge/plan/list` | 已对接 |
| 收益试算 | `GET /api/pledge/plan/quote` | 后端已具备，首页直接用 list 返回的试算字段 |
| 报单明细 | `GET /api/pledge/order/status`、`/pending` | 支付状态已对接，列表页待接 |
| 立即报单 | `POST /api/pledge/order/buy` + 钱包支付 | 已对接 |

> 公告前台接口不暴露数据库主键：新增 `article_no` 公开 ID（自动生成），前台统一用该 ID；
> `list`/`recommend` 只返回 `id/title/createdAt`，`view` 只返回 `id/title/content/createdAt`。
> 首页推荐由后台 `is_recommend` 控制，默认推荐。

### 3. 报单支付

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/chain/config` | 后台配置的支付链、RPC、合约、平台收款地址 |
| POST | `/api/pledge/order/buy` | 提交报单，参数 `levelName`、`periodMonths` |
| GET | `/api/pledge/order/list` | 当前用户正式报单列表，支持 `page/pageSize` |
| GET | `/api/pledge/order/status?orderNo=` | 轮询支付处理状态 |
| GET | `/api/chain/withdrawal/options` | 可提币网络/手续费 |
| POST | `/api/chain/withdrawal/apply` | 提交提币申请并冻结资产 |
| GET | `/api/chain/withdrawal/list` | 当前用户提币记录 |

流程：

1. 注册/登录时把钱包地址与 user_id 写入 Redis Hash `users:wallet:uid`，写入失败注册/登录失败。
2. `buy` 写入 `hg_addon_pledge_add_order`，同一用户只允许一条 `status=0` 报单；新报单会把旧待处理报单置为 `status=2`。
3. 前端按 `/api/chain/config` 切换钱包网络，向平台地址支付套餐原价。
4. chain 扫块校验合约和平台地址，按 `(hash, log_index)` 写入 `chain_recharge`，通过 Redis 解析付款地址对应的 user_id。
5. chain 推送 `chain_payment` 队列，pledge 按 `user_id + amount + status=0` 匹配报单，写入 `hg_addon_pledge_order` 并把报单置为已处理。
6. 前端每 3 秒轮询 `/order/status`，成功后显示订单号。

---

## 二、其它页面 API 规划

### 1. 平台公告（NoticeView / NoticeDetailView / 首页滚动）

**后端已完成：**

| 方法 | 路径 | 入参 | 出参 |
|------|------|------|------|
| GET | `/api/cms/article/recommend` | 无 | `{list:[{id,title,createdAt}]}`，全部推荐公告，不分页 |
| GET | `/api/cms/article/list` | `page`、`pageSize`、`title?` | `{list:[{id,title,createdAt}], totalCount}` |
| GET | `/api/cms/article/view/{article_no}` | `article_no`（路径参数） | `{id,title,content,createdAt}` |

说明：

- `id` 为自动生成的 `article_no`，不是数据库主键，避免暴露自增 ID。
- 列表不返回 `cover/sort/status/type/content` 等前台用不到的字段。
- `hg_addon_cms_article` 新增 `is_recommend`（默认 1=推荐）；后台公告编辑页增加“推荐首页”开关。
- 前端公告列表、公告详情已改为真实接口；首页使用上下滚动组件展示全部推荐公告。
- 正文 HTML 在详情页做了 `DOMParser` 清洗（移除 script/style/iframe/事件属性等）后再渲染。

> 已有数据库需要执行一次升级 SQL：
>
> ```sql
> ALTER TABLE hg_addon_cms_article
>   ADD COLUMN article_no varchar(32) NOT NULL DEFAULT '' COMMENT '前台公开ID' AFTER id,
>   ADD COLUMN is_recommend tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否推荐：1=推荐，0=不推荐' AFTER status;
> UPDATE hg_addon_cms_article SET article_no = LOWER(SUBSTRING(MD5(CONCAT(id, RAND(), NOW())), 1, 20)) WHERE article_no = '';
> ALTER TABLE hg_addon_cms_article
>   ADD UNIQUE KEY uk_article_no (article_no),
>   ADD KEY idx_recommend (is_recommend, status);
> ```

---

### 2. 我的资产（AssetsView）

现有通用资产账户/流水表已经可用：

- `hg_sys_asset_account`：`user_id`、`available`、`frozen`
- `hg_sys_asset_log`：`change_type`、`biz_type`、`biz_id`、`biz_sn`、变动前后余额
- 主模块已有 `service.Asset()`：`Account`、`AccountList`、`Logs`、`Stat`、`Change`

**建议在 users 插件新增前台接口：**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/users/asset/account` | 当前用户资产账户（available/frozen） |
| GET | `/api/users/asset/logs` | `page/pageSize/changeType/bizType/createdAt[]` |
| GET | `/api/users/asset/stat` | 总资产、今日入账/出账、待释放 |

实现要点：用户 ID 必须从登录上下文取，**禁止**前端传 `userId`；`users` 插件直接调用主模块 `service.Asset()`，不要重复写账务逻辑。

---

### 3. 报单明细 + 立即报单（OrderDetailView / HomeView / BuyModal）

这是目前缺口最大的模块，建议在 `pledge` 插件内新增订单域。

#### 表设计（新增）

`hg_addon_pledge_order`：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| order_no | varchar(32) | 业务单号，唯一 |
| user_id | bigint | 用户 ID |
| plan_id / level_id | bigint | 套餐 / 挡位 |
| amount | decimal(20,8) | 入金金额 |
| period_months | int | 周期 |
| monthly_reward_rate | decimal(10,4) | 下单时锁定月化率 |
| status | tinyint | 0=待支付 1=进行中 2=已完成 3=已取消 4=提前退本 |
| pay_chain | varchar(30) | 支付链 |
| pay_symbol | varchar(20) | 币种 |
| pay_tx_hash | varchar(100) | 支付交易哈希 |
| pay_from | varchar(100) | 付款地址 |
| paid_at | datetime | 支付确认时间 |
| start_at / end_at | datetime | 收益起止 |
| released_amount | decimal(20,8) | 已释放收益 |
| created_at / updated_at | datetime | 标准字段 |

收益释放记录建议**复用 `hg_sys_asset_log`**（`biz_type='pledge_reward'`、`biz_id=order_id`），不要重复造流水表；订单表只存汇总值。

#### 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/pledge/order/create` | `{planId, payChain, paySymbol}` → 返回订单 + 收款地址/金额 |
| GET | `/api/pledge/order/list` | `page/pageSize/status`，仅当前用户 |
| GET | `/api/pledge/order/view` | `id`，校验归属 |
| POST | `/api/pledge/order/cancel` | 仅待支付可取消 |
| POST | `/api/pledge/order/earlyRefund` | 提前退本，按 `early_exit_penalty_rate` 试算并落单 |
| GET | `/api/pledge/order/summary` | 入金总额、今日释放、累计释放 |
| GET | `/api/pledge/order/rewardList` | 收益释放明细（或走 asset/logs） |

#### 支付流程（建议先做方案 A）

**方案 A（推荐，复用链插件）：**

1. 用户点击报单 → `order/create` 创建 `待支付` 订单，返回平台收款地址、金额、链。
2. 前端拉起钱包 `eth_sendTransaction` 转账 USDT。
3. 链插件扫块匹配 `from/to/amount/txHash` → 回调 pledge 订单模块确认支付。
4. 订单变为 `进行中`，同时写 `hg_sys_asset_log`（入金）。
5. 定时任务按日释放收益 → 写资产流水，更新 `released_amount`。

**方案 B（后续）：** 合约下单，后端校验交易回执与事件日志。需要先确定合约地址、链、ABI。

> 在订单模块完成前，首页“报单明细”保持空态；“立即报单”只做登录校验并提示“接入中”。

---

### 4. 收益明细（EarningsView）

数据来源有两类：

| 类型 | 来源 |
|------|------|
| 静态收益 / 释放收益 | `hg_sys_asset_log`（biz_type=pledge_reward） |
| 直推/团队收益 | 新增佣金表或 `hg_sys_asset_log`（biz_type=pledge_commission） |

**接口建议：**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/users/asset/logs` | 通用收支流水，前端按 `bizType` 过滤 |
| GET | `/api/pledge/order/summary` | 报单金额、待释放 |
| GET | `/api/users/community/summary` | 直推收益、团队收益汇总 |

若需要“收益总额/今日释放/累计释放”聚合，建议单独提供 `/api/pledge/reward/summary`，避免前端拉全量流水计算。

---

### 5. 我的社区（CommunityView）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/users/community/summary` | UID、等级、直推人数、团队人数、直推收益、团队收益 |
| GET | `/api/users/community/members` | 直推成员列表（uid、注册时间、等级、有效报单） |
| GET | `/api/users/community/invite` | 邀请地址/邀请码/邀请链接 |

关系数据已具备：`hg_addon_users_tree`、`hg_addon_users_user.pid`；佣金需要新增结算逻辑（见下）。

---

### 6. 在线提币（WithdrawView）

链插件已有 `hg_addon_chain_withdraw`（`uid/chain/symbol/from/to/pay_amount/receive_amount/fee/status/hash`）和后台审核接口，但**缺少前台接口**。

**建议在 chain 插件新增：**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/chain/config` | 可用链、币种、合约地址、最小/最大提币、手续费、开关 |
| GET | `/api/chain/deposit/config` | 入金收款地址、链、币种（报单/充值用） |
| POST | `/api/chain/withdraw/create` | `{chain, symbol, amount, toAddress}`，扣减可用余额并冻结 |
| GET | `/api/chain/withdraw/list` | 当前用户提币记录 |
| GET | `/api/chain/withdraw/view` | 详情（校验归属） |
| GET | `/api/chain/recharge/list` | 当前用户充值记录 |

实现要点：

- 费率、最小额、开关全部从 `hg_addon_chain_contract` 读取，前端不硬编码。
- 创建提币：**先校验资产可用余额 → 资产冻结 → 落提现单**，同一事务。
- 审核通过/打币失败时做解冻或扣减，保证账务一致。
- 提币地址必须做 EVM 地址校验；不同链的地址格式后续扩展。

---

### 7. 用户中心 / 全局

**已扩展 `GET /api/users/center/info`：**

```json
{
  "user": {
    "uid": "85285352",
    "nickname": "0x7a3f...9c4e",
    "avatar": "",
    "address": "0x7a3f...",
    "levelName": "普通会员"
  }
}
```

前端 `SideNav` 的 UID / 等级、头部钱包地址已改为真实数据，不再依赖 demo 常量。

**登录安全增强（后续）：**

- 当前 nonce 为随机串，已做 5 分钟缓存 + 一次性消费 + 防重放；
- 建议升级为 SIWE 风格消息：`domain + address + nonce + issuedAt + expirationTime + statement`，前端按同格式签名；
- 增加 `chainId` 上报，便于后续按链区分账户；
- 多钱包绑定：同一用户可绑定多个钱包地址（`hg_addon_users_user_account` 已支持多凭证）。

---

## 三、后台调整清单（按优先级）

### P0（首页闭环，当前）

- [x] CMS 公告公开列表/详情接口
- [x] 钱包 nonce/login/me/logout 前端对接
- [x] 首页挡位/套餐/试算数据对接
- [ ] 首页报单明细空态 + 登录校验（已完成空态）
- [x] CMS 公告公开接口 + 首页推荐接口 + 公开 ID

### P1（报单 + 资产）

- [ ] `hg_addon_pledge_order` 建表 + entity/dao/do 生成
- [ ] pledge 订单 create/list/view/cancel/earlyRefund/summary
- [ ] 链上支付确认（复用 chain 扫块/充值记录）
- [ ] users 前台资产接口（account/logs/stat）
- [ ] 收益释放定时任务（写资产流水）
- [ ] 首页报单明细接真实订单

### P2（社区 + 提币 + 公告页）

- [ ] users 社区 summary/members/invite
- [ ] 直推/团队佣金结算
- [ ] chain 前台 config/deposit/withdraw/recharge 接口
- [ ] 提币冻结/审核/解冻账务闭环
- [x] 公告列表/详情页替换假数据 + 正文 XSS 清洗

### P3（体验与安全）

- [x] `center/info` 返回钱包地址/等级名称
- [ ] SIWE 消息格式 + chainId
- [ ] 多钱包绑定/切换
- [ ] 接口聚合（首页 summary）减少请求数
- [ ] 错误码/文案统一

---

## 四、接口约定

- 所有前台接口统一走 `/api/<addon>/...`（插件前缀由 `addons.RouterPrefix` 生成）。
- 插件内 C 端逻辑放 `logic/<domain>/`，后台逻辑放 `logic/sys/`；C 端 service 单独放 `service/<domain>.go`，禁止把 C 端方法塞进 `ISysXxx`。
- 需要登录的接口放在插件 `ApiAuth` 中间件之后；用户身份从 `contexts.GetUsersUserId(ctx)` 获取。
- 统一响应：`{code:0,message:"操作成功",data:{...}}`；失败 `code:-1`，登录失效 `code:61`。
- 前端请求统一走 `h5/src/utils/http.ts` 的 `http.*`，禁止页面直接使用 axios。
- 金额字段统一 `decimal(20,8)`，前端展示保留 2 位；计算一律后端完成，前端只做展示/试算校验。
- 字典/枚举由后端提供，前端不硬编码状态文案。

---

## 五、当前本地验证

```bash
# 后端（8000 被占用时可换端口）
cd server && go run main.go

# 前端
cd h5 && pnpm dev
```

已验证：

- `GET /api/pledge/plan/list` 返回 15 个启用套餐（仅前端使用字段，无数据库主键）
- `GET /api/pledge/plan/quote?levelName=V3&periodMonths=3` 返回 5000/3 个月/总回报 6350
- `GET /api/cms/article/list`、`/api/cms/article/recommend`、`/api/cms/article/view/{article_no}` 返回已发布/推荐公告（公开 ID，无数据库主键）
- `POST /api/users/auth/login`（`users.wallet.skipVerify=true`）+ `/api/users/center/info` 正常
- `pnpm run typecheck`、`pnpm run build` 通过
