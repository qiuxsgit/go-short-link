#!/bin/bash

# 进入web目录
cd web

# 安装依赖
echo "安装前端依赖..."
yarn

# 构建前端代码
echo "构建前端代码..."
yarn build

# 创建静态文件目录
mkdir -p ../static

# 复制构建后的文件到static目录
echo "复制构建后的文件到static目录..."
cp -r build/* ../static/

echo "前端构建完成！"