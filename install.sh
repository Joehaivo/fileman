#!/bin/bash

# FileMan 一键安装脚本
# 用法: curl -fsSL https://raw.githubusercontent.com/Joehaivo/fileman/main/install.sh | bash

set -e

REPO="Joehaivo/fileman"
BINARY_NAME="fm"
INSTALL_DIR="/usr/local/bin"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        *)          echo "unknown";;
    esac
}

# 检测架构
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64";;
        arm64|aarch64)  echo "arm64";;
        *)              echo "unknown";;
    esac
}

# 获取最新版本号
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "无法获取最新版本号"
        exit 1
    fi
    echo "$version"
}

# 下载并安装
install_fm() {
    local os="$1"
    local arch="$2"
    local version="$3"

    local archive_name="${BINARY_NAME}-${os}-${arch}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

    info "正在下载 ${archive_name}..."
    info "下载地址: ${download_url}"

    # 创建临时目录
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf ${tmp_dir}" EXIT

    local tmp_file="${tmp_dir}/${archive_name}"

    # 下载文件
    if ! curl -fsSL "$download_url" -o "$tmp_file"; then
        error "下载失败，请检查网络连接或版本是否存在"
        exit 1
    fi

    # 检查安装目录权限
    if [ ! -w "$INSTALL_DIR" ]; then
        warn "需要管理员权限安装到 ${INSTALL_DIR}"
        SUDO="sudo"
    else
        SUDO=""
    fi

    # 安装
    info "正在安装到 ${INSTALL_DIR}/${BINARY_NAME}..."
    $SUDO mv "$tmp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    $SUDO chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    success "安装完成！"
}

# 主函数
main() {
    echo ""
    echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║       FileMan 安装程序               ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"
    echo ""

    # 检测系统
    local os
    os=$(detect_os)
    if [ "$os" = "unknown" ]; then
        error "不支持的操作系统: $(uname -s)"
        exit 1
    fi
    info "检测到操作系统: ${os}"

    # 检测架构
    local arch
    arch=$(detect_arch)
    if [ "$arch" = "unknown" ]; then
        error "不支持的架构: $(uname -m)"
        exit 1
    fi
    info "检测到架构: ${arch}"

    # 获取版本
    local version
    version=$(get_latest_version)
    info "最新版本: ${version}"

    # 安装
    install_fm "$os" "$arch" "$version"

    # 验证安装
    echo ""
    if command -v ${BINARY_NAME} &> /dev/null; then
        success "运行 '${BINARY_NAME} --version' 验证安装"
        $(${BINARY_NAME} --version 2>/dev/null || echo "fm ${version}")
        echo ""
        echo -e "${GREEN}现在你可以使用 '${BINARY_NAME}' 命令启动 FileMan 了！${NC}"
        echo ""
    else
        warn "请重新打开终端或运行 'hash -r' 后使用 ${BINARY_NAME} 命令"
    fi
}

main "$@"