# kaf-cli 安装指南

## 包管理器安装（推荐）

### Windows - winget

```shell
winget install kaf-cli
```

### Arch Linux - AUR

```shell
yay -S kaf-cli
```

也支持龙芯 LoongArch64 架构。

## 手动安装

### 1. 下载

从 [GitHub Releases](https://github.com/ystyle/kaf-cli/releases/latest) 下载对应平台的压缩包。

各平台文件名：

| 文件名 | 平台 |
|--------|------|
| `kaf-cli_*_windows_amd64.zip` | Windows 64 位（含 kindlegen） |
| `kaf-cli_*_windows_386.zip` | Windows 32 位（含 kindlegen） |
| `kaf-cli_*_darwin_amd64.zip` | macOS Intel（含 kindlegen） |
| `kaf-cli_*_darwin_arm64.zip` | macOS Apple Silicon（含 kindlegen） |
| `kaf-cli_*_linux_amd64.zip` | Linux 64 位 |
| `kaf-cli_*_linux_arm64.zip` | Linux ARM64 |
| `kaf-cli_*_linux_loong64.zip` | 龙芯 LoongArch64 |
| `kaf-cli_*_wasm_wasip1.zip` | WASM WASI 版本 |

### 2. 解压

```shell
# Linux/Mac
unzip kaf-cli_*.zip
chmod +x kaf-cli

# Windows - 右键解压即可
```

### 3. 放到 PATH 目录（可选，方便在任何位置使用）

**Windows:**
- 把 `kaf-cli.exe` 和 `kindlegen.exe` 放到 `C:\Windows\` 下

**Linux/Mac:**
```shell
# 方法1: 复制到系统目录
sudo cp kaf-cli /usr/local/bin/

# 方法2: 添加到用户 PATH
# 在 ~/.bashrc 或 ~/.zshrc 末尾添加：
export PATH="$HOME/application:$PATH"
```

## kindlegen 安装（MOBI 格式支持）

Windows 和 macOS 版本的压缩包已自带 kindlegen，无需额外安装。

Linux 需要单独下载：
1. 从 [GitHub kindlegen 备份](https://github.com/ystyle/kaf-cli/releases/tag/kindlegen) 下载
2. 解压后把 `kindlegen` 放到与 `kaf-cli` 同一目录，或放到 PATH 目录中

## MCP 版本安装

从 [GitHub MCP Releases](https://github.com/ystyle/kaf-cli/releases?q=mcp-) 下载 kaf-mcp。文件名格式与 kaf-cli 类似，以 `kaf-mcp_` 开头。

## 百度网盘下载

所有版本（包括手机版 kaf 和 kaf-wifi）可从百度网盘下载：
`https://pan.baidu.com/s/1EPkLJ7WIJYdYtRHBEMqw0w?pwd=h4np`

## 验证安装

```shell
kaf-cli
# 应显示版本号和参数帮助信息
```
