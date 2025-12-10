#!/bin/bash

# 设置 Go 代理
export GOPROXY=https://goproxy.io,direct

# 禁用 CGO 以避免编译错误
export CGO_ENABLED=0

# 清理并下载依赖
echo "整理依赖..."
go mod tidy

# 构建项目
echo "构建项目..."
go build -o go-short-link

echo "构建完成！"
