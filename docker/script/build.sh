#!/usr/bin/env bash
# Usage: ./build.sh {RUN_NAME}

RUN_NAME="$1"
ROOT_DIR=$(pwd)

if [ -z "$RUN_NAME" ]; then
    echo "Error: Service name is required."
    exit 1
fi

# 进入模块
cd "./cmd/${RUN_NAME}" || exit
# 创建out目录
mkdir -p ${ROOT_DIR}/output/${RUN_NAME}

# 基于环境变量决定构建还是测试
if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    go build -o ${ROOT_DIR}/output/bin/${RUN_NAME}/${RUN_NAME}
else
    go test -c -covermode=set -o ${ROOT_DIR}/output/bin/${RUN_NAME}/${RUN_NAME} -coverpkg=./...
fi