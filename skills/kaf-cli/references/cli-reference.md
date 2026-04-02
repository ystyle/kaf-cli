# kaf-cli CLI 参数详细参考

## 完整参数列表

```
kaf-cli [选项]
```

### 输入输出

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-filename` | string | (必填) | txt 文件路径，必须以 `.txt` 结尾 |
| `-out` | string | 书名 | 输出文件名，不需要包含格式后缀 |
| `-format` | string | `epub` | 输出格式：`epub`、`mobi`、`azw3`、`all`。可通过环境变量 `KAF_CLI_FORMAT` 修改默认值 |

### 书籍信息

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-bookname` | string | 文件名 | 书名。支持从知轩藏书格式文件名自动提取：`《书名》作者：作者.txt` |
| `-author` | string | `YSTYLE` | 作者名 |
| `-lang` | string | `zh` | 语言设置：en,de,fr,it,es,zh,ja,pt,ru,nl。可通过环境变量 `KAF_CLI_LANG` 修改默认值 |

### 封面设置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-cover` | string | `cover.png` | 封面图片。支持：本地图片路径（jpg/png）、`orly`（生成 O'Rly 风格封面，需联网）、`none`（不使用封面） |
| `-cover-orly-color` | string | 随机 | O'Rly 封面主题色，支持 1-16 编号和 hex 颜色代码 |
| `-cover-orly-idx` | int | 随机 | O'Rly 封面动物图案，0-41，查看图案：https://orly.nanmu.me |

### 章节匹配

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-match` | string | 自动 | 章节标题匹配正则表达式。不填自动识别常见格式 |
| `-volume-match` | string | `^第[0-9一二三四五六七八九十零〇百千两 ]+[卷部]` | 卷匹配规则。设为 `false` 禁用卷识别 |
| `-exclude` | string | 自动 | 排除无效章节/卷的正则表达式 |
| `-max` | uint | `35` | 标题最大字数，超过的标题会被忽略 |
| `-unknow-title` | string | `章节正文` | 无法识别的章节默认名称 |

### 排版样式

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-align` | string | `center` | 标题对齐方式：`left`、`center`、`right`。可通过环境变量 `KAF_CLI_ALIGN` 修改默认值 |
| `-indent` | uint | `2` | 段落缩进字数 |
| `-bottom` | string | `1em` | 段落间距，单位可以为 `em`、`px` |
| `-line-height` | string | `1.5rem` | 行高，用于控制行间距 |
| `-font` | string | 无 | 嵌入字体文件路径，嵌入后 epub 正文将使用该字体 |
| `-separate-chapter-number` | bool | `false` | 分离章节序号和标题样式，序号单独一行显示 |
| `-custom-css-file` | string | 无 | 自定义 CSS 文件路径，覆盖默认样式 |

### 其他

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-tips` | bool | `true` | 是否在书中添加软件教程 |

> **macOS 注意**: `-tips` 参数设为 false 的方法为 `kaf-cli -filename 小说.txt -tips=0`

## 简单模式

当只传入一个 `.txt` 文件参数时，进入简单模式：

```shell
kaf-cli 全职法师.txt
```

等同于拖拽模式，自动使用默认参数。

## 自定义章节匹配规则

匹配规则使用 Go 正则表达式语法。

### 自动识别的格式

kaf-cli 默认可以自动识别以下章节格式：
- 第X章、第X回、第X节、第X集、第X幕、第X卷、第X部
- Section xxx、Chapter xxx、Page xxx
- 纯数字标题（如 `123`）
- 引子、楔子、序章、最终章
- 番外

### 自定义匹配示例

| 章节格式 | 匹配规则 |
|----------|----------|
| 第x节 | `-match "第.{1,8}节"` |
| Section 1~100 | `-match "Section \d+"` |
| Chapter xxx | `-match "Chapter .{1,8}"` |

### 正则语法速查

| 符号 | 含义 |
|------|------|
| `.` | 任意字符 |
| `\d+` | 一个或多个数字 |
| `\w` | 中英文字符 |
| `{1,8}` | 前面的字符出现 1-8 次 |
| `*` | 前面的字符出现 0 次或多次 |
| `+` | 前面的字符出现 1 次或多次 |

## 环境变量

| 变量名 | 作用 | 示例 |
|--------|------|------|
| `KAF_CLI_ALIGN` | 默认标题对齐方式 | `left`、`center`、`right` |
| `KAF_CLI_LANG` | 默认语言 | `zh`、`en`、`ja` 等 |
| `KAF_CLI_FORMAT` | 默认输出格式 | `epub`、`mobi`、`azw3`、`all` |

## 自定义 CSS

通过 `-custom-css-file` 指定 CSS 文件覆盖默认样式。

### 可用 CSS 选择器

```css
/* 卷名样式 */
h2.volume {
    border-top: 5px solid red;
    border-bottom: 5px solid red;
    color: darkred;
}

/* 章节标题样式 */
h3.title {
    color: blue;
    border-bottom: 3px solid blue;
}

/* 章节序号样式 */
h3.title span.chapter-number {
    color: green;
}

/* 正文段落样式 */
.content {
    line-height: 2;
    text-indent: 4em;
}

/* 整体样式 */
body {
    margin: 2em;
}
```

### 使用方法

```shell
kaf-cli -filename 小说.txt -custom-css-file mystyle.css
```
