#!/usr/bin/env bash
#
# 上传二进制 + 运行期文件，并在服务器上重启服务（手工部署用，CI 不调用）。
#
# 用法：
#   ./deploy/deploy.sh <user@host> [远程目录] [--no-restart]
#
# 示例：
#   ./deploy/deploy.sh root@1.2.3.4
#   ./deploy/deploy.sh root@1.2.3.4 /opt/ltc
#   ./deploy/deploy.sh root@1.2.3.4 --no-restart   # 首次部署，先不启动
#
# 前置条件（服务器，只做一次）：
#   1. 已安装 ltc.service（见 deploy/ltc.service）
#   2. 已存在 <远程目录>/manifest/config/config.yaml
#
# 说明：
#   - 永远不覆盖服务器上的 manifest/config/config.yaml
#   - 二进制部署为 <远程目录>/main，先传成 main.new，再 mv 覆盖，避免 "Text file busy"
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$ROOT/server/temp/hotgo"

# 页面前缀（与 web/.env.production 保持一致）
PAGE_PREFIX="$(sed -n 's/^[[:space:]]*VITE_PUBLIC_PATH=//p' "$ROOT/web/.env.production" | tail -1 | tr -d '\r')"
PAGE_PREFIX="${PAGE_PREFIX:-/admin}"
PAGE_PREFIX="/${PAGE_PREFIX#/}"
PAGE_PREFIX="${PAGE_PREFIX%/}"

NO_RESTART=0
POSITIONAL=()
for arg in "$@"; do
  case "$arg" in
    --no-restart) NO_RESTART=1 ;;
    *) POSITIONAL+=("$arg") ;;
  esac
done

TARGET="${POSITIONAL[0]:-}"
REMOTE_DIR="${POSITIONAL[1]:-/opt/ltc}"

ok()   { printf '\033[1;32m==>\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31mERROR:\033[0m %s\n' "$*" >&2; exit 1; }

[[ -n "$TARGET" ]] || fail "用法：./deploy/deploy.sh <user@host> [远程目录] [--no-restart]"
[[ -f "$BIN" ]]    || fail "未找到 ${BIN}，请先执行 ./deploy/build.sh"
command -v rsync   >/dev/null 2>&1 || fail "本机缺少 rsync"

ok "创建远程目录 $REMOTE_DIR"
ssh "$TARGET" "mkdir -p '$REMOTE_DIR/manifest/config' '$REMOTE_DIR/resource/public/attachment' '$REMOTE_DIR/logs' '$REMOTE_DIR/storage'"

ok "上传二进制"
scp -q "$BIN" "$TARGET:$REMOTE_DIR/main.new"

ok "同步 manifest（保留 config.yaml）"
rsync -az --exclude 'config/config.yaml' "$ROOT/server/manifest/config/"   "$TARGET:$REMOTE_DIR/manifest/config/"
rsync -az                                "$ROOT/server/manifest/i18n/"     "$TARGET:$REMOTE_DIR/manifest/i18n/"
rsync -az                                "$ROOT/server/manifest/migrations/" "$TARGET:$REMOTE_DIR/manifest/migrations/"
rsync -az                                "$ROOT/server/manifest/seeders/"  "$TARGET:$REMOTE_DIR/manifest/seeders/"

ok "切换二进制"
ssh "$TARGET" "set -e
  cd '$REMOTE_DIR'
  if [ ! -f manifest/config/config.yaml ]; then
    echo '缺少 $REMOTE_DIR/manifest/config/config.yaml，请先从 config.example.yaml 复制并修改' >&2
    rm -f main.new
    exit 1
  fi
  chmod +x main.new
  mv -f main.new main
"

if [[ "$NO_RESTART" == "1" ]]; then
  ok "已上传，按 --no-restart 跳过重启"
  cat <<EOF

首次部署请到服务器执行：
  cd $REMOTE_DIR
  sudo -u www-data ./main migrate && sudo -u www-data ./main seed && sudo -u www-data ./main menu
  systemctl enable --now ltc
EOF
  exit 0
fi

ok "重启服务并做健康检查"
ssh "$TARGET" "set -e
  systemctl restart ltc
  sleep 3
  systemctl --no-pager --lines=5 status ltc || true
  curl -fsS -o /dev/null "http://127.0.0.1:8000${PAGE_PREFIX}/"
  echo "健康检查通过：${PAGE_PREFIX}/ 返回 200"
"
