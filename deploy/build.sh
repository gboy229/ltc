#!/usr/bin/env bash
#
# 一键构建：后台前端 -> server/resource/public/<VITE_PUBLIC_PATH> -> 单文件二进制
#
# 用法：
#   ./deploy/build.sh            # 打服务器用的 linux/amd64
#   ./deploy/build.sh --local    # 打本机可执行版本，用于本地预览
#
# 产物：
#   server/temp/hotgo            linux/amd64（部署到服务器）
#   server/temp/hotgo-local      当前系统/架构（本地直接运行）
#
# 依赖：
#   node >= 20.19、pnpm 10、Go 1.25+、gf CLI
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEB_DIR="$ROOT/web"
SERVER_DIR="$ROOT/server"

ok()   { printf '\033[1;32m==>\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31mERROR:\033[0m %s\n' "$*" >&2; exit 1; }

LOCAL_BUILD=0
for arg in "$@"; do
  case "$arg" in
    --local) LOCAL_BUILD=1 ;;
    -h|--help) sed -n '2,15p' "$0"; exit 0 ;;
    *) fail "未知参数：${arg}（可用：--local）" ;;
  esac
done

if [[ "$LOCAL_BUILD" == "1" ]]; then
  OUT_BIN="$SERVER_DIR/temp/hotgo-local"
  BUILD_LABEL="本地预览 $(go env GOOS)/$(go env GOARCH)"
else
  OUT_BIN="$SERVER_DIR/temp/hotgo"
  BUILD_LABEL="linux/amd64（服务器）"
fi

command -v pnpm >/dev/null 2>&1 || fail "缺少 pnpm，请先安装：npm i -g pnpm"
command -v gf   >/dev/null 2>&1 || fail "缺少 gf CLI，可在 server 目录执行 make cli 安装"
command -v go   >/dev/null 2>&1 || fail "缺少 Go，需要 1.25+（go.mod 要求）"

# 页面前缀从 .env.production 的 VITE_PUBLIC_PATH 读取，
# 保证「构建出来的资源路径」和「拷贝目录」一致。
PAGE_PREFIX="$(sed -n 's/^[[:space:]]*VITE_PUBLIC_PATH=//p' "$WEB_DIR/.env.production" | tail -1 | tr -d '\r')"
PAGE_PREFIX="${PAGE_PREFIX:-/admin}"
PAGE_PREFIX="/${PAGE_PREFIX#/}"
PAGE_PREFIX="${PAGE_PREFIX%/}"
[[ "$PAGE_PREFIX" =~ ^/[A-Za-z0-9_-]+$ ]] || fail "VITE_PUBLIC_PATH 不合法：${PAGE_PREFIX}"
ADMIN_DIR="$SERVER_DIR/resource/public${PAGE_PREFIX}"

# 安全检查：只允许操作 resource/public 下的子目录
[[ "$ADMIN_DIR" == "$SERVER_DIR/resource/public/"* ]] || fail "静态目录异常：$ADMIN_DIR"

ok "页面前缀：${PAGE_PREFIX}"
ok "1/3 安装并构建后台前端"
cd "$WEB_DIR"
# 本地已有依赖时跳过安装：pnpm 可能因 store 路径变化想重建 node_modules，
# 非交互环境会直接报错。CI 上 node_modules 不存在，会正常安装。
if [[ "${FORCE_INSTALL:-0}" == "1" || ! -d "$WEB_DIR/node_modules" ]]; then
  pnpm install --frozen-lockfile
else
  ok "检测到 node_modules，跳过安装（需要重装请用 FORCE_INSTALL=1）"
fi
pnpm build

[[ -f "$WEB_DIR/dist/index.html" ]] || fail "前端构建失败：缺少 web/dist/index.html"

ok "2/3 同步前端到 server/resource/public${PAGE_PREFIX}"
# 换过前缀时清掉旧目录，避免旧页面还能访问
if [[ "$PAGE_PREFIX" != "/admin" && -d "$SERVER_DIR/resource/public/admin" ]]; then
  ok "清理旧页面目录 resource/public/admin"
  rm -rf "$SERVER_DIR/resource/public/admin"
fi
rm -rf "$ADMIN_DIR"
mkdir -p "$ADMIN_DIR"
cp -R "$WEB_DIR/dist/." "$ADMIN_DIR/"

ok "3/3 编译并把资源打进二进制（${BUILD_LABEL}）"
cd "$SERVER_DIR"
# gf build 会读取 hack/config.yaml，并在覆盖 packed.go 时询问，自动回答 y
if [[ "$LOCAL_BUILD" == "1" ]]; then
  printf 'y\n' | gf build -s "$(go env GOOS)" -a "$(go env GOARCH)" -o "$OUT_BIN"
else
  printf 'y\n' | gf build
fi

[[ -f "$OUT_BIN" ]] || fail "未找到构建产物：$OUT_BIN"

ok "构建完成"
printf '产物：%s\n' "$OUT_BIN"
printf '大小：%s\n' "$(du -h "$OUT_BIN" | awk '{print $1}')"
file "$OUT_BIN" 2>/dev/null || true

if [[ "$LOCAL_BUILD" == "1" ]]; then
  cat <<EOF

本地预览：
  ./deploy/run-local.sh
  注意：默认监听 :8000，如与正在运行的 gf run 冲突，先停掉其中一个。
EOF
else
  cat <<'EOF'

下一步：
  首次部署：./deploy/deploy.sh <user@host> --no-restart
           然后到服务器执行 migrate / seed / menu（见 deploy/README.md）
  日常发版：./deploy/deploy.sh <user@host>
EOF
fi
