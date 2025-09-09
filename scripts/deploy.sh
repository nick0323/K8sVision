#!/bin/bash

# K8sVision 部署脚本
# 使用方法: ./scripts/deploy.sh [environment] [action]
# 环境: dev, staging, prod
# 操作: deploy, upgrade, rollback, delete

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

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安装"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        log_error "docker 未安装"
        exit 1
    fi
    
    log_success "依赖检查通过"
}

# 检查 Kubernetes 连接
check_k8s_connection() {
    log_info "检查 Kubernetes 连接..."
    
    if ! kubectl cluster-info &> /dev/null; then
        log_error "无法连接到 Kubernetes 集群"
        exit 1
    fi
    
    log_success "Kubernetes 连接正常"
}

# 创建命名空间
create_namespace() {
    local namespace=$1
    log_info "创建命名空间: $namespace"
    
    if kubectl get namespace "$namespace" &> /dev/null; then
        log_warning "命名空间 $namespace 已存在"
    else
        kubectl create namespace "$namespace"
        log_success "命名空间 $namespace 创建成功"
    fi
}

# 部署配置
deploy_config() {
    local environment=$1
    local namespace=$2
    
    log_info "部署配置..."
    
    # 创建 ConfigMap
    if [ -f "config-$environment.yaml" ]; then
        kubectl create configmap k8svision-config \
            --from-file=config.yaml="config-$environment.yaml" \
            -n "$namespace" \
            --dry-run=client -o yaml | kubectl apply -f -
        log_success "ConfigMap 部署成功"
    else
        log_warning "配置文件 config-$environment.yaml 不存在，使用默认配置"
        kubectl create configmap k8svision-config \
            --from-file=config.yaml=config.yaml \
            -n "$namespace" \
            --dry-run=client -o yaml | kubectl apply -f -
    fi
    
    # 创建 Secret
    kubectl create secret generic k8svision-auth \
        --from-literal=username="${K8SVISION_USERNAME:-admin}" \
        --from-literal=password="${K8SVISION_PASSWORD:-admin123}" \
        --from-literal=jwt-secret="${K8SVISION_JWT_SECRET:-k8svision-secret-key}" \
        -n "$namespace" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    log_success "Secret 部署成功"
}

# 部署 RBAC
deploy_rbac() {
    local namespace=$1
    
    log_info "部署 RBAC..."
    kubectl apply -f k8s/rbac.yaml
    log_success "RBAC 部署成功"
}

# 部署应用
deploy_app() {
    local environment=$1
    local namespace=$2
    
    log_info "部署应用..."
    
    # 部署后端
    kubectl apply -f k8s/backend-deployment.yaml
    kubectl apply -f k8s/backend-service.yaml
    
    # 部署前端
    kubectl apply -f k8s/frontend-deployment.yaml
    kubectl apply -f k8s/frontend-service.yaml
    
    # 部署 Ingress
    if [ -f "k8s/ingress-$environment.yaml" ]; then
        kubectl apply -f "k8s/ingress-$environment.yaml"
    else
        kubectl apply -f k8s/ingress.yaml
    fi
    
    # 部署 HPA
    if [ "$environment" = "prod" ]; then
        kubectl apply -f k8s/hpa.yaml
        log_info "HPA 部署成功"
    fi
    
    # 部署网络策略
    kubectl apply -f k8s/network-policy.yaml
    
    log_success "应用部署成功"
}

# 等待部署完成
wait_for_deployment() {
    local namespace=$1
    
    log_info "等待部署完成..."
    
    kubectl wait --for=condition=available --timeout=300s deployment/k8svision-backend -n "$namespace"
    kubectl wait --for=condition=available --timeout=300s deployment/k8svision-frontend -n "$namespace"
    
    log_success "部署完成"
}

# 检查部署状态
check_deployment() {
    local namespace=$1
    
    log_info "检查部署状态..."
    
    echo "=== Pods ==="
    kubectl get pods -n "$namespace"
    
    echo "=== Services ==="
    kubectl get svc -n "$namespace"
    
    echo "=== Ingress ==="
    kubectl get ingress -n "$namespace"
    
    echo "=== HPA ==="
    kubectl get hpa -n "$namespace" 2>/dev/null || echo "HPA 未部署"
}

# 升级应用
upgrade_app() {
    local environment=$1
    local namespace=$2
    local image_tag=${3:-latest}
    
    log_info "升级应用到版本: $image_tag"
    
    # 更新镜像
    kubectl set image deployment/k8svision-backend k8svision="k8svision:$image_tag" -n "$namespace"
    kubectl set image deployment/k8svision-frontend k8svision-frontend="k8svision-frontend:$image_tag" -n "$namespace"
    
    # 等待升级完成
    kubectl rollout status deployment/k8svision-backend -n "$namespace"
    kubectl rollout status deployment/k8svision-frontend -n "$namespace"
    
    log_success "应用升级成功"
}

# 回滚应用
rollback_app() {
    local namespace=$1
    
    log_info "回滚应用..."
    
    kubectl rollout undo deployment/k8svision-backend -n "$namespace"
    kubectl rollout undo deployment/k8svision-frontend -n "$namespace"
    
    # 等待回滚完成
    kubectl rollout status deployment/k8svision-backend -n "$namespace"
    kubectl rollout status deployment/k8svision-frontend -n "$namespace"
    
    log_success "应用回滚成功"
}

# 删除应用
delete_app() {
    local namespace=$1
    
    log_warning "删除应用..."
    
    read -p "确定要删除应用吗? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        kubectl delete namespace "$namespace"
        log_success "应用删除成功"
    else
        log_info "取消删除"
    fi
}

# 主函数
main() {
    local environment=${1:-dev}
    local action=${2:-deploy}
    local namespace="k8svision-$environment"
    
    log_info "开始 $action 操作，环境: $environment"
    
    case $action in
        deploy)
            check_dependencies
            check_k8s_connection
            create_namespace "$namespace"
            deploy_config "$environment" "$namespace"
            deploy_rbac "$namespace"
            deploy_app "$environment" "$namespace"
            wait_for_deployment "$namespace"
            check_deployment "$namespace"
            ;;
        upgrade)
            check_dependencies
            check_k8s_connection
            upgrade_app "$environment" "$namespace" "$3"
            check_deployment "$namespace"
            ;;
        rollback)
            check_dependencies
            check_k8s_connection
            rollback_app "$namespace"
            check_deployment "$namespace"
            ;;
        delete)
            check_dependencies
            check_k8s_connection
            delete_app "$namespace"
            ;;
        status)
            check_deployment "$namespace"
            ;;
        *)
            echo "使用方法: $0 [environment] [action]"
            echo "环境: dev, staging, prod"
            echo "操作: deploy, upgrade, rollback, delete, status"
            exit 1
            ;;
    esac
    
    log_success "操作完成"
}

# 运行主函数
main "$@"
