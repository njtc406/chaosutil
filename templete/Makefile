.PHONY: clean help

APP = app
APP_NAME = ${APP}

## 编译打包
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct" \
	GOPRIVATE="gitlab.chivalrous-tech.com" \
	go build -v -a -ldflags "-w -s" -o ./release/${APP_NAME} ./src

## 清理缓存和打包结果
clean:
	@go clean -cache
	@if [ -f ./release/${APP_NAME} ] ; then rm ./release/${APP_NAME}; fi

## 打印帮助文档
help:
	@echo "make - 编译 Go 代码, 生成linux二进制文件"
	@echo "make clean - 移除编译的二进制文件"