@echo off
REM K8sVision 构建脚本 (Windows)
REM 使用方法: scripts\build.bat [version] [platform] [action]

setlocal enabledelayedexpansion

REM 设置颜色
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "NC=[0m"

REM 日志函数
:log_info
echo %BLUE%[INFO]%NC% %~1
goto :eof

:log_success
echo %GREEN%[SUCCESS]%NC% %~1
goto :eof

:log_warning
echo %YELLOW%[WARNING]%NC% %~1
goto :eof

:log_error
echo %RED%[ERROR]%NC% %~1
goto :eof

REM 获取版本号
:get_version
if "%~1" neq "" (
    set "VERSION=%~1"
) else (
    for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set "VERSION=%%i"
    if "!VERSION!"=="" set "VERSION=dev-%date:~0,4%%date:~5,2%%date:~8,2%-%time:~0,2%%time:~3,2%%time:~6,2%"
)
goto :eof

REM 构建前端
:build_frontend
call :log_info "构建前端..."
cd frontend

REM 安装依赖
call :log_info "安装前端依赖..."
call npm ci
if errorlevel 1 goto :error

REM 设置版本号
set "VITE_APP_VERSION=%VERSION%"
set "VITE_APP_BUILD_TIME=%date:~0,4%-%date:~5,2%-%date:~8,2%T%time:~0,2%:%time:~3,2%:%time:~6,2%Z"

REM 构建前端
call :log_info "构建前端应用..."
call npm run build
if errorlevel 1 goto :error

cd ..
call :log_success "前端构建完成"
goto :eof

REM 构建后端
:build_backend
call :log_info "构建后端..."

REM 设置构建参数
set "LDFLAGS=-X main.version=%VERSION% -X main.buildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%T%time:~0,2%:%time:~3,2%:%time:~6,2%Z -w -s"

REM 设置目标平台
set "TARGET_OS=windows"
set "TARGET_ARCH=amd64"

if "%~1" neq "" (
    if "%~1"=="linux/amd64" (
        set "TARGET_OS=linux"
        set "TARGET_ARCH=amd64"
    ) else if "%~1"=="linux/arm64" (
        set "TARGET_OS=linux"
        set "TARGET_ARCH=arm64"
    ) else if "%~1"=="darwin/amd64" (
        set "TARGET_OS=darwin"
        set "TARGET_ARCH=amd64"
    ) else if "%~1"=="darwin/arm64" (
        set "TARGET_OS=darwin"
        set "TARGET_ARCH=arm64"
    ) else if "%~1"=="windows/amd64" (
        set "TARGET_OS=windows"
        set "TARGET_ARCH=amd64"
    ) else (
        call :log_error "不支持的平台: %~1"
        exit /b 1
    )
)

REM 设置环境变量
set "GOOS=%TARGET_OS%"
set "GOARCH=%TARGET_ARCH%"
set "CGO_ENABLED=0"

REM 构建后端
call :log_info "构建后端应用 (%TARGET_OS%/%TARGET_ARCH%)..."
if not exist "dist" mkdir dist
go build -ldflags="%LDFLAGS%" -o "dist\k8svision-%TARGET_OS%-%TARGET_ARCH%.exe" main.go
if errorlevel 1 goto :error

call :log_success "后端构建完成"
goto :eof

REM 构建 Docker 镜像
:build_docker
call :log_info "构建 Docker 镜像..."

REM 构建标签
set "TAG=k8svision:%VERSION%"
set "LATEST_TAG=k8svision:latest"

REM 构建镜像
if "%~1" neq "" (
    call :log_info "构建多平台镜像: %~1"
    docker buildx build --platform "%~1" -t "%TAG%" -t "%LATEST_TAG%" .
) else (
    call :log_info "构建当前平台镜像"
    docker build -t "%TAG%" -t "%LATEST_TAG%" .
)
if errorlevel 1 goto :error

call :log_success "Docker 镜像构建完成: %TAG%"
goto :eof

REM 运行测试
:run_tests
call :log_info "运行测试..."

REM 后端测试
call :log_info "运行后端测试..."
go test -v ./...
if errorlevel 1 goto :error

REM 前端测试
call :log_info "运行前端测试..."
cd frontend
call npm test -- --coverage --watchAll=false
if errorlevel 1 goto :error
cd ..

call :log_success "测试完成"
goto :eof

REM 代码检查
:run_lint
call :log_info "运行代码检查..."

REM Go 代码检查
call :log_info "检查 Go 代码..."
where golangci-lint >nul 2>&1
if errorlevel 1 (
    call :log_warning "golangci-lint 未安装，跳过 Go 代码检查"
) else (
    golangci-lint run
    if errorlevel 1 goto :error
)

REM 前端代码检查
call :log_info "检查前端代码..."
cd frontend
call npm run lint
if errorlevel 1 goto :error
cd ..

call :log_success "代码检查完成"
goto :eof

REM 生成文档
:generate_docs
call :log_info "生成文档..."

REM 生成 API 文档
where swag >nul 2>&1
if errorlevel 1 (
    call :log_warning "swag 未安装，跳过 API 文档生成"
) else (
    call :log_info "生成 Swagger 文档..."
    swag init -g main.go -o docs
    if errorlevel 1 goto :error
)

call :log_success "文档生成完成"
goto :eof

REM 清理构建文件
:clean
call :log_info "清理构建文件..."

REM 清理 Go 构建文件
go clean -cache
if exist "dist" rmdir /s /q dist

REM 清理前端构建文件
cd frontend
if exist "dist" rmdir /s /q dist
cd ..

REM 清理 Docker 镜像
docker image prune -f

call :log_success "清理完成"
goto :eof

REM 错误处理
:error
call :log_error "构建失败"
exit /b 1

REM 主函数
:main
set "VERSION="
set "PLATFORM="
set "ACTION=build"

if "%~1" neq "" set "VERSION=%~1"
if "%~2" neq "" set "PLATFORM=%~2"
if "%~3" neq "" set "ACTION=%~3"

call :get_version "%VERSION%"

call :log_info "开始构建，版本: %VERSION%"

if "%ACTION%"=="build" (
    call :build_frontend
    call :build_backend "%PLATFORM%"
    call :build_docker "%PLATFORM%"
) else if "%ACTION%"=="frontend" (
    call :build_frontend
) else if "%ACTION%"=="backend" (
    call :build_backend "%PLATFORM%"
) else if "%ACTION%"=="docker" (
    call :build_docker "%PLATFORM%"
) else if "%ACTION%"=="test" (
    call :run_tests
) else if "%ACTION%"=="lint" (
    call :run_lint
) else if "%ACTION%"=="docs" (
    call :generate_docs
) else if "%ACTION%"=="clean" (
    call :clean
) else if "%ACTION%"=="all" (
    call :run_lint
    call :run_tests
    call :build_frontend
    call :build_backend "%PLATFORM%"
    call :build_docker "%PLATFORM%"
    call :generate_docs
) else (
    echo 使用方法: %0 [version] [platform] [action]
    echo 版本: 可选，默认为 git tag 或时间戳
    echo 平台: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
    echo 操作: build, frontend, backend, docker, test, lint, docs, clean, all
    exit /b 1
)

call :log_success "构建完成"
goto :eof

REM 运行主函数
call :main %*
