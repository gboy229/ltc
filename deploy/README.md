# 部署说明（HotGo + 后台 web + h5）

> 目标：让你在任何一台 Linux 服务器上，用最少步骤跑起来，并且知道每一步在干什么。

## 0. 一图看懂

```
构建机（你的 Mac / GitHub Actions）                 服务器（Linux）
──────────────────────────────────               ─────────────────────────────
web/  ──pnpm build──▶  server/resource/public/<VITE_PUBLIC_PATH>
                              │
                              ▼
                     gf build（把 resource 打进二进制）
                              │
                     server/temp/hotgo  ──上传──▶   /opt/ltc/main
                                                    /opt/ltc/manifest/config/config.yaml
                                                    MySQL + Redis
                                                    systemd 启动
                                                    nginx 反代 80/443
```

**部署只有三样东西：**

| 东西 | 说明 | 放哪 |
|---|---|---|
| 二进制 | 已经包含后台页面，不用单独传 dist | `/opt/ltc/main` |
| 配置 | `config.yaml`，数据库/Redis/密钥 | `/opt/ltc/manifest/config/` |
| 依赖 | MySQL、Redis、nginx | 服务器上装 |

> 后台是 Hash 路由 + 页面前缀（当前 `/mMe6PXqfRLhjWCiK`），所以 Go 服务能直接托管，**不需要 nginx 做 SPA 回退**。
> h5 是 History 路由，**必须**由 nginx 做回退。

---

## 1. 服务器一次性准备（约 20 分钟）

以下以 Ubuntu 为例，CentOS 把 `apt` 换成 `dnf`、`www-data` 换成 `nginx`。

### 1.1 装依赖

```bash
sudo apt update
sudo apt install -y mysql-server redis-server nginx rsync
sudo systemctl enable --now mysql redis-server nginx
```

### 1.2 建库

```bash
sudo mysql -e "CREATE DATABASE hotgo DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
sudo mysql -e "CREATE USER 'hotgo'@'127.0.0.1' IDENTIFIED BY '换成强密码';"
sudo mysql -e "GRANT ALL PRIVILEGES ON hotgo.* TO 'hotgo'@'127.0.0.1'; FLUSH PRIVILEGES;"
```

Redis 建议设密码（可选）：

```bash
sudo sed -i 's/^# requirepass .*/requirepass 换成强密码/' /etc/redis/redis.conf
sudo systemctl restart redis-server
```

### 1.3 建目录 + 配置文件

```bash
sudo mkdir -p /opt/ltc/{manifest/config,resource/public/attachment,logs,storage}
sudo mkdir -p /opt/www/h5
sudo chown -R www-data:www-data /opt/ltc /opt/www/h5
```

第一次先把模板传上去：

```bash
# 在项目根目录执行
scp -r server/manifest root@你的服务器IP:/opt/ltc/
ssh root@你的服务器IP "chown -R www-data:www-data /opt/ltc/manifest"
```

到服务器上生成配置：

```bash
cd /opt/ltc
sudo -u www-data cp manifest/config/config.example.yaml manifest/config/config.yaml
sudo -u www-data vim manifest/config/config.yaml
```

**至少改这几项：**

| 配置项 | 改成 |
|---|---|
| `database.default.link` | `mysql:hotgo:密码@tcp(127.0.0.1:3306)/hotgo?loc=Local&parseTime=true&charset=utf8mb4` |
| `database.default.debug` | `false` |
| `redis.default.address` / `pass` | 你的 Redis 地址和密码 |
| `system.debug` | `false` |
| `system.mode` | `"product"` |
| `token.secretKey` | 随机长字符串 |
| `tcp.client.cron.secretKey`、`tcp.client.auth.secretKey` | 随机字符串 |
| `chain.privateKeySecret` | 删除，改用环境变量 `CHAIN_PRIVATE_KEY_SECRET` |
| `users.wallet.skipVerify` | `false` |
| `server.pprofEnabled` | `false` |
| `server.openapiPath` / `swaggerPath` | 置空（或由 nginx 限制来源） |

