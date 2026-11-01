# 区块管理插件（chain）待办与修复记录

> 更新于：2026-10-27
> 约定：`[x]` 已完成，`[ ]` 待处理，`[-]` 本轮进行中

## 已完成

- [x] P0-1 私钥不再出现在 API / 页面
  - 链配置列表、详情改为显式字段模型，SQL 不查询 `transfer_private`
  - 生成地址接口只返回 address
- [x] P0-2 私钥加密落库 + 生成地址改用 adapters + 地址历史
  - AES-256-GCM（`enc:v1:` 前缀），密钥来自 `chain.privateKeySecret`
  - `adapters.NewChainAdapter(driver).GenerateAddress()`
  - 生成结果写入 `hg_addon_chain_address`，同步更新 `chain_config.transfer_address`
- [x] P0-3 充值/提现不可修改、不可删除，提现仅审核
  - 移除 recharge/withdraw 的 delete 链路与菜单权限
  - 审核仅允许「审核中 → 审核通过 / 审核拒绝」
- [x] 菜单同步修复
  - 移除确定性 ID / `MenuIdBase`，按 `name` 判断新增或更新，ID 由数据库自增

## 已完成（P1 1-3，2026-10-27）

- [x] P1-1 `chain_config` 的 `driver` 唯一索引改为 `(driver, chain)` 唯一
  - 迁移逻辑：`addons/chain/migrations/hg_addon_chain_config.sql`
  - 当前库已应用：删除 `chain_config_driver_unique`，新增 `chain_config_driver_chain_unique`
  - 实测：同一 `evm` 驱动可新增第二条链配置
- [x] P1-2 `chainconfig` / `contract` 批量删除统一为 `ids []int64`
  - 后端改为 `WhereIn` 批量删除，前端单删/批删统一传 `ids`
  - 实测：批量删除后记录数为 0
- [x] P1-3 提现审核增加 `version` 乐观锁 + 审核人/审核时间/备注
  - 新增字段：`audit_by`、`audit_at`、`remark`；迁移逻辑：`addons/chain/migrations/hg_addon_chain_withdraw.sql`
  - 审核时 `version` 不匹配返回「已被其他操作修改」；重复审核被拒绝
  - 前端审核弹窗支持填写备注，列表展示审核人/审核时间/备注
  - 实测：审核后 status=1、version=1、remark 正确，重复审核报错

## 已完成（合约模块可运营优化，2026-10-27）

- [x] 列表只展示运营关注字段
  - 链、币种、合约地址、精度、提币限额（最小~最大）、手续费（百分比/固定）、充币/提币开关、创建时间
  - 扫块游标类系统字段（`browser_page`、`browser_limit`、`start_block`）不再返回和展示
- [x] 编辑页按 Tab 分组
  - 基础信息 / 提币规则 / 手续费 / 功能开关
  - 每个参数都补充了作用说明和填写示例
- [x] 链名改为下拉选择
  - 新增 `/chain/chainconfig/options` 接口，返回已配置的链
  - 避免手输 `BEP20` 而实际配置为 `bsc` 时添加失败
- [x] 合约地址校验改为链驱动
  - 前端只做必填/长度校验，不再写死 EVM 的 `0x` 正则
  - 后端通过 `ChainAdapter.IsValidAddress` 校验，按链配置的 `driver` 自动选择
  - 新增 Tron 等链时只需实现 adapter + factory + 驱动选项，合约页面无需改动
- [x] 手续费口径明确
  - 百分比：`rate=1` 表示 1%，手续费 = 提币数量 × 费率 ÷ 100
  - 固定个数：`rate=1` 表示每笔收取 1 个代币
- [x] 系统字段从编辑请求中移除
  - `ContractUpdateFields` / `ContractInsertFields` 不再包含扫块游标
- [x] 后端校验补齐
  - 合约地址格式、精度/小数位范围、限额关系、手续费范围、开关枚举
  - 链配置必须存在；同链合约地址唯一（大小写不敏感）
  - 链名、合约地址、精度创建后不可修改（前后端双重校验），避免扫块游标和历史流水错乱
- [x] 数据库
  - 新增唯一索引 `chain_contract_chain_address_unique (chain, address)`
  - `min_amount` / `max_amount` / `verify_amount` / `rate` 调整为 `decimal(20,8)`
- [x] 实测
  - 新增合约成功；重复地址、未知链、非法地址均被拒绝
  - 列表模型不包含扫块系统字段

## 数据库变更（2026-10-27）

- `hg_addon_chain_config`
  - 删除索引 `chain_config_driver_unique`
  - 新增唯一索引 `chain_config_driver_chain_unique (driver, chain)`
  - 删除 `browser_api_url`
  - 新增 `confirmations int NOT NULL DEFAULT 15`
- `hg_addon_chain_contract`
  - 新增唯一索引 `chain_contract_chain_address_unique (chain, address)`
  - `min_amount` / `max_amount` / `verify_amount` / `rate` 改为 `decimal(20,8)`
  - 删除 `browser_page`、`browser_limit`、`start_block`
  - 新增 `scan_block bigint unsigned`、`scan_block_hash varchar(66)`
