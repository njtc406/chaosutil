# Makefile
# Author: YR
# Version: 1.0
# Date: 2023-10-21
.PHONY: release debug clean help

SHELL = /bin/bash

## app
APP = app-query
## app名称
APP_NAME = ${APP}
## 源码路径
SRC_PATH = ../../cmd/query
## 编译导出路径
TARGET_PATH_RELEASE = ../../build/output/query/release
TARGET_PATH_DEBUG = ../../build/output/query/debug
TARGET_PATH = ${TARGET_PATH_RELEASE}
## 配置文件路径
TARGET_CONFIG_PATH = ${TARGET_PATH}/configs
CONFIG_PATH = ../../configs
## 需要拷贝的配置
CONFIGS = public query
## 系统类型
SYS_OS=linux
## cpu架构
CPU_ARC=amd64
## GOPRIVATE配置 用于配置私有仓库
PRI_GO=""

## 编译正式版本
release: create-target-dir clean build-os-release copy-configs

## 编译调试版本
debug: TARGET_PATH = ${TARGET_PATH_DEBUG}
debug: create-target-dir clean build-os-debug copy-configs

## 创建导出目录
create-target-dir:
	@mkdir -p $(TARGET_PATH)
	@mkdir -p $(TARGET_CONFIG_PATH)

## 拷贝配置文件
copy-configs:
	@for dir in $(CONFIGS); do \
		mkdir -p "$(TARGET_CONFIG_PATH)/$$dir"; \
		cp -rf "$(CONFIG_PATH)/$$dir/"* "$(TARGET_CONFIG_PATH)/$$dir/"; \
	done

## 编译可执行文件
build-os-release:
	CGO_ENABLED=0 GOOS=${SYS_OS} GOARCH=${CPU_ARC} \
	GOPROXY="https://goproxy.cn,direct" \
	GOPRIVATE=${PRI_GO} \
	go build -v -a -ldflags "-w -s" -o ${TARGET_PATH}/${APP_NAME} ${SRC_PATH}

## 编译可执行文件
build-os-debug:
	CGO_ENABLED=0 GOOS=${SYS_OS} GOARCH=${CPU_ARC} \
	GOPROXY="https://goproxy.cn,direct" \
	GOPRIVATE=${PRI_GO} \
	go build -tags "swagger" -v -a -ldflags "-w -s" -o ${TARGET_PATH}/${APP_NAME} ${SRC_PATH}

## 清理缓存和打包结果
clean:
	@go clean -cache
	@if [ -d ${TARGET_PATH_RELEASE} ] ; then rm -fr ${TARGET_PATH_RELEASE}/*; fi
	@if [ -d ${TARGET_PATH_DEBUG} ] ; then rm -fr ${TARGET_PATH_DEBUG}/*; fi

## 打印帮助文档
help:
	@echo "make - 编译 Go 代码, 生成可执行文件"
	@echo "make debug - 编译 Go 代码, 生成debug状态可执行文件"
	@echo "make clean - 清除编译的可执行文件和编译缓存"