> 不想装 Redis：把 `cache.adapter` 改成 `file`，`queue.driver` 改成 `disk`。
> 但单机也建议装 Redis，配置最省心。

### 1.4 安装 systemd

```bash
# 在本地项目根目录执行（把文件传上去）
scp deploy/ltc.service root@你的服务器IP:/etc/systemd/system/ltc.service

# 然后在服务器上执行
sudo systemctl daemon-reload
sudo systemctl enable ltc        # 只设置开机自启，先不启动
```

敏感变量写到 `/opt/ltc/.env`：

```bash
sudo -u www-data tee /opt/ltc/.env >/dev/null <<'EOF'
CHAIN_PRIVATE_KEY_SECRET=换成32位以上随机字符串
EOF
```

### 1.5 安装 nginx

```bash
# 在本地项目根目录执行
scp deploy/nginx.conf root@你的服务器IP:/etc/nginx/conf.d/ltc.conf

# 然后在服务器上执行：改 server_name、root 路径
sudo vim /etc/nginx/conf.d/ltc.conf
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl reload nginx
```

如果暂时不部署 h5，把 `nginx.conf` 里的 `root /opt/www/h5;` 和 `location / { try_files ... }` 换成 `return 404;`。

---

## 2. 发版（每次）

### 2.1 本机构建

```bash
./deploy/build.sh
```

它做三件事：

1. `web` 安装依赖并构建，产物拷到 `server/resource/public/<页面前缀>`
2. `gf build` 把 `resource` 打进二进制
3. 输出 `server/temp/hotgo`（linux/amd64，约 60MB；已用 `-s -w` 裁剪，部署时重命名为 `main`）

> 构建机需要：Node 20.19+、pnpm 10、Go 1.25+、gf CLI。
> 没装 gf：`cd server && make cli`。
>
> 耗时参考：前端约 15 秒，后端首次约 1~3 分钟（首次会下载 Go 1.25 工具链），之后约 30 秒。
> 本地已有 `node_modules` 时会跳过安装；需要重装用 `FORCE_INSTALL=1 ./deploy/build.sh`。

构建成功后会输出：

```text
==> 构建完成
产物：.../server/temp/hotgo
大小：约 60M
...: ELF 64-bit LSB executable, x86-64 ...
```

> **这个产物只能在 Linux 上运行。** 在 macOS 直接 `./hotgo` 会报
> `zsh: exec format error`，因为它是 Linux ELF，不是 Mac 可执行文件。

### 2.1.1 本地预览

想在 Mac 上直接跑打好的包，用 `--local` 打一个当前系统的版本：

```bash
./deploy/build.sh --local
./deploy/run-local.sh
```

- 产物是 `server/temp/hotgo-local`（在当前机器上可直接执行）；
- 默认监听 `:8000`，如果 `gf run` 正在跑会端口冲突，先停掉其中一个；
- **必须从 `server` 目录启动**：`manifest/config`、`resource`、`logs` 都是相对工作目录的路径。
  在 `server/temp` 里直接 `./hotgo-local` 会报
  `casbin model file does not exist: manifest/config/casbin.conf`，进程活着但 8000 不监听。
  `run-local.sh` 已经帮你 `cd` 到 `server` 再启动。

访问地址：<http://127.0.0.1:8000/mMe6PXqfRLhjWCiK>（根路径 `/` 返回 403 是正常的）。

### 2.2 上传

首次（先不重启，因为数据库还没初始化）：

```bash
./deploy/deploy.sh root@你的服务器IP --no-restart
```

之后每次：

```bash
./deploy/deploy.sh root@你的服务器IP
```

脚本会：传二进制 → 同步 `manifest`（**不覆盖服务器上的 config.yaml**）→ 切换 → 重启 → `curl /<前缀>` 健康检查。

### 2.3 首次初始化（只做一次）

```bash
ssh root@你的服务器IP
cd /opt/ltc

sudo -u www-data ./main migrate    # 建表
sudo -u www-data ./main seed       # 种子数据（默认管理员 admin / 123456）
sudo -u www-data ./main menu       # 同步菜单

sudo systemctl start ltc
```

