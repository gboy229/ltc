## 区块管理

### 简介




### 使用说明




### 迁移或安装

1、安装 HotGo (2.19.0及以上)

项目介绍：https://github.com/bufanyun/hotgo

2、将当前插件项目拷贝进 HotGo 根目录的 server/addons 目录下

3、在 HotGo 根目录的 server/addons/modules 目录下创建go文件:chain.go，内容如下：
```go
package modules

import _ "hotgo/addons/chain"
```

4、HotGo 后台进入 开发工具->插件管理->找到 区块管理 (chain) 进行安装

5、重启服务即可生效


### 前台接口

无需登录：

- `GET /api/chain/config`：返回钱包支付所需的链ID、链名称、原生币、RPC、区块浏览器和启用充币的合约。
- `wallet_switchEthereumChain` 只需要 `chainId`；`wallet_addEthereumChain` 才需要 `chainName`、`nativeCurrency`、`rpcUrls`、`blockExplorerUrls`。

扫块消费者会把平台收款地址的 Transfer 事件写入 `hg_addon_chain_recharge`，并按 `(hash, log_index)` 幂等，然后推送 `chain_payment` 队列给业务插件。

需要登录的提币接口：

- `GET /api/chain/withdrawal/options`：可提币网络、币种、手续费、最小/最大数量。
- `POST /api/chain/withdrawal/apply`：提交提币申请，参数 `chain`、`symbol`、`toAddress`、`amount`。
- `GET /api/chain/withdrawal/list`：当前用户提币记录，支持 `page/pageSize`。

提币申请会冻结用户对应的资产，后台审核通过后由打款流程结算。需要在后台合约配置中开启 `withdraw_switch`。

### RPC 要求

扫块依赖 `eth_getLogs`。部分公共测试网 RPC（例如 BSC 测试网的 `data-seed-prebsc-1-s1.binance.org:8545`）会直接返回 `limit exceeded`，即使只查询一个区块。

BSC 测试网可使用：

```text
https://bsc-testnet.publicnode.com
```

生产环境建议使用支持 `eth_getLogs` 的付费/专属 RPC，避免公共节点限流。

### 常用命令行

```shell
# 接口维护-gen service
gf gen service -s=addons/chain/logic -d=addons/chain/service

```
