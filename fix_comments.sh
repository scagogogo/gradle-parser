#!/bin/bash

# 修复所有中文注释，确保以句号结尾
find pkg -name "*.go" -exec sed -i '' 's|// \([^/]*[^。.]\)$|// \1。|g' {} \;

# 修复特定的注释格式问题
find pkg -name "*.go" -exec sed -i '' 's|// \([^/]*\)\.。$|// \1。|g' {} \;

echo "Fixed comment formatting"