**登录后第一件事：改掉 admin 的密码。**

### 2.4 验证

```bash
curl -I http://你的服务器IP/mMe6PXqfRLhjWCiK/          # 期望 200
sudo systemctl status ltc
sudo tail -f /opt/ltc/logs/server/*.log
```

### 2.5 h5（独立前端）

```bash
cd h5
pnpm install
pnpm build
rsync -az dist/ root@你的服务器IP:/opt/www/h5/
```

h5 的接口地址是 `/api`，由 nginx 反代到 Go，和后台同域，不需要改跨域配置。

> 以上是手工方式；push 后 CI 会自动构建并发布 h5（见第 3 节），不需要再手动 rsync。

### 2.6 修改后台页面路径

当前页面路径由 `web/.env.production` 的 `VITE_PUBLIC_PATH` 决定（现在是 `/mMe6PXqfRLhjWCiK`），API 固定 `/admin`。两者是分开的：

- `VITE_PUBLIC_PATH`：页面和 JS/CSS 的路径；
- `VITE_GLOB_API_URL_PREFIX`：接口前缀，保持 `/admin`，不要跟着页面前缀改。

**方式一：重新打包（页面和资源一起搬）**

1. 改 `web/.env.production` 的 `VITE_PUBLIC_PATH=/新路径`；
2. `./deploy/build.sh`：脚本会自动把 dist 拷到 `resource/public/新路径`，并清理旧的 `resource/public/admin`；
3. 上传二进制、重启服务；
4. 把 [deploy/nginx.conf](nginx.conf) 里 `[1]` 的 `/mMe6PXqfRLhjWCiK` 换成新路径，reload。

**方式二：不改打包，只换入口（临时换路径）**

在 [deploy/nginx.conf](nginx.conf) 里启用注释的 `/secret` 段，改成想要的路径：

```nginx
location = /secret {
    return 302 /secret/;
}
location = /secret/ {
    proxy_pass http://127.0.0.1:8000/mMe6PXqfRLhjWCiK/index.html;
    ...
}
```

然后 `sudo nginx -t && sudo systemctl reload nginx`，访问 `http://域名/secret/`。

页面 HTML 仍从当前构建前缀（`VITE_PUBLIC_PATH`）的 `index.html` 出，JS/CSS 也在该前缀下，API 仍走 `/admin`，所以不用重新打包、也不用重启 Go。

注意：`config.yaml` 的 `router.admin.prefix` 保持 `/admin`，不要和页面前缀混在一起。

---

## 3. GitHub Actions（push 自动部署）

### 3.1 触发规则：改哪个目录，就部署哪个

| 改动范围 | 触发的工作流 | 做什么 |
|---|---|---|
| `server/**`、`web/**` | `deploy-server.yml` | 构建后台前端 + 二进制 → 上传 → 切换 → 重启 → 健康检查 |
| `h5/**` | `deploy-h5.yml` | 构建 h5 → 上传 → 替换 `/opt/www/h5`（**不重启 Go**） |
| 只改 `deploy/` 下脚本或 `*.md` | 不触发 | 需要时在 Actions 页面手动 `Run workflow` |
| 打 tag `v*` | `release.yml` | 构建并创建 Release，**不部署** |

> `web/**` 必须触发服务端部署：后台页面会被 `gf build` 打进 Go 二进制，只传 dist 不生效。
> 一次 push 同时改了 `web/` 和 `h5/`，两个工作流会并行跑，互不影响。

工作流文件：

| 文件 | 作用 |
|---|---|
| `.github/workflows/build-server.yml` | 构建 + 可选部署/发布（不产生 artifact） |
| `.github/workflows/deploy-server.yml` | push 自动部署 Go 服务 |
| `.github/workflows/deploy-h5.yml` | push 自动部署 h5 |
| `.github/workflows/release.yml` | 打 tag 发 Release |

### 3.2 Actions 存储：为什么不占 500MB