- `hg_addon_chain_withdraw`
  - 新增 `audit_by bigint(20) NOT NULL DEFAULT 0`（审核人ID）
  - 新增 `audit_at datetime NULL`（审核时间）
  - 新增 `remark varchar(255) NOT NULL DEFAULT ''`（审核备注）
- 上述变更由 `addons/chain/migrations/*.sql` 建表语句覆盖（当前库已是最新结构），插件版本 `v1.0.4`

## 已完成（充值记录列表优化，2026-10-27）

- [x] 列表按运营视角重排
  - ID（方便直接定位记录）、用户ID（`#uid`）、用户钱包（转出地址）、充值地址（转入）、币种、充值数量（绿色高亮）、类型、交易哈希、充值时间
  - 合约地址、块高等次要字段移入「详情」弹窗
- [x] 详情弹窗
  - 展示记录ID、用户ID、链/币种、类型、数量、用户钱包、充值地址、合约地址、块高、哈希、时间
  - 地址/哈希支持一键复制
- [x] 搜索全部展开，并补充参数
  - 用户关键词（用户ID/钱包地址）、链、币种、类型、用户钱包、充值地址、合约地址、交易哈希、金额区间、块高区间、充值时间范围
  - 关键词：纯数字按充值 `uid` 精确匹配，其余按 `from_address` / `to_address` 模糊匹配
- [x] 后端
  - 新增充值类型字典 `ChainRechargeTypeOptions`
  - 列表增加金额/块高/时间范围、合约地址筛选
  - 不依赖主模块用户表，充值记录本身没有用户名/姓名/手机号字段
- [x] 实测
  - 505 条记录分页正常；金额/块高/合约/地址/时间筛选均正常
  - 钱包地址关键词：完整/部分转出地址 140 条，充值地址 504 条，无匹配地址 0 条
  - 纯数字关键词按 `uid` 精确匹配（当前 `uid=17` 共 139 条）
- [ ] 后续如需展示用户名：先确认业务用户来源表（`hg_admin_member` 是后台用户；users 插件表字段是 `nickname`），再接对应服务

## 已完成（提现记录列表优化，2026-10-27）

- [x] 列表按运营视角重排
  - ID、用户ID（`#uid`）、用户钱包（转入地址）、出币地址（转出地址）、币种、申请金额、到账金额（绿色高亮）、手续费、状态、申请时间
  - 金额列不重复带 symbol；列表不展示交易哈希
  - 审核时间、合约地址、交易哈希等移入「详情」弹窗
- [x] 详情弹窗
  - 展示金额、地址、哈希、状态、审核时间、申请时间
  - 不展示审核人、审核备注、版本号、更新时间
  - 地址/哈希支持一键复制
- [x] 搜索全部展开，并补充参数
  - 用户关键词（用户ID/钱包地址）、链、币种、状态、用户钱包、出币地址、合约地址、交易哈希、金额区间、申请时间范围、审核时间范围
- [x] 后端
  - 列表增加出币地址、合约地址、金额、申请/审核时间筛选
  - 关键词：纯数字按 `uid` 精确匹配，其余按 `from_address` / `to_address` 模糊匹配
  - 详情接口增加记录不存在提示
- [x] 测试数据
  - 已向 `hg_addon_chain_withdraw` 插入 14 条演示数据（ID 2~15），覆盖审核中/已审核/已打币/失败/完成/取消
  - 覆盖 usdt/bnb、不同金额与时间；需要清理时可执行 `DELETE FROM hg_addon_chain_withdraw WHERE id BETWEEN 2 AND 15`
  - `type` 列已由业务删除，entity/do/dao 已重新生成，前后端 type 相关代码已清理
- [x] 实测
  - 总记录 14 条；status=0 3 条；keyword=17 3 条；钱包地址 3 条；金额/时间筛选正常

## 进行中（扫块 Phase 1，2026-10-27）

- [x] 数据结构
  - `hg_addon_chain_config`：删除 `browser_api_url`，新增 `confirmations`（默认 15）
  - `hg_addon_chain_contract`：删除 `browser_page`、`browser_limit`、`start_block`，新增 `scan_block`、`scan_block_hash`
  - 不新增 cursor 表；游标直接放合约表
- [x] Adapter
  - `GetBlockNumber`、`GetBlockHash`、`ScanTransferLogs`
  - EVM 用 `eth_getLogs`：块高区间 + 合约地址 + `Transfer` 事件，不按 `to` 过滤
  - 扫描端只负责解析和推送该合约的全部 Transfer；是否属于平台/用户由消费端判断
  - 只处理标准 ERC20 `Transfer`，不处理内部转账/原生币
  - 严格校验 3 个 topic + 32 字节 data；ERC721/非标准 Transfer 跳过
