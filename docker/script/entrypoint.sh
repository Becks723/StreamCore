#!/usr/bin/env bash
CURDIR=$(pwd)

export KITEX_RUNTIME_ROOT=$CURDIR
export KITEX_LOG_DIR="$CURDIR/log"

if [ ! -d "$KITEX_LOG_DIR/app" ]; then
    mkdir -p "$KITEX_LOG_DIR/app"
fi

if [ ! -d "$KITEX_LOG_DIR/rpc" ]; then
    mkdir -p "$KITEX_LOG_DIR/rpc"
fi

# 环境变量 SERVICE 将由调用者设置
exec "$CURDIR/output/bin/$SERVICE/$SERVICE"