GitHub Free 的 **制品（artifact）额度只有 500MB**，而且和 Packages 共享，是**账号级**的。
一次服务端构建的压缩包约 38MB，如果每次 push 都存 artifact、保留 14 天，十几次就会把额度占满，
之后所有仓库的上传都会报 `Artifact storage quota has been hit`。

所以本项目的做法是：

- **不生成 artifact**：构建、打包、scp、重启都在同一个 job 里完成，二进制只在 runner 的临时磁盘上存在；
- h5 同理，构建完直接发布，不存 artifact；
- pnpm / Go 的缓存走**独立的 10GB/仓库缓存额度**，自动淘汰，不占 500MB 制品额度，可以放心用。

> 以后如果有人加 `actions/upload-artifact`，必须写 `retention-days: 1`，否则很容易把额度打满。

### 3.3 用哪个用户部署

**推荐 root**：本项目是单机部署，部署动作要同时写 `/opt/ltc`、`/opt/www/h5` 并执行 `systemctl`，用普通用户要配三处权限，收益有限。
务必使用**专用密钥**（只给这个仓库、可随时吊销），不要复用个人密钥。

如果 `SSH_USER` 不是 root（例如云主机默认的 `ubuntu`），该用户必须能**免密 sudo**：脚本会用 `sudo` 写 `/opt/ltc`、`/opt/www/h5` 并执行 `systemctl`。

```bash
# 在服务器上确认（有输出就说明可以）
sudo -n true && echo "有免密 sudo"
```

没有免密 sudo 的话，二选一：

- 改用 `SSH_USER=root`（推荐，最省事）；
- 给部署用户配置免密 sudo（注意这等价于 root 权限）：

```bash
echo 'ubuntu ALL=(ALL) NOPASSWD: ALL' | sudo tee /etc/sudoers.d/ltc-deploy
sudo chmod 440 /etc/sudoers.d/ltc-deploy
```

> 部署脚本会自动读取 `ltc.service` 的 `User=` 和 nginx 的 `user` 指令，把 `/opt/ltc/main`、`/opt/www/h5`
> 的属主改成实际运行用户。服务器用 `nginx` / `www` 时，记得同步改 `ltc.service` 的 `User=` / `Group=`。

### 3.4 Secrets / Variables

仓库 **Settings → Secrets and variables → Actions**：

| 类型 | 名称 | 必填 | 说明 |
|---|---|---|---|
| Secret | `SSH_HOST` | ✅ | 服务器 IP/域名 |
| Secret | `SSH_USER` | ✅ | 部署用户（root 或 deployer） |
| Secret | `SSH_KEY` | ✅ | 私钥全文（建议 ed25519） |
| Secret | `SSH_PORT` | 否 | 默认 22 |
| Variable | `H5_CHECK_URL` | 否 | h5 健康检查地址，默认 `http://127.0.0.1/` |
| Variable | `HEALTH_RETRIES` | 否 | 服务端健康检查等待秒数，默认 `60` |
| Variable | `STRICT_HEALTH_CHECK` | 否 | `1` = 检查失败自动回滚；默认 `0` 只告警、保留新版本 |

> 服务器上有多个 vhost 时，把 `H5_CHECK_URL` 设为真实域名（如 `https://example.com/`），否则健康检查可能命中别的站点。

### 3.5 服务器一次性准备

CI 只负责「传包 + 重启」，下面这些仍要手动做一次：

1. 目录：`/opt/ltc/{manifest/config,resource/public/attachment,logs,storage}`、`/opt/www/h5`；
2. 配置：按第 1.3 节准备 `/opt/ltc/manifest/config/config.yaml`（CI 永远不会覆盖它）；
3. 服务：按第 1.4 节安装 `ltc.service`，按第 1.5 节配置 nginx（h5 的 `root /opt/www/h5`）；
4. 首次初始化（CI 不会替你做）：

```bash
cd /opt/ltc
sudo -u www-data ./main migrate
sudo -u www-data ./main seed     # 只执行一次
sudo -u www-data ./main menu
sudo systemctl enable --now ltc
```

