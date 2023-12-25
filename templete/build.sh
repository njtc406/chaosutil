#!/bin/bash
# build.sh
# Author: YR
# Version: 1.0
# Date: 2023-10-21

## app
APP=app-query
## app名称
APP_NAME=${APP}
## 源码路径
SRC_PATH=../../cmd/query
## 编译导出路径
TARGET_PATH=../../build/output/query
TARGET_APP_PATH=${TARGET_PATH}/app
## 需要拷贝的配置
CONFIGS=(
query
)
## 配置文件路径
TARGET_CONFIG_PATH=${TARGET_APP_PATH}/configs
CONFIG_PATH=../../configs
## docker-compose.yml文件路径
DOCKER_COMPOSE_PATH=./docker-compose.yml
## 系统类型
SYS_OS=linux
## cpu架构
CPU_ARC=amd64
## GOPRIVATE配置 用于配置私有仓库
PRI_GO=""

function release() {
    create_target_dir
    clean
    build-os-release
    copy-configs
}

function debug() {
    create_target_dir
    clean
    build-os-debug
    copy-configs
}

function clean() {
    rm -rf ${TARGET_APP_PATH}
}

function help() {
    echo "Usage: $0 [release|debug|clean|help]"
}

function create_target_dir() {
    mkdir -p ${TARGET_APP_PATH}
}

function copy-configs() {
    # shellcheck disable=SC2068
    for dir in ${CONFIGS[@]}; do
        mkdir -p ${TARGET_CONFIG_PATH}/$dir
        cp -rf ${CONFIG_PATH}/$dir/* ${TARGET_CONFIG_PATH}/$dir/
    done

    if [ -f ${DOCKER_COMPOSE_PATH} ]; then
        cp -f ${DOCKER_COMPOSE_PATH} ${TARGET_PATH}/
    fi
}

function build-os-release() {
    CGO_ENABLED=0 GOOS=${SYS_OS} GOARCH=${CPU_ARC} \
    GOPROXY="https://goproxy.cn,direct" \
    GOPRIVATE=${PRI_GO} \
    go build -v -a -ldflags "-w -s" -o ${TARGET_APP_PATH}/${APP_NAME} ${SRC_PATH}
}

function build-os-debug() {
    CGO_ENABLED=0 GOOS=${SYS_OS} GOARCH=${CPU_ARC} \
    GOPROXY="https://goproxy.cn,direct" \
    GOPRIVATE=${PRI_GO} \
    go build -v -a -ldflags "-w -s" -o ${TARGET_APP_PATH}/${APP_NAME} ${SRC_PATH}
}

case "$1" in
    release)
        release
        ;;
    debug)
        debug
        ;;
    clean)
        clean
        ;;
    help)
        help
        ;;
    *)
        help
        ;;
esac
