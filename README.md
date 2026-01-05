# biliCLI - B站视频本地音频播放器

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Rust](https://img.shields.io/badge/rust-%23000000.svg?style=flat&logo=rust&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)
![Python](https://img.shields.io/badge/python-%2314354C.svg?style=flat&logo=python&logoColor=white)

## 🎯 项目简介

biliCLI 是一个专注于音频播放的终端式B站视频的音频分片播放工具，专为技术爱好者设计。它能够智能地将下载的B站视频分片文件组织成结构化目录树，并提供美观的TUI界面进行高效浏览和音频播放管理。无需图形界面或内置浏览器，性能友好

### ✨ 核心特色

- 🚀 **高性能扫描**：基于Rust的快速文件索引引擎
- 🎨 **优雅界面**：现代化的终端用户界面
- 🎵 **智能音频播放**：自动识别音频文件，提供沉浸式播放体验
- ⏯️ **播放控制**：支持暂停/继续、主动退出等交互控制
- 🎬 **视觉特效**：播放时显示十六进制刷屏效果
- 🔍 **便捷浏览**：快速定位音频内容
- 📊 **结构化组织**：自动构建清晰的音频目录树
- ⚡ **轻量高效**：多语言混合架构，性能卓越

### 🔧 开发环境

> **开发说明**：本项目使用 **Neovim** 作为主要编辑器开发，在 Windows 11 + WSL 环境下测试。使用 Windows 11 客户端缓存视频，WSL 运行管理工具。支持多平台运行，欢迎不同环境的用户反馈使用体验。
>

## 🚀 快速开始

### 📋 环境要求

#### 必需依赖
- **Rust** - 用于编译高性能文件扫描器
- **Go 1.20+** - 用于构建TUI界面
- **Python 3.8+** - 用于辅助脚本
- **ffplay** - 来自 ffmpeg，用于视频播放
- **xxd** - 用于十六进制显示

#### 推荐工具
```bash
# Ubuntu/Debian
sudo apt-get install pkg-config build-essential ffmpeg vim

# macOS
brew install pkg-config ffmpeg vim
```

### 🛠️ 安装步骤

#### ⚠️ 重要：路径配置说明
本项目各组件存在硬编码路径依赖，请务必按照以下结构编译和部署：

```
biliCLI/                    # 项目根目录
├── buildtree/              # Rust项目目录
│   └── target/release/     # 编译输出目录（必须）
│       └── buildtree       # Rust扫描器可执行文件
├── cmd/tui/mytui           # Go TUI程序（必须）
├── play                    # 播放脚本（必须）
├── fake_hex                # 十六进制显示工具（必须）
├── launch                  # 启动器（必须）
├── 01_read_config.py       # 配置读取脚本（必须）
├── 02_find_m4s.py          # 文件查找脚本（必须）
├── 03_detect_av.py         # 音视频检测脚本（必须）
├── 04_play.py              # Python播放脚本（必须）
├── buildtree/tree.json     # 索引文件（生成）
└── config.json             # 配置文件（必须）
```

#### 1. 获取项目代码
```bash
git clone https://github.com/diyiliumin/biliCLI.git
cd biliCLI
```

#### 2. 编译所有组件

##### 方法1：使用Makefile（推荐）
```bash
# 一键编译所有组件
make all

# 或者分别编译
make rust    # 编译Rust扫描器
make go      # 编译Go TUI
make c-tools # 编译C辅助工具
```

##### 方法2：手动编译（必须按顺序执行）
```bash
# 编译 Rust 扫描器（必须在此目录编译）
cd buildtree
cargo build --release
# 确保生成：buildtree/target/release/buildtree

# 编译 Go TUI（必须输出到指定位置）
cd ../cmd/tui
go build -o mytui
# 确保生成：cmd/tui/mytui

# 编译 C 辅助工具（必须在根目录编译）
cd ../..
gcc fake_hex.c -o fake_hex    # 必须在根目录
gcc launcher.c -o launch      # 必须在根目录
# 确保生成：fake_hex 和 launch
```

##### 方法3：验证编译结果
```bash
# 检查所有必需文件是否存在
ls -la buildtree/target/release/buildtree  # Rust扫描器
ls -la cmd/tui/mytui                       # Go TUI
ls -la fake_hex launch                     # C工具
ls -la 01_read_config.py 02_find_m4s.py 03_detect_av.py 04_play.py play  # 脚本文件
```

#### 3. 配置文件
编辑 `config.json`，设置你的B站视频下载目录：
```json
{
  "root": "/path/to/your/bilibili/downloads"
}
```

**常见路径示例：**
- **Windows**: `C:/Users/用户名/Videos/Bilibili`
- **Linux**: `~/Videos/Bilibili`
- **macOS**: `~/Movies/Bilibili`

#### 4. 路径依赖检查
```bash
# 运行依赖检查脚本（如果有）
./check_dependencies.sh

# 手动检查关键路径
echo "检查tree.json路径:"
ls -la buildtree/tree.json
echo "检查配置文件:"
ls -la config.json
echo "检查所有组件:"
ls -la buildtree/target/release/buildtree cmd/tui/mytui fake_hex launch play
```

## 📖 使用指南

### 🎯 基本流程

#### 首次使用 - 构建音频索引
```bash
# 运行扫描器构建索引（必须在项目根目录执行）
cd buildtree && cargo run --release
# 确保生成：buildtree/tree.json
```

#### 启动TUI界面
```bash
# 使用启动器（推荐方式）- 必须在项目根目录
./launch

```

#### 播放音频内容
```bash
# 通过TUI界面选择项目按p键播放
# 或直接调用play脚本（必须在项目根目录）
./play <CID>
```

#### 路径使用注意事项
```bash
# ⚠️ 重要：所有操作必须在项目根目录执行
cd /path/to/biliCLI  # 先进入项目根目录

# 正确用法：
./launch                    # ✅ 正确
./play 12345               # ✅ 正确

# 错误用法：
cd cmd/tui && ./mytui      # ❌ 错误 - 路径依赖会失败
./buildtree/target/release/buildtree  # ❌ 错误 - 相对路径错误
```

### ⌨️ 控制指南

#### TUI界面快捷键
| 快捷键 | 功能描述 |
|--------|----------|
| **j/k** | 上下移动光标 |
| **h** | 收起目录节点 |
| **l/Enter** | 展开目录节点或播放选中项 |
| **Space** | 选择/取消选中项目 |
| **p** | 播放选中项 |
| **q/Ctrl+C** | 退出程序 |

#### 播放时交互控制
| 快捷键 | 功能描述 |
|--------|----------|
| **p/P** | 暂停/继续播放 |
| **x/X** | 退出播放 |

### 🎮 播放模式

- **📺 顺序播放**：按照目录结构顺序播放视频
- **🎲 随机播放**：随机打乱播放顺序，发现惊喜
- **🎵 音频模式**：智能识别音频文件，提供沉浸式音频播放体验
- **⏯️ 播放控制**：支持暂停/继续播放功能（按p键）
- **🎬 视觉特效**：播放时显示十六进制刷屏效果
- **⌨️ 交互控制**：支持用户主动退出（按x键）

### 📁 目录结构

```
视频分组 (Group)
├── 视频系列 (Title)
│   ├── 分集1 (Item)
│   ├── 分集2 (Item)
│   └── 分集3 (Item)
└── 另一个系列
    ├── ost1
    └── ost2
```

## 🔧 常见问题解答

### ❌ 路径相关错误

#### ❓ 错误："找不到 tree.json 文件"
**原因**：运行路径不正确或索引未生成
**解决方案**：
```bash
# 确保在项目根目录
cd /path/to/biliCLI
# 构建索引
cd buildtree && cargo run --release
# 检查文件是否存在
ls -la buildtree/tree.json
```

#### ❓ 错误："找不到 01_read_config.py" 或其他脚本
**原因**：不在项目根目录运行
**解决方案**：
```bash
# 必须回到项目根目录
cd /path/to/biliCLI
# 重新运行
./launch
```

#### ❓ 错误："找不到 buildtree/target/release/buildtree"
**原因**：Rust组件未正确编译或路径错误
**解决方案**：
```bash
# 检查编译结果
cd buildtree
ls -la target/release/
# 如果不存在，重新编译
cargo build --release
```

### ❓ 启动时报错 "未检测到 buildtree/tree.json"
**解决方案**：运行 `cd buildtree && cargo run --release` 构建音频索引

### ❓ 播放时只有声音没有画面
**原因说明**：play脚本设计为音频播放器，会自动过滤视频文件，只播放音频内容

### ❓ TUI界面显示异常
**解决方案**：确保终端支持UTF-8编码，建议使用现代化终端模拟器

### ❓ 无法找到音频文件
**排查步骤**：
1. 检查 `config.json` 中的路径配置是否正确
2. 确认音频目录下存在 `videoInfo.json` 文件
3. 验证音频文件格式是否受支持（.m4s格式）

### ❓ 如何更新音频库
**操作方法**：添加新音频后，重新运行 `cd buildtree && cargo run --release` 构建索引

## 🧹 维护与清理

```bash
# 删除索引文件（需要重新构建）
rm buildtree/tree.json

# 清理Rust构建缓存
cd buildtree && cargo clean

# 清理Go编译文件
rm -rf cmd/tui/mytui

# 清理异常退出的播放进程
pkill -f ffplay
pkill -f play
```

## 🏗️ 项目架构

```
biliCLI/
├── buildtree/          # Rust高性能扫描器：分析视频文件生成索引
│   ├── Cargo.toml     # Rust项目配置
│   └── src/
│       └── main.rs    # 扫描器主程序
├── cmd/tui/           # Go TUI：现代化终端用户界面
│   └── main.go        # TUI主程序
├── internal/          # Go内部库
│   ├── model/         # 数据模型定义
│   ├── tree/          # 树形结构处理
│   └── ui/            # 用户界面组件
├── *.py               # Python辅助脚本集合
├── *.c                # C语言辅助工具
├── .nvim/             # Neovim配置文件（优化开发体验）
├── config.json        # 主配置文件
├── Makefile          # 自动化构建脚本
└── README.md          # 项目文档
```

## 🤝 贡献指南

我们欢迎所有形式的贡献，包括但不限于：

### 🐛 问题反馈
- 在 [Issues](https://github.com/diyiliumin/biliCLI/issues) 页面提交问题
- 提供详细的复现步骤和环境信息
- 附上相关的错误日志或截图

### 💡 功能建议
- 描述清楚使用场景和预期行为
- 解释该功能对用户价值的理解
- 欢迎提交原型设计或参考实现

### 🔧 代码贡献
1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的修改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

### 📋 开发规范
- 保持代码风格一致性
- 添加适当的注释说明
- 更新相关文档
- 确保通过基础测试

## 📄 许可证

本项目采用 MIT 许可证开源 - 详见 [LICENSE](LICENSE) 文件

## ⚖️ 免责声明

- 本工具仅为本地视频文件管理器，**不提供任何视频下载功能**
- 用户需自行通过官方渠道下载视频内容
- 使用B站视频内容请遵守相关版权规定和服务条款
- 项目仅供学习和交流使用

## 🙏 致谢

感谢以下开源项目的支持：
- [Rust](https://www.rust-lang.org/) - 高性能系统编程语言
- [Go](https://golang.org/) - 简洁高效的编程语言
- [Python](https://www.python.org/) - 强大的脚本语言
- [FFmpeg](https://ffmpeg.org/) - 多媒体处理框架

---

**⭐ 如果这个项目对你有帮助，请给我们一个Star！**

**Made with ❤️ by [ayazumi](https://github.com/diyiliumin) · 2026**