- [x] Scanner
  - 按链加载 `recharge_switch=1` 的合约
  - `confirmations` 控制只扫 `latest - confirmations`
  - 每批 1000 块，单次最多 5 批
  - `eth_getLogs` 遇到 `limit exceeded` 等范围/结果限制时自动二分区间重试
  - `scan_block=0` 时首次从当前确认块开始；需要历史时手工设置 `scan_block` 并清空 `scan_block_hash`
  - 用 `scan_block_hash` 做 reorg 检测，异常回退 100 块
  - 推送成功后更新游标，避免丢事件
- [x] 队列
  - topic：`chain:scan:transfer`
  - 队列驱动切换为 `redis`
  - 消费端 `queues/scan.go` 目前只记录日志，业务归属/落库/入账留待后续
- [x] 定时任务
  - `chain_scan`，每 5 秒，单例策略
  - `Install/Upgrade/UnInstall` 同步/删除插件定时任务
- [x] 前端
  - 链配置编辑页删除浏览器API地址，新增确认数（BSC 15 / ETH 12 提示）
  - 链配置列表增加确认数列
- [x] 诊断命令
  - `server/cmd/chainscan`：读取链/合约配置，直接调用 `ScanTransferLogs` 打印块高、from/to、raw 金额、换算金额、txHash、logIndex
  - 支持 `-contract` 临时指定未入库合约、`-rpc` 临时覆盖 RPC、`-from/-to` 指定区间、`-json` 输出
  - 固定把 `eth_getLogs` 原始 JSON 响应追加写入 `server/logs/chainscan/transfers.jsonl`，不需要 `-log` 参数
  - 日志包含 meta/scan/raw_request 行 + 原始 JSON-RPC 响应（address/topics/data/blockNumber 等原样保留）
  - 不再写入解析后的 Transfer 结构；`-max` 只限制控制台
- [x] 验证
  - `go build ./...`、`go vet ./addons/chain/...`、`go test ./addons/chain/...`、前端 ESLint 通过
  - 用本地假 JSON-RPC 验证了 `eth_getLogs` 解析（合约、from/to、amount、txHash、logIndex、blockHeight）
  - 用假 RPC + 真实 DB/Redis 跑通了「扫块 → 队列 → 更新 scan_block/scan_block_hash」，事务回滚
  - 真实链上：公共 RPC 已通，`scan_block` 推进到 `130012399`；随后因 USDT 全量日志触发 `limit exceeded`
  - 已加二分区间重试，并用假 RPC 验证（大区间报 limit 时自动拆小后返回数据）
  - id=1 合约地址原来是 43 位无效地址，已修正为真实 BSC USDT
  - 当前运行中的后端进程需要重启才能加载二分重试代码
- [ ] 待确认/后续
  - 合约地址：id=1 的 USDT 地址与真实 BSC USDT 不一致，id=5/7 是测试地址；扫描前需确认
  - 平台收款地址不再参与扫描过滤；消费端需要按 `to` 判断是否属于平台地址
  - USDT 等高频合约会推送全量 Transfer，消费端需关注队列量和过滤性能
  - 队列消费端业务逻辑（from -> uid、幂等、落库、入账）
  - `chain_recharge` 暂未加 `log_index`/唯一索引（按你的要求先不动）

## 后续

- [x] P1-4 补齐 chain 的 `migrations/` SQL 建表文件（与 users 插件对齐）
- [ ] P1-4 如需默认数据再补 `seeders/`（当前 chain 无默认数据）
  - 建表 SQL 不要创建 `driver` 唯一索引，并包含 `audit_by/audit_at/remark`
  - 合约表使用 `(chain, address)` 唯一索引，金额字段 `decimal(20,8)`
  - 充值/提现流水的金额精度目前是 `decimal(20,4)`，P3 提现流程实现时建议统一到 8 位
- [ ] P2-1 用户充值地址分配（`chain_address` 按 uid + chain 生成/复用）
- [ ] P2-2 扫块任务 + 队列，充值落 `chain_recharge`（hash 幂等）
- [ ] P2-3 充值确认后入账（需先定确认数、归集策略）
- [ ] P3-1 提现申请 API（开关/限额/手续费/精度校验）
  - 手续费口径：百分比 `rate=1` 表示 1%，`fee = 提币数量 × rate ÷ 100`；固定个数时 `fee = rate`
- [ ] P3-2 审核通过后签名打币（解密私钥 → SendTransaction → 回写 hash）
- [ ] P3-3 回执轮询，更新已打币/完成/失败
- [ ] 安全加固：密钥环境变量/KMS、AAD 绑定 chain+address、密钥轮换、历史数据迁移
- [ ] 清理：模板 Test 接口、补 chain 配置菜单、审计日志

## 待确认的业务规则

1. 充值归集策略：用户地址直接入账，还是归集到 `transfer_address`？
2. 到账确认数：EVM 链几个区块确认后入账？
3. 提现打币方式：审核通过后自动打币，还是再点一次「打币」人工触发？
