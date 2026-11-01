#!/usr/bin/env bash
#
# 启动本机预览版二进制（server/temp/hotgo-local）
#
# 用法：
#   ./deploy/run-local.sh
#
# 说明：
#   HotGo 的 manifest/config、resource、logs 都是相对“工作目录”的路径，
#   必须在 server 目录下启动，否则会读不到 casbin.conf，HTTP 服务起不来。
#   这个脚本会先 cd 到 server，再 exec 二进制。
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVER_DIR="$ROOT/server"
BIN="$SERVER_DIR/temp/hotgo-local"

if [[ ! -f "$BIN" ]]; then
  echo "未找到 ${BIN}" >&2
  echo "请先执行：./deploy/build.sh --local" >&2
  exit 1
fi

if [[ ! -f "$SERVER_DIR/manifest/config/config.yaml" ]]; then
  echo "缺少 ${SERVER_DIR}/manifest/config/config.yaml" >&2
  exit 1
fi

if [[ ! -f "$SERVER_DIR/manifest/config/casbin.conf" ]]; then
  echo "缺少 ${SERVER_DIR}/manifest/config/casbin.conf" >&2
  exit 1
fi

cd "$SERVER_DIR"
exec "$BIN" "$@"