**首次上线前先手动跑通一次**，确认后台能打开、能登录，再让 CI 接管。登录后第一件事：改掉 admin 默认密码。

### 3.6 日常发版

```bash
git add -A
git commit -m "feat: xxx"
git push origin master
```

- 改了 `server/` 或 `web/`：约 4~8 分钟上线（构建 3~5 分钟 + 传输/重启 1 分钟）；
- 改了 `h5/`：约 1~2 分钟上线，Go 服务不重启；
- 改了文档：不触发，需要时在 Actions 页面手动 `Run workflow`。

### 3.7 数据库迁移（默认手动）

CI **不会**执行 `migrate`，只更新二进制并重启。版本包含表结构变更时，手动执行：

```bash
ssh root@服务器
cd /opt/ltc
sudo -u www-data ./main migrate
sudo -u www-data ./main menu        # 菜单有变动时
```

也可以在 Actions → deploy-server → Run workflow 时勾选 `migrate`。
`migrate` 有 `hg_schema_migrations` 版本表，重复执行安全；`seed` 只在首次执行。

### 3.8 回滚

- 默认：健康检查失败只告警，**保留新版本**，不会自动回滚；
- 想要自动回滚：把仓库 Variable `STRICT_HEALTH_CHECK` 设为 `1`；
- 手动：

```bash
cd /opt/ltc
cp -f main.prev main
systemctl restart ltc
```

- h5：发布前会备份到 `/tmp/h5-backup.*`，健康检查失败自动回滚；手动回滚用上一版 `dist` 重新发布即可。

### 3.9 手工和 CI 怎么选

- 本地调试、紧急修服务器：继续用 `./deploy/build.sh` + `./deploy/deploy.sh`；
- 日常发版：`git push`，CI 自动完成；
- 两边共用同一套服务器端逻辑：[remote-deploy-server.sh](remote-deploy-server.sh)。

## 4. 常见问题

### 4.1 前端改了，页面没变

因为 `gf build` 把 `resource` 打包进了二进制，而 GoFrame 的查找顺序是 **打包资源 > 磁盘**。
**必须重新执行 `./deploy/build.sh` 并替换二进制**，只传 dist 无效。

确认方法：

```bash
grep -o 'assets/index-[^"]*\.js' web/dist/index.html
curl -s http://127.0.0.1:8000/mMe6PXqfRLhjWCiK/index.html | grep -o 'assets/index-[^"]*\.js'
```

两个 hash 必须一致。

### 4.2 后台页面返回 404

- 构建时用了 `go build` 而不是 `gf build`，资源没打进包；
- 或者二进制里没有 `resource/public/<页面前缀>`（构建前没拷 dist）。

### 4.3 启动报数据库/Redis 连接失败

- 检查 `config.yaml` 的账号密码、库名；
- 确认 MySQL/Redis 已启动、监听 127.0.0.1；
- 不想用 Redis：`cache.adapter: file` + `queue.driver: disk`。

### 4.4 日志/上传目录没权限

```bash
sudo chown -R www-data:www-data /opt/ltc/{logs,storage,resource}
```

### 4.5 执行命令提示找不到 config

所有相对路径（`manifest/...`、`logs/...`、`resource/...`）都以**进程工作目录**为准。
必须在 `/opt/ltc` 下执行，systemd 里已经用 `WorkingDirectory` 固定。

### 4.6 h5 刷新页面 404

h5 是 History 路由，必须在 nginx 配 `try_files $uri $uri/ /index.html;`。

> h5 部署会保留 `.user.ini`、`.well-known`、`.htaccess`（宝塔等面板生成，`.user.ini` 常带 `chattr +i`），
> 不会因为 `rsync --delete` 删不掉而失败。

### 4.7 自动部署失败

先看 Actions 里失败的是哪一步：

- **SSH/SCP 失败**：本机先手动 `ssh` 一次，确认密钥和端口；服务器安全组放行 SSH 端口；
- **健康检查失败**：看日志 `sudo journalctl -u ltc -n 100 --no-pager`，以及 `/opt/ltc/logs/server/*.log`；
  默认会保留新版本（只告警），修好后重新 push；想自动回滚就设 `STRICT_HEALTH_CHECK=1`；
