# kaf-mcp MCP 服务器使用指南

kaf-mcp 是 kaf-cli 的 MCP (Model Context Protocol) 服务器版本，可以集成到支持 MCP 的 AI 助手中（如 Cherry Studio、Claude Desktop、Cursor 等），通过 AI 对话直接转换小说。

## 安装 kaf-mcp

从 [GitHub Releases](https://github.com/ystyle/kaf-cli/releases?q=mcp-) 下载对应平台的 kaf-mcp。

各平台文件名：

| 文件名 | 平台 |
|--------|------|
| `kaf-mcp_*_windows_amd64.zip` | Windows 64 位 |
| `kaf-mcp_*_windows_386.zip` | Windows 32 位 |
| `kaf-mcp_*_darwin_amd64.zip` | macOS Intel |
| `kaf-mcp_*_darwin_arm64.zip` | macOS Apple Silicon |
| `kaf-mcp_*_linux_amd64.zip` | Linux 64 位 |
| `kaf-mcp_*_linux_arm64.zip` | Linux ARM64 |
| `kaf-mcp_*_linux_loong64.zip` | 龙芯 LoongArch64 |

解压后得到 `kaf-mcp` 可执行文件。

## MCP 配置

### 通用配置

在 AI 工具的 MCP 服务器配置中添加：

- **命令**: kaf-mcp 可执行文件的完整路径
- **环境变量**: `KAF_DIR=xxx`（xxx 替换为小说 txt 文件所在目录，转换后的文件也会输出到此目录）

### Cherry Studio 配置示例

1. 打开 Cherry Studio 设置
2. 添加 MCP 服务器
3. 命令填 kaf-mcp 路径
4. 在环境变量中添加 `KAF_DIR`，值为小说目录路径
5. 保存配置

### Claude Desktop 配置示例

在 `claude_desktop_config.json` 中添加：

```json
{
  "mcpServers": {
    "kaf-mcp": {
      "command": "/path/to/kaf-mcp",
      "env": {
        "KAF_DIR": "/path/to/novels"
      }
    }
  }
}
```

### Cursor 配置示例

在 Cursor 设置的 MCP 部分添加类似配置。

## MCP 提供的工具

### kaf_convert

将 TXT 小说文件转换为 EPUB/AZW3/MOBI 电子书格式。

**参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `filename` | string | 是 | txt 小说文件路径，支持相对路径（相对于 KAF_DIR） |
| `bookname` | string | 否 | 书名，为空时自动从文件名识别 |
| `author` | string | 否 | 作者，为空时自动从文件名识别 |
| `match` | string | 否 | 章节匹配正则表达式，不填自动识别，例：`第.{1,8}章` |
| `format` | string | 否 | 输出格式：`epub`、`mobi`、`azw3`、`all`（全部格式）。默认 epub |
| `cover` | string | 否 | 封面：本地图片路径（jpg/png）、`orly`（在线生成 O'Rly 风格封面，需联网）、`none`（无封面）。默认使用 cover.png |
| `cover-orly-color` | string | 否 | O'Rly 封面主题色，1-16 编号或 hex 颜色代码（如 `#ff6600`），不填随机 |
| `cover-orly-idx` | number | 否 | O'Rly 封面动物图案编号 0-41，不填随机 |
| `font` | string | 否 | 嵌入字体文件路径，嵌入后 epub 正文使用该字体 |
| `custom-css-file` | string | 否 | 自定义 CSS 文件路径，覆盖默认样式。可用选择器：`h2.volume`、`h3.title`、`.content`、`body` |
| `lang` | string | 否 | 书籍语言：en,de,fr,it,es,zh,ja,pt,ru,nl。默认 zh |
| `indent` | number | 否 | 段落缩进字数，默认 2 |
| `align` | string | 否 | 标题对齐：left、center、right。默认 center |
| `line-height` | string | 否 | 行高/行间距，如 `1.5rem`、`2`。默认 1.5rem |
| `bottom` | string | 否 | 段落间距，单位可为 em、px。默认 1em |
| `max` | number | 否 | 标题最大字数，超过的会被忽略。默认 35 |
| `out` | string | 否 | 输出文件名（不含后缀），默认为书名 |
| `volume-match` | string | 否 | 卷匹配正则，设为 `false` 禁用卷识别。默认自动匹配 |
| `exclude` | string | 否 | 排除无效章节/卷的正则表达式 |
| `unknow-title` | string | 否 | 未知章节默认名称，默认：章节正文 |
| `separate-chapter-number` | boolean | 否 | 是否分离章节序号和标题样式（序号单独一行），默认 false |
| `tips` | boolean | 否 | 是否在书中添加软件教程，默认 true |

**返回值：** 转换成功后返回生成的电子书文件路径。

### office_website

获取 kaf-mcp 的官方网站/代码仓库地址。

## 使用示例

配置完成后，在 AI 对话中直接说明需求即可：

```
帮我把 全职法师.txt 转成 epub
```

```
把目录下的 斗破苍穹.txt 转成电子书，作者是天蚕土豆
```

AI 会自动调用 kaf_convert 工具完成转换，并返回生成的 epub 文件路径。

## 文件名自动识别

MCP 版本同样支持从知轩藏书格式文件名自动提取书名和作者：
- `《希灵帝国》（校对版全本）作者：远瞳.txt`
- `《希灵帝国》作者：远瞳.txt`

## 注意事项

- MCP 版本支持 epub、mobi、azw3、all 四种输出格式，通过 `format` 参数指定
- `KAF_DIR` 环境变量决定了文件读取和输出的目录
- 相对路径的文件会从 `KAF_DIR` 目录读取
- 转换成功后 AI 助手会使用 markdown 链接格式返回文件路径，方便点击查看
