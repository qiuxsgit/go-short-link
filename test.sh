#!/bin/bash

# 创建短链接
echo "创建短链接..."
response=$(curl -s -X POST http://localhost:8080/short-link/create \
  -H "Content-Type: application/json" \
  -d '{"link":"https://www.example.com", "expire": 60}')

echo "响应: $response"

# 从响应中提取短链接
shortlink=$(echo $response | grep -o 'http://[^"]*')
echo "生成的短链接: $shortlink"

# 访问短链接（会重定向到原始链接）
echo -e "\n使用短链接..."
echo "curl -v $shortlink"
echo "这将重定向到 https://www.example.com"