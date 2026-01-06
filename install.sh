#!/bin/bash
# biliCLI 一键安装脚本
# 适用于不想编译的用户

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🎵 bilimusicplayer-cli 一键安装脚本${NC}"
echo -e "${YELLOW}适用于免编译安装${NC}"
echo ""

# 检测操作系统
OS="unknown"
ARCH="unknown"

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]]; then
    OS="windows"
else
    echo -e "${RED}❌ 不支持的操作系统: $OSTYPE${NC}"
    exit 1
fi

# 检测架构
case $(uname -m) in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ 不支持的架构: $(uname -m)${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}📋 系统信息:${NC}"
echo "  操作系统: $OS"
echo "  架构: $ARCH"
echo ""

# GitHub仓库信息
GITHUB_REPO="diyiliumin/bilimusicplayer-cli"
GITHUB_API="https://api.github.com/repos/$GITHUB_REPO/releases/latest"

# 获取最新版本
echo -e "${YELLOW}🔍 检查最新版本...${NC}"
LATEST_RELEASE=$(curl -s "$GITHUB_API" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}❌ 无法获取最新版本信息${NC}"
    echo -e "${YELLOW}请检查网络连接或手动下载:${NC}"
    echo "https://github.com/$GITHUB_REPO/releases"
    exit 1
fi

echo -e "${GREEN}✅ 最新版本: $LATEST_RELEASE${NC}"

# 创建安装目录
INSTALL_DIR="$HOME/.local/bin/bilimusicplayer-cli"
mkdir -p "$INSTALL_DIR"

echo -e "${YELLOW}📦 安装目录: $INSTALL_DIR${NC}"

# 下载URL格式
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_RELEASE/biliCLI-${OS}-${ARCH}.tar.gz"

echo -e "${YELLOW}⬇️  下载预编译二进制...${NC}"
echo "下载地址: $DOWNLOAD_URL"

# 创建临时目录
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# 下载并解压
if command -v curl &> /dev/null; then
    curl -L -o biliCLI.tar.gz "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -O biliCLI.tar.gz "$DOWNLOAD_URL"
else
    echo -e "${RED}❌ 需要 curl 或 wget 来下载文件${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 解压文件...${NC}"
tar -xzf biliCLI.tar.gz

# 检查解压结果
if [ ! -f "buildtree/target/release/buildtree" ] || [ ! -f "cmd/tui/mytui" ] || [ ! -f "play" ]; then
    echo -e "${RED}❌ 下载的文件不完整或格式错误${NC}"
    echo -e "${YELLOW}请手动下载并解压:${NC}"
    echo "https://github.com/$GITHUB_REPO/releases"
    exit 1
fi

echo -e "${GREEN}✅ 文件完整性检查通过${NC}"

# 复制文件到安装目录
echo -e "${YELLOW}📂 安装文件...${NC}"
cp -r . "$INSTALL_DIR/"

# 创建符号链接
echo -e "${YELLOW}🔗 创建快捷方式...${NC}"
mkdir -p "$HOME/.local/bin"

# 创建启动脚本
cat > "$HOME/.local/bin/bilimusicplayer-cli" << 'EOF'
#!/bin/bash
# bilimusicplayer-cli 启动器
INSTALL_DIR="$HOME/.local/bin/bilimusicplayer-cli"
cd "$INSTALL_DIR"
./launch "$@"
EOF

chmod +x "$HOME/.local/bin/bilicli"

# 创建配置文件模板
if [ ! -f "$INSTALL_DIR/config.json" ]; then
    cat > "$INSTALL_DIR/config.json" << 'EOF'
{
  "root": "/path/to/your/bilibili/downloads"
}
EOF
fi

echo ""
echo -e "${GREEN}🎉 安装完成！${NC}"
echo ""
echo -e "${YELLOW}📋 后续步骤:${NC}"
echo "1. 编辑配置文件: $INSTALL_DIR/config.json"
echo "2. 设置下载目录路径"
echo "3. 运行: $HOME/.local/bin/bilicli"
echo ""
echo -e "${YELLOW}⚠️  重要提醒:${NC}"
echo "- 首次使用需要构建索引，在程序中按提示操作"
echo "- 确保已安装依赖: ffplay, python3"
echo "- 所有组件必须在安装目录内，不要移动单独文件"
echo ""

# 添加到PATH的建议
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo -e "${YELLOW}💡 建议将 $HOME/.local/bin 添加到 PATH:${NC}"
    echo "echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
    echo "source ~/.bashrc"
fi

# 清理
cd /
rm -rf "$TEMP_DIR"

echo -e "${GREEN}✨ 享受 bilimusicplayer-cli 吧！${NC}"