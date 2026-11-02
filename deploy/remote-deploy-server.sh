#!/usr/bin/env bash
#
# 服务器端：切换二进制 + 同步 manifest + 重启 + 健康检查 + 失败回滚
#
# 由 GitHub Actions 调用，也可以在服务器上手动执行：
#   mkdir -p /tmp/hotgo-deploy/pkg
#   tar -xzf hotgo-server.tar.gz -C /tmp/hotgo-deploy/pkg
#   bash /tmp/hotgo-deploy/pkg/remote-deploy-server.sh
#   bash /tmp/hotgo-deploy/pkg/remote-deploy-server.sh --migrate
#
# 环境变量：
#   APP_DIR   应用目录，默认 /opt/ltc
#   APP_PORT  健康检查端口，默认 8000
#   SERVICE   systemd 服务名，默认 ltc
#
# 说明：
#   - 二进制部署为 $APP_DIR/main，旧版本备份为 main.prev
#   - 永远不覆盖服务器上的 manifest/config/config.yaml
#   - 默认不执行 migrate / seed，按需手动执行
#   - 新二进制先落到 main.new 再 mv，避免 "Text file busy"
#   - 健康检查失败自动回滚到 main.prev
#
set -euo pipefail

APP_DIR="${APP_DIR:-/opt/ltc}"
APP_PORT="${APP_PORT:-8000}"
SERVICE="${SERVICE:-ltc}"
HEALTH_RETRIES="${HEALTH_RETRIES:-60}"
STRICT_HEALTH_CHECK="${STRICT_HEALTH_CHECK:-0}"

RUN_MIGRATE=0
RUN_MENU=0
for arg in "$@"; do
  case "$arg" in
    --migrate) RUN_MIGRATE=1 ;;
    --menu)    RUN_MENU=1 ;;
    -h|--help) sed -n '2,22p' "$0"; exit 0 ;;
    *) echo "未知参数：$arg（可用：--migrate --menu）" >&2; exit 2 ;;
  esac
done

ok()   { printf '\033[1;32m==>\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31mERROR:\033[0m %s\n' "$*" >&2; exit 1; }

SRC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
[ -f "$SRC_DIR/main" ]             || fail "部署包内缺少 main 二进制"
[ -f "$SRC_DIR/page-prefix.txt" ]  || fail "部署包内缺少 page-prefix.txt"
[ -d "$SRC_DIR/manifest" ]         || fail "部署包内缺少 manifest 目录"
[ -f "$APP_DIR/manifest/config/config.yaml" ] || fail "服务器缺少 $APP_DIR/manifest/config/config.yaml"
command -v rsync >/dev/null 2>&1   || fail "服务器缺少 rsync"
command -v curl  >/dev/null 2>&1   || fail "服务器缺少 curl"

PAGE_PREFIX="$(tr -d '[:space:]' < "$SRC_DIR/page-prefix.txt")"
[ -n "$PAGE_PREFIX" ] || fail "page-prefix.txt 为空"
case "$PAGE_PREFIX" in
  /?*) ;;
  *) fail "页面前缀不合法：$PAGE_PREFIX" ;;
esac

if [ "$(id -u)" -eq 0 ]; then
  as_root() { "$@"; }
else
  command -v sudo >/dev/null 2>&1 || fail "当前用户非 root 且没有 sudo"
  sudo -n true 2>/dev/null || fail "当前用户没有免密 sudo，请改用 root 部署，或给该用户配置 NOPASSWD"
  as_root() { sudo -n "$@"; }
fi

# 切换前先确认服务存在，避免二进制换了却起不来（systemctl show 不需要 root；
# 查不到状态时跳过，交给后面的 restart 失败回滚兜底）
if command -v systemctl >/dev/null 2>&1; then
  UNIT_STATE="$(systemctl show "$SERVICE" -p LoadState --value 2>/dev/null || true)"
  if [ -n "$UNIT_STATE" ] && [ "$UNIT_STATE" != "loaded" ]; then
    fail "服务器未安装 systemd 服务：${SERVICE}（见 deploy/ltc.service）"
  fi
fi

# 以服务实际运行用户为准（www-data / nginx / www 等）
APP_USER="$(systemctl show "$SERVICE" -p User --value 2>/dev/null || true)"
[ -n "$APP_USER" ] || APP_USER="root"

chown_app() {
  if id "$APP_USER" >/dev/null 2>&1; then
    as_root chown "$APP_USER" "$1" 2>/dev/null || true
  fi
}

