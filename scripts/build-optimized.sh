#!/bin/bash

# K8sVision 优化构建脚本
# 使用优化后的Dockerfile构建镜像

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
REGISTRY="ghcr.io/nick0323"
PROJECT_NAME="k8svision"
VERSION=${1:-"latest"}
BUILD_TYPE=${2:-"separate"}  # separate 或 full
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🚀 开始构建 K8sVision 优化镜像...${NC}"
echo -e "${YELLOW}版本: ${VERSION}${NC}"
echo -e "${YELLOW}构建类型: ${BUILD_TYPE}${NC}"
echo -e "${YELLOW}构建时间: ${BUILD_DATE}${NC}"
echo -e "${YELLOW}Git提交: ${GIT_COMMIT}${NC}"
echo ""

# 构建参数
BUILD_ARGS="--build-arg BUILD_DATE=${BUILD_DATE} --build-arg GIT_COMMIT=${GIT_COMMIT}"

if [ "$BUILD_TYPE" = "full" ]; then
    # 构建完整镜像（包含前端）
    echo -e "${BLUE}📦 构建完整镜像（后端+前端）...${NC}"
    docker build \
        ${BUILD_ARGS} \
        --target runtime \
        -f Dockerfile.full \
        -t ${REGISTRY}/${PROJECT_NAME}-full:${VERSION} \
        -t ${REGISTRY}/${PROJECT_NAME}-full:latest \
        .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ 完整镜像构建成功${NC}"
    else
        echo -e "${RED}❌ 完整镜像构建失败${NC}"
        exit 1
    fi
else
    # 构建后端镜像
    echo -e "${BLUE}📦 构建后端镜像...${NC}"
    docker build \
        ${BUILD_ARGS} \
        --target runtime \
        -t ${REGISTRY}/${PROJECT_NAME}-backend:${VERSION} \
        -t ${REGISTRY}/${PROJECT_NAME}-backend:latest \
        .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 后端镜像构建成功${NC}"
else
    echo -e "${RED}❌ 后端镜像构建失败${NC}"
    exit 1
fi

# 构建前端镜像
echo -e "${BLUE}📦 构建前端镜像...${NC}"
docker build \
    ${BUILD_ARGS} \
    -t ${REGISTRY}/${PROJECT_NAME}-frontend:${VERSION} \
    -t ${REGISTRY}/${PROJECT_NAME}-frontend:latest \
    ./frontend

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 前端镜像构建成功${NC}"
else
    echo -e "${RED}❌ 前端镜像构建失败${NC}"
    exit 1
fi

# 显示镜像信息
echo ""
echo -e "${BLUE}📊 镜像信息:${NC}"
echo -e "${YELLOW}后端镜像:${NC}"
docker images ${REGISTRY}/${PROJECT_NAME}-backend:${VERSION} --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

echo -e "${YELLOW}前端镜像:${NC}"
docker images ${REGISTRY}/${PROJECT_NAME}-frontend:${VERSION} --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

# 安全扫描（如果安装了trivy）
if command -v trivy &> /dev/null; then
    echo ""
    echo -e "${BLUE}🔍 执行安全扫描...${NC}"
    echo -e "${YELLOW}扫描后端镜像:${NC}"
    trivy image ${REGISTRY}/${PROJECT_NAME}-backend:${VERSION} --severity HIGH,CRITICAL --exit-code 0
    
    echo -e "${YELLOW}扫描前端镜像:${NC}"
    trivy image ${REGISTRY}/${PROJECT_NAME}-frontend:${VERSION} --severity HIGH,CRITICAL --exit-code 0
fi

echo ""
echo -e "${GREEN}🎉 所有镜像构建完成！${NC}"
echo -e "${YELLOW}使用以下命令运行:${NC}"
echo -e "  docker-compose up -d"
echo ""
echo -e "${YELLOW}推送镜像到仓库:${NC}"
echo -e "  docker push ${REGISTRY}/${PROJECT_NAME}-backend:${VERSION}"
echo -e "  docker push ${REGISTRY}/${PROJECT_NAME}-frontend:${VERSION}"
