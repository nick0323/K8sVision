#!/bin/bash

# K8sVision 多架构构建脚本
# 支持 AMD64 和 ARM64 架构

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
PLATFORMS="linux/amd64,linux/arm64"

echo -e "${BLUE}🚀 开始多架构构建 K8sVision...${NC}"
echo -e "${YELLOW}版本: ${VERSION}${NC}"
echo -e "${YELLOW}平台: ${PLATFORMS}${NC}"
echo ""

# 检查是否安装了buildx
if ! docker buildx version &> /dev/null; then
    echo -e "${RED}❌ Docker Buildx 未安装，请先安装 Docker Buildx${NC}"
    exit 1
fi

# 创建并使用buildx构建器
echo -e "${BLUE}🔧 设置构建器...${NC}"
docker buildx create --name k8svision-builder --use 2>/dev/null || docker buildx use k8svision-builder

# 构建并推送多架构镜像
echo -e "${BLUE}📦 构建并推送多架构镜像...${NC}"

# 构建后端镜像
echo -e "${YELLOW}构建后端镜像 (${PLATFORMS})...${NC}"
docker buildx build \
    --platform ${PLATFORMS} \
    --file Dockerfile.multiarch \
    --target runtime \
    --tag ${REGISTRY}/${PROJECT_NAME}-backend:${VERSION} \
    --tag ${REGISTRY}/${PROJECT_NAME}-backend:latest \
    --push \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 后端多架构镜像构建并推送成功${NC}"
else
    echo -e "${RED}❌ 后端多架构镜像构建失败${NC}"
    exit 1
fi

# 构建前端镜像
echo -e "${YELLOW}构建前端镜像 (${PLATFORMS})...${NC}"
docker buildx build \
    --platform ${PLATFORMS} \
    --file frontend/Dockerfile \
    --tag ${REGISTRY}/${PROJECT_NAME}-frontend:${VERSION} \
    --tag ${REGISTRY}/${PROJECT_NAME}-frontend:latest \
    --push \
    ./frontend

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 前端多架构镜像构建并推送成功${NC}"
else
    echo -e "${RED}❌ 前端多架构镜像构建失败${NC}"
    exit 1
fi

# 显示镜像信息
echo ""
echo -e "${BLUE}📊 多架构镜像信息:${NC}"
echo -e "${YELLOW}后端镜像:${NC}"
docker buildx imagetools inspect ${REGISTRY}/${PROJECT_NAME}-backend:${VERSION}

echo -e "${YELLOW}前端镜像:${NC}"
docker buildx imagetools inspect ${REGISTRY}/${PROJECT_NAME}-frontend:${VERSION}

echo ""
echo -e "${GREEN}🎉 多架构镜像构建完成！${NC}"
echo -e "${YELLOW}镜像已推送到:${NC}"
echo -e "  ${REGISTRY}/${PROJECT_NAME}-backend:${VERSION}"
echo -e "  ${REGISTRY}/${PROJECT_NAME}-frontend:${VERSION}"
echo ""
echo -e "${YELLOW}支持的架构:${NC}"
echo -e "  - linux/amd64 (Intel/AMD 64位)"
echo -e "  - linux/arm64 (ARM 64位)"
