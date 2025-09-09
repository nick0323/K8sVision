#!/bin/bash

# K8sVision 构建脚本
# 使用方法: ./scripts/build.sh [version] [platform]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 获取版本号
get_version() {
    if [ -n "$1" ]; then
        echo "$1"
    else
        git describe --tags --always --dirty 2>/dev/null || echo "dev-$(date +%Y%m%d-%H%M%S)"
    fi
}

# 构建前端
build_frontend() {
    local version=$1
    log_info "构建前端..."
    
    cd frontend
    
    # 安装依赖
    log_info "安装前端依赖..."
    npm ci
    
    # 设置版本号
    export VITE_APP_VERSION="$version"
    export VITE_APP_BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    
    # 构建前端
    log_info "构建前端应用..."
    npm run build
    
    cd ..
    log_success "前端构建完成"
}

# 构建后端
build_backend() {
    local version=$1
    local platform=$2
    log_info "构建后端..."
    
    # 设置构建参数
    local ldflags="-X main.version=$version -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -w -s"
    
    # 设置目标平台
    local target_os="linux"
    local target_arch="amd64"
    
    if [ -n "$platform" ]; then
        case $platform in
            linux/amd64)
                target_os="linux"
                target_arch="amd64"
                ;;
            linux/arm64)
                target_os="linux"
                target_arch="arm64"
                ;;
            darwin/amd64)
                target_os="darwin"
                target_arch="amd64"
                ;;
            darwin/arm64)
                target_os="darwin"
                target_arch="arm64"
                ;;
            windows/amd64)
                target_os="windows"
                target_arch="amd64"
                ;;
            *)
                log_error "不支持的平台: $platform"
                exit 1
                ;;
        esac
    fi
    
    # 设置环境变量
    export GOOS="$target_os"
    export GOARCH="$target_arch"
    export CGO_ENABLED=0
    
    # 构建后端
    log_info "构建后端应用 (${target_os}/${target_arch})..."
    go build -ldflags="$ldflags" -o "dist/k8svision-${target_os}-${target_arch}" main.go
    
    log_success "后端构建完成"
}

# 构建 Docker 镜像
build_docker() {
    local version=$1
    local platform=$2
    log_info "构建 Docker 镜像..."
    
    # 构建标签
    local tag="k8svision:$version"
    local latest_tag="k8svision:latest"
    
    # 构建镜像
    if [ -n "$platform" ]; then
        log_info "构建多平台镜像: $platform"
        docker buildx build --platform "$platform" -t "$tag" -t "$latest_tag" .
    else
        log_info "构建当前平台镜像"
        docker build -t "$tag" -t "$latest_tag" .
    fi
    
    log_success "Docker 镜像构建完成: $tag"
}

# 运行测试
run_tests() {
    log_info "运行测试..."
    
    # 后端测试
    log_info "运行后端测试..."
    go test -v ./...
    
    # 前端测试
    log_info "运行前端测试..."
    cd frontend
    npm test -- --coverage --watchAll=false
    cd ..
    
    log_success "测试完成"
}

# 代码检查
run_lint() {
    log_info "运行代码检查..."
    
    # Go 代码检查
    log_info "检查 Go 代码..."
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
    else
        log_warning "golangci-lint 未安装，跳过 Go 代码检查"
    fi
    
    # 前端代码检查
    log_info "检查前端代码..."
    cd frontend
    npm run lint
    cd ..
    
    log_success "代码检查完成"
}

# 生成文档
generate_docs() {
    log_info "生成文档..."
    
    # 生成 API 文档
    if command -v swag &> /dev/null; then
        log_info "生成 Swagger 文档..."
        swag init -g main.go -o docs
    else
        log_warning "swag 未安装，跳过 API 文档生成"
    fi
    
    log_success "文档生成完成"
}

# 清理构建文件
clean() {
    log_info "清理构建文件..."
    
    # 清理 Go 构建文件
    go clean -cache
    rm -rf dist/
    
    # 清理前端构建文件
    cd frontend
    rm -rf dist/
    cd ..
    
    # 清理 Docker 镜像
    docker image prune -f
    
    log_success "清理完成"
}

# 主函数
main() {
    local version=$(get_version "$1")
    local platform=$2
    local action=${3:-build}
    
    log_info "开始构建，版本: $version"
    
    case $action in
        build)
            build_frontend "$version"
            build_backend "$version" "$platform"
            build_docker "$version" "$platform"
            ;;
        frontend)
            build_frontend "$version"
            ;;
        backend)
            build_backend "$version" "$platform"
            ;;
        docker)
            build_docker "$version" "$platform"
            ;;
        test)
            run_tests
            ;;
        lint)
            run_lint
            ;;
        docs)
            generate_docs
            ;;
        clean)
            clean
            ;;
        all)
            run_lint
            run_tests
            build_frontend "$version"
            build_backend "$version" "$platform"
            build_docker "$version" "$platform"
            generate_docs
            ;;
        *)
            echo "使用方法: $0 [version] [platform] [action]"
            echo "版本: 可选，默认为 git tag 或时间戳"
            echo "平台: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64"
            echo "操作: build, frontend, backend, docker, test, lint, docs, clean, all"
            exit 1
            ;;
    esac
    
    log_success "构建完成"
}

# 运行主函数
main "$@"