- **缺少配置**：确认 `/opt/ltc/manifest/config/config.yaml` 已存在、`ltc.service` 已安装；
- **h5 健康检查失败**：确认 nginx 已配置 `root /opt/www/h5`，必要时设置变量 `H5_CHECK_URL`；
- **迁移没跑**：CI 默认只更新二进制，表结构变更要手动 `cd /opt/ltc && sudo -u www-data ./main migrate`。

### 4.8 服务启动失败：status=217/USER

`ltc.service` 里的 `User=www-data` 在服务器上不存在（CentOS / 宝塔面板常见，用户是 `nginx` 或 `www`）。

```bash
id nginx
id www
```

把 `/etc/systemd/system/ltc.service` 里的 `User=` / `Group=` 改成实际用户，然后：

```bash
sudo chown -R <用户>:<用户> /opt/ltc/{logs,storage,resource} /opt/ltc/main
sudo systemctl daemon-reload
sudo systemctl restart ltc
```

### 4.9 编译报 Go 缓存目录没有权限

错误类似：

```text
open /Users/xxx/Library/Caches/go-build/...: operation not permitted
```

说明当前环境不允许写默认缓存目录（沙箱、受限容器等）。把缓存指到项目内：

```bash
GOCACHE="$PWD/.cache/go-build" ./deploy/build.sh
```

`.cache/` 已加入 `.gitignore`。普通本机终端直接 `./deploy/build.sh` 即可，不需要这一步。

### 4.10 在线安装插件成功，但没有建表

插件建表的 SQL 已经通过 `go:embed` 打进二进制（[server/addons/embed.go](../server/addons/embed.go)），程序优先从内嵌文件读取，不再依赖服务器上的 `addons/` 目录。

如果仍然没建表：

1. 确认跑的是新构建的二进制（push 后 CI 会重新打包）；
2. 到后台把插件**卸载再安装**，迁移才会重新执行；
3. 旧部署可能在服务器留有 `/opt/ltc/addons`，可以删掉（新代码优先用内嵌 SQL）。

> 插件建表 SQL 使用 `CREATE TABLE IF NOT EXISTS`，重复执行安全。
> 新增插件或修改 SQL 后，生产环境必须重新构建二进制才会生效。
> 开发模式（`system.mode: develop`）下如果仓库有 `addons/` 目录，会优先读磁盘，改 SQL 立即生效；生产模式（`product`）读内嵌 SQL。

### 4.11 报 `BUILD_LABEL: unbound variable`

macOS 自带的是 bash 3.2，在 UTF-8 locale 下，`$VAR` 后面紧跟中文全角标点（如 `（`、`）`、`，`）时，会把标点字节当成变量名的一部分，于是报 `unbound variable`。

脚本里已经统一改成 `${VAR}` 写法。以后自己加代码时，变量后面接中文也请用 `${VAR}`。

---

## 5. 相关文件

| 文件 | 用途 |
|---|---|
| `deploy/build.sh` | 一键构建前端 + 二进制 |
| `deploy/run-local.sh` | 本机启动预览版（自动 cd 到 server） |
| `deploy/deploy.sh` | 手工上传 + 切换 + 重启 + 健康检查 |
| `deploy/remote-deploy-server.sh` | 服务器端部署脚本（CI 与手工共用，含回滚） |
| `deploy/remote-deploy-h5.sh` | 服务器端 h5 发布脚本（含回滚） |
| `deploy/ltc.service` | systemd 单元 |
| `deploy/nginx.conf` | nginx 反代示例（后台 + h5 + WebSocket） |
| `.github/workflows/build-server.yml` | GitHub Actions：构建 + 可选部署/发布（无 artifact） |
| `.github/workflows/deploy-server.yml` | GitHub Actions：push 自动部署 Go 服务 |
| `.github/workflows/deploy-h5.yml` | GitHub Actions：push 自动部署 h5 |
| `.github/workflows/release.yml` | GitHub Actions：打 tag 生成 Release |
