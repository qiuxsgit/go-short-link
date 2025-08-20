#!/bin/bash

# 测试配置
ADMIN_PORT=8081
ACCESS_PORT=8082

# 创建多个短链接进行测试
echo "创建多个短链接进行测试..."

# 创建第一个短链接（通过管理API）
echo "创建短链接 1..."
response1=$(curl -s -X POST http://localhost:$ADMIN_PORT/short-link/create \
  -H "Content-Type: application/json" \
  -d '{"link":"https://www.example.com/page1", "expire": 3600}')

echo "响应 1: $response1"
shortlink1=$(echo $response1 | grep -o 'http://[^"]*')
echo "生成的短链接 1: $shortlink1"

# 创建第二个短链接（通过管理API）
echo -e "\n创建短链接 2..."
response2=$(curl -s -X POST http://localhost:$ADMIN_PORT/short-link/create \
  -H "Content-Type: application/json" \
  -d '{"link":"https://www.example.com/page2", "expire": 3600}')

echo "响应 2: $response2"
shortlink2=$(echo $response2 | grep -o 'http://[^"]*')
echo "生成的短链接 2: $shortlink2"

# 测试IP白名单（使用非本地IP，应该被拒绝）
echo -e "\n测试IP白名单（使用伪造的非白名单IP）..."
response_denied=$(curl -s -X POST http://localhost:$ADMIN_PORT/short-link/create \
  -H "Content-Type: application/json" \
  -H "X-Forwarded-For: 8.8.8.8" \
  -d '{"link":"https://www.example.com/page3", "expire": 3600}')

echo "IP白名单测试响应: $response_denied"

# 访问第一个短链接（通过访问API）
echo -e "\n访问短链接 1..."
# 提取短码
shortcode1=$(echo $shortlink1 | grep -o '[^/]*$')
echo "短码 1: $shortcode1"
echo "curl -v http://localhost:$ACCESS_PORT/s/$shortcode1"
echo "这将重定向到 https://www.example.com/page1"

# 访问第二个短链接（通过访问API）
echo -e "\n访问短链接 2..."
# 提取短码
shortcode2=$(echo $shortlink2 | grep -o '[^/]*$')
echo "短码 2: $shortcode2"
echo "curl -v http://localhost:$ACCESS_PORT/s/$shortcode2"
echo "这将重定向到 https://www.example.com/page2"

# 创建更多短链接以测试缓存
echo -e "\n创建更多短链接以测试缓存..."
for i in {1..10}; do
  echo "创建短链接 $((i+2))..."
  curl -s -X POST http://localhost:$ADMIN_PORT/short-link/create \
    -H "Content-Type: application/json" \
    -d "{\"link\":\"https://www.example.com/page$((i+2))\", \"expire\": 3600}" > /dev/null
done

echo -e "\n已创建多个短链接，可以使用以下命令测试重定向:"
echo "curl -v http://localhost:$ACCESS_PORT/s/$shortcode1"
echo "curl -v http://localhost:$ACCESS_PORT/s/$shortcode2"

# 测试缓存性能
echo -e "\n测试缓存性能..."
echo "重复访问第一个短链接10次（应该从缓存中获取）:"
for i in {1..10}; do
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:$ACCESS_PORT/s/$shortcode1
done

# 测试两个服务是否正常工作
echo -e "\n测试服务状态..."
echo "管理API服务状态:"
curl -s -o /dev/null -w "状态码: %{http_code}\n" http://localhost:$ADMIN_PORT/short-link/create

echo "访问API服务状态:"
curl -s -o /dev/null -w "状态码: %{http_code}\n" http://localhost:$ACCESS_PORT/s/test
