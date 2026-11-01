#!/usr/bin/env bash
#
# 服务器端：发布 h5 静态文件到 /opt/www/h5
#
# 由 GitHub Actions 调用，也可以手动执行：
#   mkdir -p /tmp/hotgo-h5-deploy/pkg
#   tar -xzf h5-dist.tar.gz -C /tmp/hotgo-h5-deploy/pkg
#   bash /tmp/hotgo-h5-deploy/pkg/remote-deploy-h5.sh
#
# 环境变量：
#   APP_DIR       h5 发布目录，默认 /opt/www/h5（与 nginx root 一致）
#   H5_CHECK_URL  健康检查地址，默认 http://127.0.0.1/
#
# 说明：
#   - h5 是纯静态文件，不重启 Go 服务
#   - 发布前备份、发布后健康检查，失败自动回滚
#   - 需要服务器上有 rsync、curl
#
set -euo pipefail

APP_DIR="${APP_DIR:-/opt/www/h5}"
CHECK_URL="${H5_CHECK_URL:-http://127.0.0.1/}"

ok()   { printf '\033[1;32m==>\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31mERROR:\033[0m %s\n' "$*" >&2; exit 1; }

SRC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/dist"
[ -f "$SRC_DIR/index.html" ] || fail "部署包内缺少 dist/index.html"
command -v rsync >/dev/null 2>&1 || fail "服务器缺少 rsync"
command -v curl  >/dev/null 2>&1 || fail "服务器缺少 curl"

if [ "$(id -u)" -eq 0 ]; then
  as_root() { "$@"; }
else
  command -v sudo >/dev/null 2>&1 || fail "当前用户非 root 且没有 sudo"
  as_root() { sudo -n "$@"; }
fi

if command -v flock >/dev/null 2>&1; then
  exec 9>/tmp/hotgo-h5-deploy.lock
  flock -n 9 || fail "已有 h5 部署正在进行"
fi

mkdir -p "$APP_DIR"

BACKUP="/tmp/h5-backup.$$"
rm -rf "$BACKUP"
if [ -d "$APP_DIR" ] && [ -n "$(ls -A "$APP_DIR" 2>/dev/null || true)" ]; then
  ok "备份当前 h5 到 $BACKUP"
  cp -a "$APP_DIR" "$BACKUP"
fi

ok "同步 h5 到 $APP_DIR"
rsync -a --delete --checksum "$SRC_DIR/" "$APP_DIR/"
as_root chown -R www-data:www-data "$APP_DIR" 2>/dev/null || true

ok "健康检查 $CHECK_URL"
if curl -fsS -o /dev/null --max-time 5 "$CHECK_URL"; then
  rm -rf "$BACKUP"
  ok "h5 部署成功"
  exit 0
fi

printf '\033[1;31mERROR:\033[0m h5 健康检查失败，开始回滚\n' >&2
if [ -d "$BACKUP" ]; then
  rm -rf "$APP_DIR"
  mv "$BACKUP" "$APP_DIR"
  as_root chown -R www-data:www-data "$APP_DIR" 2>/dev/null || true
  echo "已回滚到上一版 h5" >&2
else
  echo "没有可回滚的备份（首次部署），保留已上传文件，请检查 nginx 配置" >&2
fi
exit 1