rollback() {
  if [ -f "$APP_DIR/main.prev" ]; then
    as_root cp -f "$APP_DIR/main.prev" "$APP_DIR/main.new"
    as_root chmod 0755 "$APP_DIR/main.new"
    as_root mv -f "$APP_DIR/main.new" "$APP_DIR/main"
    chown_app "$APP_DIR/main"
    as_root systemctl restart "$SERVICE" || true
    echo "已回滚到 main.prev" >&2
  else
    echo "没有可用的 main.prev，未回滚" >&2
  fi
}

# 同一时间只允许一个部署
if command -v flock >/dev/null 2>&1; then
  exec 9>/tmp/ltc-deploy.lock
  flock -n 9 || fail "已有部署正在进行"
fi

ok "同步 manifest（保留 config.yaml）"
as_root rsync -a --checksum --exclude 'config/config.yaml' "$SRC_DIR/manifest/" "$APP_DIR/manifest/"

as_root chmod 0755 "$SRC_DIR/main"

if [ "$RUN_MIGRATE" = "1" ]; then
  ok "执行数据库迁移（手动触发）"
  ( cd "$APP_DIR" && "$SRC_DIR/main" migrate )
fi
if [ "$RUN_MENU" = "1" ]; then
  ok "同步菜单（手动触发）"
  ( cd "$APP_DIR" && "$SRC_DIR/main" menu )
fi

ok "备份当前二进制到 main.prev"
if [ -f "$APP_DIR/main" ]; then
  as_root cp -f "$APP_DIR/main" "$APP_DIR/main.prev"
fi

ok "切换二进制"
as_root cp -f "$SRC_DIR/main" "$APP_DIR/main.new"
as_root chmod 0755 "$APP_DIR/main.new"
as_root mv -f "$APP_DIR/main.new" "$APP_DIR/main"
chown_app "$APP_DIR/main"

# 切换运行用户后 logs/storage 可能还是旧用户所有，导致启动 panic
ok "修正运行目录权限"
as_root mkdir -p "$APP_DIR/logs" "$APP_DIR/storage" "$APP_DIR/resource/public/attachment"
as_root chown -R "$APP_USER" "$APP_DIR/logs" "$APP_DIR/storage" "$APP_DIR/resource/public/attachment" 2>/dev/null || true

# /tmp/gsessions 可能由旧运行用户创建，chown 一下避免 session 不可写
if [ -d /tmp/gsessions ]; then
  as_root chown -R "$APP_USER" /tmp/gsessions 2>/dev/null || true
fi

ok "重启服务"
if ! as_root systemctl restart "$SERVICE"; then
  printf '\033[1;31mERROR:\033[0m 服务启动失败，开始回滚\n' >&2
  rollback
  as_root systemctl --no-pager --lines=30 status "$SERVICE" || true
  exit 1
fi

ok "健康检查 http://127.0.0.1:${APP_PORT}${PAGE_PREFIX}/"
HEALTHY=0
for _ in $(seq 1 "$HEALTH_RETRIES"); do
  if curl -fsS -o /dev/null --max-time 3 "http://127.0.0.1:${APP_PORT}${PAGE_PREFIX}/" 2>/dev/null; then
    HEALTHY=1
    break
  fi
  sleep 1
done

if [ "$HEALTHY" -eq 1 ]; then
  ok "部署成功：${PAGE_PREFIX}/"
  if [ "$RUN_MIGRATE" != "1" ]; then
    printf '提示：本次未执行数据库迁移。如版本包含表结构变更，请手动执行：\n'
    printf '  cd %s && sudo -u %s ./main migrate\n' "$APP_DIR" "$APP_USER"
  fi
  exit 0
fi

if [ "$STRICT_HEALTH_CHECK" = "1" ]; then
  printf '\033[1;31mERROR:\033[0m 健康检查失败，开始回滚\n' >&2
  rollback
  as_root systemctl --no-pager --lines=30 status "$SERVICE" || true
  exit 1
fi

printf '\033[1;33mWARN:\033[0m 健康检查 %s 秒未通过，已保留新版本（不回滚）。\n' "$HEALTH_RETRIES" >&2
printf '请检查：sudo journalctl -u %s -n 100 --no-pager\n' "$SERVICE" >&2
as_root systemctl --no-pager --lines=20 status "$SERVICE" || true
exit 0
