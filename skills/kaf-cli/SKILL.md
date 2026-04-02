---
name: kaf-cli
description: |
  kaf-cli 是一个把 TXT 小说文本转成 EPUB/AZW3/MOBI 电子书格式的命令行工具。
  当用户需要：转换小说 txt 到电子书格式、批量生成 epub、设置书籍封面、匹配章节目录、自定义电子书样式、
  使用 kaf-cli 或 kaf-mcp 命令、了解 kaf-cli 的安装和参数用法时，都应该使用此 skill。
---

# kaf-cli - TXT 转电子书工具

kaf-cli 可以把 TXT 文本小说快速转换为 EPUB、AZW3、MOBI 格式的电子书。支持自动识别编码、章节、书名作者，也支持自定义封面、CSS 样式等高级功能。

## 工具安装

详细安装说明请参考 `references/installation.md`。

### 快速安装

| 平台 | 方式 |
|------|------|
| Windows | `winget install kaf-cli` 或从 [GitHub Releases](https://github.com/ystyle/kaf-cli/releases/latest) 下载 |
| Linux (Arch) | `yay -S kaf-cli` |
| macOS | 从 [GitHub Releases](https://github.com/ystyle/kaf-cli/releases/latest) 下载 darwin 版本 |
| Linux (通用) | 从 [GitHub Releases](https://github.com/ystyle/kaf-cli/releases/latest) 下载对应架构 |

### 下载文件名格式

- `kaf-cli_{version}_windows_amd64.zip` - Windows 64 位
- `kaf-cli_{version}_windows_386.zip` - Windows 32 位
- `kaf-cli_{version}_darwin_amd64.zip` - macOS Intel
- `kaf-cli_{version}_darwin_arm64.zip` - macOS Apple Silicon
- `kaf-cli_{version}_linux_amd64.zip` - Linux 64 位
- `kaf-cli_{version}_linux_arm64.zip` - Linux ARM64
- `kaf-cli_{version}_linux_loong64.zip` - 龙芯 LoongArch64

Windows 和 macOS 版本自带 kindlegen，可直接生成 mobi 格式。Linux 需单独安装 kindlegen 才能生成 mobi。

## 基本使用

### 拖拽模式（最简单）

把 txt 文件直接拖到 `kaf-cli.exe`（Windows）或 `kaf-cli`（Linux/Mac）可执行文件上即可自动转换。如果目录下有 `cover.png` 文件会自动添加为封面。

### 命令行简单模式

```shell
kaf-cli 全职法师.txt
```

### 命令行完整模式

```shell
kaf-cli -filename 全职法师.txt -author 乱 -format epub
```

## 常用参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-filename` | (必填) | txt 文件路径 |
| `-author` | `YSTYLE` | 作者 |
| `-bookname` | 文件名 | 书名，默认从文件名自动识别 |
| `-cover` | `cover.png` | 封面：本地图片路径、`orly`（在线生成 O'Rly 风格封面）、`none`（无封面） |
| `-format` | `epub` | 输出格式：`epub`、`mobi`、`azw3`、`all`（全部格式） |
| `-match` | 自动 | 章节匹配正则表达式，不填自动识别 |
| `-volume-match` | 自动 | 卷匹配正则，设为 `false` 禁用卷识别 |
| `-lang` | `zh` | 语言：en,de,fr,it,es,zh,ja,pt,ru,nl |
| `-indent` | `2` | 段落缩进字数 |
| `-align` | `center` | 标题对齐：left、center、right |
| `-out` | 书名 | 输出文件名（不含后缀） |
| `-font` | 无 | 嵌入字体文件路径 |
| `-line-height` | `1.5rem` | 行高/行间距 |
| `-bottom` | `1em` | 段落间距 |
| `-max` | `35` | 标题最大字数 |
| `-custom-css-file` | 无 | 自定义 CSS 文件路径 |
| `-separate-chapter-number` | `false` | 分离章节序号和标题样式 |
| `-tips` | `true` | 添加软件教程到书中 |
| `-exclude` | 自动 | 排除无效章节的正则 |

完整参数说明和更多示例请参考 `references/cli-reference.md`。

## 使用示例

### 基本转换

```shell
# 最简用法，自动识别书名和章节
kaf-cli 小说.txt

# 指定作者
kaf-cli -filename 小说.txt -author "作者名"

# 指定输出格式
kaf-cli -filename 小说.txt -format all
```

### 自定义章节匹配

```shell
# 章节格式为"第x节"
kaf-cli -filename ebook.txt -match "第.{1,8}节"

# 章节格式为 "Section 1" ~ "Section 100"
kaf-cli -filename ebook.txt -match "Section \d+"

# 章节格式为 "Chapter xxx"
kaf-cli -filename ebook.txt -match "Chapter .{1,8}"
```

### 使用 O'Rly 风格封面

```shell
# 自动生成随机封面
kaf-cli -filename 小说.txt -cover orly

# 指定封面主题色和动物图案
kaf-cli -filename 小说.txt -cover orly -cover-orly-color "#ff6600" -cover-orly-idx 5
```

### 自定义 CSS 样式

```shell
kaf-cli -filename 小说.txt -custom-css-file mystyle.css
```

可用 CSS 类名：
- `h2.volume` - 卷名样式
- `h3.title` - 章节标题样式
- `h3.title span.chapter-number` - 章节序号样式
- `.content` - 正文段落样式
- `body` - 整体样式

### 环境变量

可以通过环境变量修改默认值：

```shell
export KAF_CLI_ALIGN=left        # 默认标题对齐方式
export KAF_CLI_LANG=zh           # 默认语言
export KAF_CLI_FORMAT=epub       # 默认输出格式
```

## 文件名自动识别

支持知轩藏书格式文件名自动提取书名和作者：
- `《希灵帝国》（校对版全本）作者：远瞳.txt` → 书名：希灵帝国，作者：远瞳
- `《希灵帝国》作者：远瞳.txt` → 书名：希灵帝国，作者：远瞳

## MCP 版本（kaf-mcp）

kaf-cli 还提供了 MCP (Model Context Protocol) 服务器版本，可集成到 AI 助手中使用。详细说明请参考 `references/mcp-usage.md`。

### MCP 快速配置

1. 从 [GitHub Releases](https://github.com/ystyle/kaf-cli/releases/latest) 下载 kaf-mcp（与 kaf-cli 在同一个 Release）
2. 在 AI 工具中配置 MCP 服务器，命令指向 kaf-mcp 路径
3. 设置环境变量 `KAF_DIR` 为小说文件所在目录
4. 在 AI 对话中直接请求转换小说即可

## 常见问题

### 中文乱码
kaf-cli 自动识别字符编码，无需手动处理。

### 章节没有生成
使用 `-match` 参数指定章节匹配规则，例如 `-match "第.{1,8}章"`。

### MOBI 格式生成
需要 kindlegen 程序。Windows 和 macOS 版本已自带，Linux 需单独下载。

### 在任意位置使用
- **Windows**: 把 `kaf-cli.exe` 放到 `C:\Windows\` 下
- **Linux/Mac**: 把 kaf-cli 放到 PATH 目录中，如 `/usr/local/bin/`
